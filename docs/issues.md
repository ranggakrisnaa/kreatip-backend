# Kreatip Backend — Issue Plan (MVP, Backend-Only)

Setiap entry siap di-paste ke GitHub Issue. Label ditulis di awal, lalu `## Description`, `## Acceptance Criteria`, `## Notes`.

Urutan milestone sengaja menyelaraskan layer clean architecture: migration → entity → repository → usecase → delivery. Dalam satu fitur, kerjakan bottom-up.

---

## Milestone 0 — Project Foundation

### #1 Init Go module + struktur clean architecture
**Labels:** `setup`, `chore`
**Description:**
Scaffold project mengikuti [khannedy/golang-clean-architecture](https://github.com/khannedy/golang-clean-architecture) dengan adaptasi di `docs/architecture.md`.
**Acceptance Criteria:**
- [ ] `go mod init github.com/<org>/kreatip-backend`
- [ ] Folder `cmd/web`, `cmd/worker`, `internal/{config,entity,model,repository,usecase,delivery,gateway}`, `db/migrations`, `test/`
- [ ] Fiber `"Hello World"` jalan di `GET /` saat `go run ./cmd/web`
- [ ] `Makefile` dengan target `run`, `build`, `test`, `migrate-up`, `migrate-down`, `lint`

### #2 Setup config loader (Viper)
**Labels:** `setup`
- [ ] `internal/config/config.go` load `config.json` + override env
- [ ] Validasi field wajib saat startup (fail-fast)
- [ ] `.env.example` + `config.json.example` committed

### #3 Setup Docker Compose untuk dev
**Labels:** `setup`, `devex`
- [ ] Postgres 16, Redis 7, MailHog
- [ ] Volume persistent untuk Postgres
- [ ] README singkat: `docker compose up -d`

### #4 Setup GORM + koneksi Postgres
**Labels:** `setup`, `db`
- [ ] `internal/config/database.go` bikin `*gorm.DB` dari config
- [ ] Pool tuning dari config
- [ ] Logger GORM diarahkan ke Logrus
- [ ] Health check query `SELECT 1` saat startup

### #5 Setup Redis client
**Labels:** `setup`
- [ ] `internal/config/redis.go` bikin `*redis.Client`
- [ ] Ping saat startup

### #6 Setup golang-migrate + baseline migration
**Labels:** `setup`, `db`
- [ ] Install `golang-migrate` di Makefile
- [ ] `db/migrations/000001_init.up.sql` & `down.sql` (baseline kosong)
- [ ] Dokumentasi cara bikin migration baru

### #7 Setup Logrus structured logger
**Labels:** `setup`, `observability`
- [ ] JSON formatter di non-dev, text di dev
- [ ] Request-id middleware + inject ke context logger

### #8 Setup validator + global error handler Fiber
**Labels:** `setup`
- [ ] `go-playground/validator/v10` di-inject ke container
- [ ] Custom error response shape `{ "errors": [{field, message}] }`
- [ ] Fiber `ErrorHandler` translate panic + validation error

### #9 Setup GitHub Actions CI
**Labels:** `setup`, `ci`
- [ ] Job: `go vet`, `golangci-lint`, `go test ./...`
- [ ] Cache module
- [ ] Jalankan di PR ke `main`

### #10 Setup Sentry integration
**Labels:** `setup`, `observability`
- [ ] Init di `cmd/web/main.go` jika DSN ada
- [ ] Recover middleware Fiber report ke Sentry

---

## Milestone 1 — Auth & User Profile

### #11 Migration: users, refresh_tokens, audit_logs
**Labels:** `db`, `auth`
- [ ] SQL up/down sesuai ERD (`docs/erd.md`)
- [ ] Index sesuai tabel di ERD

### #12 Entity + Repository: User & RefreshToken
**Labels:** `auth`
- [ ] `entity/user.go`, `entity/refresh_token.go` (GORM tags)
- [ ] `repository/user_repository.go` (Create, FindByEmail, FindByID)
- [ ] `repository/refresh_token_repository.go`
- [ ] Unit test repository pakai testcontainers

### #13 AuthUsecase: register + email verification
**Labels:** `auth`, `feature`
- [ ] `usecase/auth_usecase.go` — `Register`, `VerifyEmail`
- [ ] Hash password bcrypt cost 12
- [ ] Generate verification token (crypto/rand 32 byte, disimpan hash)
- [ ] Enqueue email via Asynq
- [ ] Unit test dengan mock repo & mock email gateway

### #14 AuthUsecase: login + JWT issue
**Labels:** `auth`, `feature`
- [ ] `Login` return access (15m) + refresh (30d)
- [ ] Refresh token disimpan sebagai SHA-256 hash
- [ ] Rate limit per-email + per-IP di controller
- [ ] Audit log entry saat login sukses/gagal

### #15 AuthUsecase: refresh token rotation
**Labels:** `auth`
- [ ] Flow sesuai `docs/flows.md#7`
- [ ] Revoke lama saat issue baru

### #16 AuthUsecase: forgot & reset password
**Labels:** `auth`
- [ ] Token TTL 30 menit, single-use
- [ ] Kirim email reset
- [ ] Invalidate semua refresh token saat reset

### #17 AuthController + routes
**Labels:** `auth`, `http`
- [ ] `POST /api/v1/auth/register`
- [ ] `POST /api/v1/auth/login`
- [ ] `POST /api/v1/auth/refresh`
- [ ] `POST /api/v1/auth/verify-email`
- [ ] `POST /api/v1/auth/forgot-password`
- [ ] `POST /api/v1/auth/reset-password`
- [ ] Integration test happy path & error

### #18 JWT auth middleware
**Labels:** `auth`, `middleware`
- [ ] Parse Bearer, verify signature, set `ctx.userID`, `ctx.role`
- [ ] Redis blacklist check per-`jti`
- [ ] Helper `GetUserID(c)` untuk controller

### #19 Migration + entity + repo: creator_profiles
**Labels:** `db`, `creator`
- [ ] Reserved usernames list (admin, api, auth, login, dll.)
- [ ] Validasi regex username `^[a-z0-9_]{3,20}$`

### #20 Creator profile CRUD (me)
**Labels:** `creator`, `feature`
- [ ] `GET /api/v1/me/profile`, `PUT /api/v1/me/profile`
- [ ] Validasi unik username
- [ ] Whitelist field yang bisa diubah

### #21 Public profile endpoint
**Labels:** `creator`, `public`
- [ ] `GET /api/v1/creators/:username` — hanya field publik
- [ ] Cache di Redis 60 detik

### #22 Avatar upload ke S3-compatible storage
**Labels:** `creator`, `storage`
- [ ] `POST /api/v1/me/avatar` multipart
- [ ] Validasi MIME + size (max 2MB)
- [ ] Gateway `storage/s3_gateway.go` interface

---

## Milestone 2 — Donation Flow

### #23 Migration: donations, payments, webhook_events
**Labels:** `db`, `donation`

### #24 Entity + Repository: Donation, Payment, WebhookEvent
**Labels:** `donation`

### #25 Payment gateway interface + Xendit impl
**Labels:** `payment`, `gateway`
- [ ] `gateway/payment/payment_gateway.go` interface: `CreateInvoice(req) Response`
- [ ] `gateway/payment/xendit_gateway.go`
- [ ] Config: secret key, callback URL
- [ ] Test dengan httptest mock server

### #26 DonationUsecase: create donation
**Labels:** `donation`, `feature`
- [ ] Hitung fee (platform % + flat + PG fee estimate)
- [ ] Simpan donation+payment dalam 1 tx
- [ ] Panggil Xendit, simpan `checkout_url`
- [ ] Validasi amount vs min/max creator

### #27 DonationController
**Labels:** `donation`, `http`
- [ ] `POST /api/v1/donations`
- [ ] `GET /api/v1/donations/:id` (status check untuk donor)
- [ ] reCAPTCHA middleware
- [ ] Rate limit per-IP

### #28 Webhook handler + signature verification
**Labels:** `payment`, `webhook`
- [ ] `POST /api/v1/webhooks/xendit`
- [ ] Verify `x-callback-token`
- [ ] Idempotency via `webhook_events.idempotency_key`
- [ ] Flow sesuai `docs/flows.md#2a`

### #29 PaymentUsecase: state transition donation
**Labels:** `payment`
- [ ] Allowed: pending→paid, pending→expired, pending→failed, paid→refunded
- [ ] Tolak transisi invalid dengan error
- [ ] Unit test matrix transitions

### #30 Email receipt untuk donor + notif creator
**Labels:** `notification`
- [ ] Asynq job `EmailDonationReceipt`
- [ ] Asynq job `EmailCreatorNewDonation`
- [ ] Template sederhana

---

## Milestone 3 — Ledger & Wallet

### #31 Migration + entity: wallets, ledger_entries
**Labels:** `db`, `ledger`

### #32 LedgerUsecase: post donation entries
**Labels:** `ledger`, `money`
- [ ] Split sesuai `docs/erd.md`
- [ ] Invariant check (sum debit == sum credit) sebelum commit
- [ ] Unit test berbagai skenario fee

### #33 WalletUsecase: balance cache update
**Labels:** `ledger`, `money`
- [ ] Pakai `SELECT ... FOR UPDATE` untuk hindari race
- [ ] Atau `UPDATE wallets SET balance = balance + ?` atomic

### #34 Endpoint wallet & ledger history
**Labels:** `wallet`, `http`
- [ ] `GET /api/v1/me/wallet`
- [ ] `GET /api/v1/me/ledger?from&to&cursor`

### #35 Reconciliation job harian
**Labels:** `job`, `ledger`
- [ ] Asynq periodic task 02:00 WIB
- [ ] Flow sesuai `docs/flows.md#5`
- [ ] Discrepancy masuk audit_logs + alert ke Sentry

---

## Milestone 4 — Withdrawal

### #36 Migration + entity: bank_accounts, withdrawals
**Labels:** `db`, `withdrawal`

### #37 CRUD bank account
**Labels:** `withdrawal`
- [ ] `POST/GET/DELETE /api/v1/me/bank-accounts`
- [ ] Validasi bank code dari whitelist
- [ ] Set default

### #38 TOTP 2FA setup
**Labels:** `security`, `auth`
- [ ] `POST /api/v1/me/2fa/enable` generate secret + QR URL
- [ ] `POST /api/v1/me/2fa/verify` konfirmasi dengan code
- [ ] `POST /api/v1/me/2fa/disable`

### #39 WithdrawalUsecase: request
**Labels:** `withdrawal`, `money`
- [ ] Flow sesuai `docs/flows.md#4`
- [ ] Wajib TOTP jika enabled
- [ ] Hold funds via ledger (creator_balance → pending_withdrawal)

### #40 Disbursement gateway (Xendit)
**Labels:** `payment`, `gateway`
- [ ] `gateway/disbursement/xendit_disbursement.go`
- [ ] Interface bisa di-mock

### #41 WithdrawalUsecase: process disbursement (worker)
**Labels:** `withdrawal`, `worker`
- [ ] Asynq job dipicu setelah admin approve
- [ ] Idempotent (pakai `withdrawal.id` sebagai external ref)

### #42 Disbursement webhook handler
**Labels:** `webhook`, `withdrawal`
- [ ] `POST /api/v1/webhooks/disbursement`
- [ ] Success → settle ledger, failed → release funds

### #43 Admin approval endpoint
**Labels:** `admin`, `withdrawal`
- [ ] `GET /api/v1/admin/withdrawals?status=requested`
- [ ] `POST /api/v1/admin/withdrawals/:id/approve`
- [ ] `POST /api/v1/admin/withdrawals/:id/reject`
- [ ] Audit log

### #44 Email notifikasi withdrawal status
**Labels:** `notification`, `withdrawal`

---

## Milestone 5 — Streamer Alert (WebSocket)

### #45 Migration + entity: alert_tokens
**Labels:** `db`, `alert`

### #46 Generate & rotate alert token
**Labels:** `alert`
- [ ] `POST /api/v1/me/alert/token` (rotate)
- [ ] `GET /api/v1/me/alert/token`

### #47 Simpan preferensi alert
**Labels:** `alert`
- [ ] `PUT /api/v1/me/alert/settings` — durasi, template, suara, dll.

### #48 WS hub + handler
**Labels:** `alert`, `realtime`
- [ ] `GET /ws/alert?token=xxx` via Fiber websocket
- [ ] In-memory hub: map userID → []conn
- [ ] Graceful disconnect & unregister
- [ ] Flow sesuai `docs/flows.md#3`

### #49 Broadcast saat donation paid
**Labels:** `alert`, `realtime`
- [ ] Hook di PaymentUsecase (setelah commit ledger)
- [ ] Publish ke hub userID = creator.user_id

---

## Milestone 6 — Dashboard & Report

### #50 Dashboard summary endpoint
**Labels:** `dashboard`
- [ ] `GET /api/v1/me/dashboard?range=7d|30d|all`
- [ ] Total donasi, count, supporter unik, rata-rata

### #51 List donations dengan filter & pagination
**Labels:** `dashboard`
- [ ] `GET /api/v1/me/donations?status&from&to&cursor`
- [ ] Cursor-based pagination

### #52 Supporter leaderboard
**Labels:** `dashboard`
- [ ] `GET /api/v1/me/supporters?range`
- [ ] Group by donor_email / donor_name

### #53 Export CSV
**Labels:** `dashboard`
- [ ] `GET /api/v1/me/donations.csv`
- [ ] Stream response

---

## Milestone 7 — Admin & Compliance

### #54 Role-based access control middleware
**Labels:** `security`, `admin`
- [ ] `RequireRole("admin")` middleware
- [ ] Seed akun admin pertama via CLI command

### #55 Admin: user management
**Labels:** `admin`
- [ ] `GET /api/v1/admin/users`
- [ ] `POST /api/v1/admin/users/:id/suspend`
- [ ] `POST /api/v1/admin/users/:id/unsuspend`

### #56 Audit log viewer
**Labels:** `admin`, `observability`
- [ ] `GET /api/v1/admin/audit-logs?actor&action&from&to`

### #57 KYC upload sederhana
**Labels:** `admin`, `compliance`
- [ ] Upload KTP → S3
- [ ] Admin review & mark verified
- [ ] Blokir withdrawal > threshold jika belum verified

### #58 Rate limiting global + per-endpoint
**Labels:** `security`
- [ ] Middleware pakai Redis sliding window
- [ ] Config per-route (donate, login, webhook exempted)

### #59 reCAPTCHA middleware
**Labels:** `security`
- [ ] Google reCAPTCHA v3
- [ ] Pasang di register, login, donate

---

## Milestone 8 — Launch Readiness

### #60 OpenAPI spec + Swagger UI
**Labels:** `docs`
- [ ] `api/openapi.yaml`
- [ ] Serve di `/docs` (non-prod only)

### #61 Load test script
**Labels:** `perf`
- [ ] k6 script: donate flow (100 RPS), webhook, dashboard
- [ ] Target: p99 < 300ms di donate

### #62 Security checklist (OWASP top 10)
**Labels:** `security`
- [ ] Review SQL injection (GORM parameter), XSS (JSON only, no HTML render), CSRF (JWT header-based), dll.

### #63 Runbook operasional
**Labels:** `docs`, `ops`
- [ ] `docs/runbook.md`: incident, webhook replay, refund manual, admin onboarding

### #64 Monitoring & alerting
**Labels:** `observability`
- [ ] Prometheus metrics: req count, latency, error rate, payment success rate
- [ ] Healthcheck endpoint `/healthz`, `/readyz`
- [ ] Alert rule contoh

### #65 Privacy policy + ToS consent
**Labels:** `compliance`
- [ ] Field `agreed_terms_version` di users
- [ ] Endpoint ambil versi terbaru ToS

---

## Ringkasan Urutan Kerja

1. **Sprint 1:** #1–#10 (foundation siap)
2. **Sprint 2:** #11–#22 (user bisa register & punya profil publik)
3. **Sprint 3:** #23–#30 (donor bisa kirim uang, masuk ke DB)
4. **Sprint 4:** #31–#35 (ledger akurat, saldo bisa dipercaya)
5. **Sprint 5:** #36–#44 (creator bisa cair)
6. **Sprint 6:** #45–#49 (alert OBS jalan)
7. **Sprint 7:** #50–#53 (dashboard untuk creator)
8. **Sprint 8:** #54–#59 (admin panel + security hardening)
9. **Sprint 9:** #60–#65 (launch prep)

Total: **65 issue**. MVP fungsional (bisa terima donasi + cair) sudah tercapai setelah Sprint 5.
