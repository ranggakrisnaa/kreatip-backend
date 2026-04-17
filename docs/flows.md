# Kreatip Backend — Flow Diagrams

Semua sequence diagram di bawah ini fokus di backend (clean architecture layer: Controller → Usecase → Repository/Gateway).

---

## 1. Register + Verifikasi Email

```mermaid
sequenceDiagram
    autonumber
    participant C as Client
    participant Ctl as AuthController
    participant UC as AuthUsecase
    participant Repo as UserRepository
    participant MQ as Asynq Queue
    participant W as EmailWorker
    participant MG as EmailGateway

    C->>Ctl: POST /api/v1/auth/register {email, password}
    Ctl->>Ctl: validate request
    Ctl->>UC: Register(ctx, req)
    UC->>Repo: FindByEmail(email)
    Repo-->>UC: not found
    UC->>UC: hash password (bcrypt)
    UC->>Repo: Create(user)
    UC->>UC: generate verification token
    UC->>MQ: Enqueue SendVerificationEmail
    UC-->>Ctl: UserResponse
    Ctl-->>C: 201 Created

    MQ->>W: dispatch job
    W->>MG: send email
    MG-->>W: ok
```

---

## 2. Donation Flow (Guest → Xendit → Webhook → Ledger → Alert)

```mermaid
sequenceDiagram
    autonumber
    participant D as Donor
    participant Ctl as DonationController
    participant UC as DonationUsecase
    participant CR as CreatorRepo
    participant DR as DonationRepo
    participant PG as XenditGateway
    participant PR as PaymentRepo

    D->>Ctl: POST /api/v1/donations {username, amount, message, method}
    Ctl->>Ctl: validate, recaptcha
    Ctl->>UC: CreateDonation(req)
    UC->>CR: FindByUsername(username)
    UC->>UC: compute fee & net amount
    UC->>DR: Insert(donation status=pending) [tx]
    UC->>PG: CreateInvoice/Charge(amount, method)
    PG-->>UC: external_id, checkout_url
    UC->>PR: Insert(payment) [tx]
    UC-->>Ctl: DonationResponse{checkout_url}
    Ctl-->>D: 200 OK (redirect ke checkout_url)
    Note over D,PG: Donor bayar di halaman Xendit
```

### 2a. Webhook Handler (uang diterima)

```mermaid
sequenceDiagram
    autonumber
    participant X as Xendit
    participant WC as WebhookController
    participant UC as PaymentUsecase
    participant WR as WebhookEventRepo
    participant DR as DonationRepo
    participant LU as LedgerUsecase
    participant WU as WalletUsecase
    participant AU as AlertUsecase
    participant WS as WS Hub (OBS)
    participant MQ as Asynq Queue

    X->>WC: POST /api/v1/webhooks/xendit (signed)
    WC->>WC: verify X-CALLBACK-TOKEN / signature
    WC->>UC: HandleWebhook(payload)
    UC->>WR: UpsertByIdempotencyKey(event)
    alt already processed
        WR-->>UC: duplicate
        UC-->>WC: 200 OK (no-op)
    else new event
        UC->>DR: UpdateStatus(donation, paid) [tx start]
        UC->>LU: PostDonationEntries(donation) [same tx]
        LU->>WU: IncrementBalanceCache
        UC->>WR: MarkProcessed [tx commit]
        UC->>AU: BroadcastDonation(donation)
        AU->>WS: push to creator room
        UC->>MQ: Enqueue EmailReceipt + EmailCreatorNotif
    end
    UC-->>WC: 200 OK
    WC-->>X: 200 OK
```

Key rules:
- **Idempotency**: `webhook_events.idempotency_key` = `provider + external_id + event_type`.
- **All-or-nothing**: update donation, insert ledger entries, mark webhook processed — dalam **satu DB transaction**.
- **Retry-safe**: Xendit akan retry jika non-2xx; handler harus idempotent.

---

## 3. Streamer Alert Overlay (OBS)

```mermaid
sequenceDiagram
    autonumber
    participant OBS as OBS Browser Source
    participant WC as WS Controller
    participant AU as AlertUsecase
    participant Repo as AlertTokenRepo
    participant Hub as WS Hub (in-memory)
    participant PU as PaymentUsecase

    OBS->>WC: WS /ws/alert?token=xxx
    WC->>AU: AuthorizeToken(token)
    AU->>Repo: FindByToken(token)
    Repo-->>AU: userID, settings
    WC->>Hub: Register(userID, conn)
    WC-->>OBS: connected

    Note over PU: Ketika donasi paid (lihat flow 2a)
    PU->>AU: BroadcastDonation(donation)
    AU->>Hub: Publish(userID, event)
    Hub->>WC: deliver
    WC->>OBS: { type: "donation", donor, amount, message }
    OBS->>OBS: render animasi + TTS via Web Speech API
```

