# Kreatip Backend — ERD

Semua nominal uang disimpan sebagai **BIGINT dalam satuan cent/rupiah terkecil** (bukan float) untuk hindari floating-point error. Timestamps: `created_at`, `updated_at`, `deleted_at` (soft delete via GORM) diasumsikan ada di semua tabel kecuali disebut lain.

## Diagram

```mermaid
erDiagram
    USERS ||--o| CREATOR_PROFILES : has
    USERS ||--o{ BANK_ACCOUNTS : owns
    USERS ||--o{ WALLETS : owns
    USERS ||--o{ WITHDRAWALS : requests
    USERS ||--o{ ALERT_TOKENS : has
    USERS ||--o{ AUDIT_LOGS : generates
    USERS ||--o{ REFRESH_TOKENS : has

    CREATOR_PROFILES ||--o{ DONATIONS : receives
    DONATIONS ||--|| PAYMENTS : has
    DONATIONS ||--o{ LEDGER_ENTRIES : produces
    WITHDRAWALS ||--o{ LEDGER_ENTRIES : produces
    WITHDRAWALS }o--|| BANK_ACCOUNTS : paid_to

    WALLETS ||--o{ LEDGER_ENTRIES : references
    PAYMENTS ||--o{ WEBHOOK_EVENTS : triggers

    USERS {
        uuid id PK
        string email UK
        string password_hash
        string role "user|creator|admin"
        bool email_verified
        string totp_secret "nullable"
        bool totp_enabled
        timestamp created_at
        timestamp updated_at
    }

    CREATOR_PROFILES {
        uuid id PK
        uuid user_id FK UK
        string username UK
        string display_name
        string avatar_url
        text bio
        jsonb social_links
        string theme_color
        bigint min_donation "default 1000"
        bigint max_donation "default 10000000"
        string thank_you_message
        bool is_active
    }

    DONATIONS {
        uuid id PK
        uuid creator_id FK "refs creator_profiles"
        string donor_name "nullable, can be anonymous"
        string donor_email "nullable"
        bigint amount "gross"
        bigint platform_fee
        bigint payment_fee
        bigint net_amount "amount - fees"
        text message
        string status "pending|paid|failed|expired|refunded"
        timestamp paid_at
        string currency "IDR"
    }

    PAYMENTS {
        uuid id PK
        uuid donation_id FK UK
        string provider "xendit|midtrans"
        string method "qris|va_bca|gopay|ovo|..."
        string external_id "provider ref id"
        string external_status
        string checkout_url
        jsonb raw_request
        jsonb raw_response
        timestamp expires_at
    }

    WEBHOOK_EVENTS {
        uuid id PK
        string provider
        string event_type
        string external_id
        string idempotency_key UK
        jsonb payload
        bool processed
        string error
        timestamp received_at
        timestamp processed_at
    }

    LEDGER_ENTRIES {
        uuid id PK
        uuid wallet_id FK
        string account "creator_balance|platform_revenue|payment_gateway|pending_settlement"
        string direction "debit|credit"
        bigint amount
        string ref_type "donation|withdrawal|adjustment"
        uuid ref_id
        string description
        timestamp occurred_at
    }

    WALLETS {
        uuid id PK
        uuid user_id FK
        string currency "IDR"
        bigint balance_available "cache, source of truth = ledger"
        bigint balance_pending
        timestamp last_reconciled_at
    }

    BANK_ACCOUNTS {
        uuid id PK
        uuid user_id FK
        string type "bank|ewallet"
        string bank_code "BCA|BNI|GOPAY|..."
        string account_number
        string account_holder_name
        bool is_verified
        bool is_default
    }

    WITHDRAWALS {
        uuid id PK
        uuid user_id FK
        uuid bank_account_id FK
        bigint amount
        bigint fee
        bigint net_amount
        string status "requested|approved|processing|success|failed|rejected"
        string external_id "xendit disbursement id"
        string failure_reason
        uuid approved_by "admin user id, nullable"
        timestamp approved_at
        timestamp completed_at
    }

    ALERT_TOKENS {
        uuid id PK
        uuid user_id FK
        string token UK "opaque, rotatable"
        jsonb settings "duration, sound, template, animation"
        timestamp revoked_at
    }

    REFRESH_TOKENS {
        uuid id PK
        uuid user_id FK
        string token_hash UK
        string user_agent
        string ip
        timestamp expires_at
        timestamp revoked_at
    }

    AUDIT_LOGS {
        uuid id PK
        uuid actor_user_id FK
        string action "login|withdraw|admin_approve|..."
        string resource_type
        uuid resource_id
        jsonb metadata
        string ip
        timestamp occurred_at
    }
```

## Index & Constraint Penting

| Table | Index / Constraint |
|---|---|
| users | `UNIQUE(email)` |
| creator_profiles | `UNIQUE(username)`, `UNIQUE(user_id)` |
| donations | `INDEX(creator_id, status, created_at DESC)`, `INDEX(status, created_at)` |
| payments | `UNIQUE(donation_id)`, `INDEX(external_id)` |
| webhook_events | `UNIQUE(idempotency_key)`, `INDEX(provider, external_id)` |
| ledger_entries | `INDEX(wallet_id, occurred_at)`, `INDEX(ref_type, ref_id)` |
| withdrawals | `INDEX(user_id, status)`, `INDEX(status, created_at)` |
| refresh_tokens | `INDEX(user_id)`, `UNIQUE(token_hash)` |
| alert_tokens | `UNIQUE(token)` |

## Aturan Ledger (double-entry)

Saat **donasi PAID** nominal Rp 10.000 dengan platform fee 5% (Rp 500) + PG fee Rp 300:

| direction | account | amount |
|---|---|---|
| debit | payment_gateway | 10.000 |
| credit | creator_balance | 9.200 |
| credit | platform_revenue | 500 |
| credit | pg_fee_expense | 300 |

Saat **withdrawal SUCCESS** Rp 9.200:

| direction | account | amount |
|---|---|---|
| debit | creator_balance | 9.200 |
| credit | payment_gateway | 9.200 |

Invariant: `SUM(credit) == SUM(debit)` per transaksi. `wallets.balance_available` = `SUM(credit) - SUM(debit)` untuk `account = creator_balance` per user — di-cache, direkonsiliasi harian.

## State Machines

**Donation / Payment:**
```
pending ──► paid ──► (ledger posted)
   │          │
   │          └──► refunded (manual admin)
   ├──► expired
   └──► failed
```

**Withdrawal:**
```
requested ──► approved ──► processing ──► success
     │           │              │
     │           └──► rejected  └──► failed
     └──► cancelled (by user, sebelum approved)
```