Catatan: WS hub cukup in-memory untuk MVP (single instance). Untuk horizontal scaling di v2 pakai Redis Pub/Sub.

---

## 4. Withdrawal Flow

```mermaid
sequenceDiagram
    autonumber
    participant C as Creator
    participant Ctl as WithdrawalController
    participant UC as WithdrawalUsecase
    participant WU as WalletUsecase
    participant LU as LedgerUsecase
    participant WR as WithdrawalRepo
    participant Admin as Admin
    participant DG as DisbursementGateway
    participant MQ as Asynq Queue

    C->>Ctl: POST /api/v1/withdrawals {amount, bank_account_id, totp_code}
    Ctl->>UC: Request(req)
    UC->>UC: verify TOTP
    UC->>WU: GetAvailableBalance(userID)
    UC->>UC: validate amount >= min, <= balance
    UC->>WR: Insert(status=requested) [tx]
    UC->>LU: HoldFunds(userID, amount) [same tx, move creator_balance → pending_withdrawal]
    UC-->>Ctl: 201
    Ctl-->>C: withdrawal pending approval

    Admin->>Ctl: POST /api/v1/admin/withdrawals/:id/approve
    Ctl->>UC: Approve(id, adminID)
    UC->>WR: UpdateStatus(approved)
    UC->>MQ: Enqueue ProcessDisbursement

    MQ->>UC: ProcessDisbursement job
    UC->>DG: Disburse(bank_account, amount)
    DG-->>UC: external_id, status=processing
    UC->>WR: Update(processing, external_id)

    Note over DG: Xendit callback webhook
    DG->>Ctl: POST /api/v1/webhooks/disbursement
    Ctl->>UC: HandleDisbursementCallback
    alt success
        UC->>LU: SettleWithdrawal [pending_withdrawal → payment_gateway]
        UC->>WR: Update(success)
    else failed
        UC->>LU: ReleaseFunds [pending_withdrawal → creator_balance]
        UC->>WR: Update(failed)
    end
```

---

## 5. Daily Reconciliation Job

```mermaid
sequenceDiagram
    autonumber
    participant Cron as Cron (Asynq scheduler)
    participant W as ReconcileWorker
    participant LR as LedgerRepo
    participant WR as WalletRepo
    participant PG as Xendit API
    participant Log as AuditLog

    Cron->>W: trigger daily 02:00 WIB
    W->>PG: Fetch settlement report (yesterday)
    W->>LR: Compute per-wallet balance from entries
    W->>WR: List all wallets
    loop per wallet
        W->>W: compare ledger sum vs wallet.balance_available
        alt mismatch
            W->>Log: write discrepancy
            W->>WR: correct balance cache
        end
    end
    W->>W: compare internal totals vs Xendit settlement
    W->>Log: write daily report
```

---

## 6. JWT Auth Middleware

```mermaid
sequenceDiagram
    autonumber
    participant C as Client
    participant MW as AuthMiddleware
    participant JWT as JWT Lib
    participant Redis as Redis (blacklist)
    participant Ctl as Controller

    C->>MW: Request + Authorization: Bearer xxx
    MW->>JWT: Parse & verify signature
    alt invalid
        MW-->>C: 401
    else valid
        MW->>Redis: GET blacklist:{jti}
        alt blacklisted
            MW-->>C: 401
        else ok
            MW->>MW: set ctx.userID, ctx.role
            MW->>Ctl: next()
        end
    end
```

---

## 7. Refresh Token Rotation

```mermaid
sequenceDiagram
    autonumber
    participant C as Client
    participant Ctl as AuthController
    participant UC as AuthUsecase
    participant Repo as RefreshTokenRepo

    C->>Ctl: POST /api/v1/auth/refresh {refresh_token}
    Ctl->>UC: Refresh(token)
    UC->>Repo: FindByHash(sha256(token))
    alt not found / revoked / expired
        UC-->>Ctl: 401
    else valid
        UC->>Repo: Revoke(old)
        UC->>UC: generate new access + refresh
        UC->>Repo: Insert(new refresh hash)
        UC-->>Ctl: tokens
    end
    Ctl-->>C: 200 {access, refresh}
```
