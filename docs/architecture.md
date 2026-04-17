# Kreatip Backend — Architecture

Mengikuti pola [khannedy/golang-clean-architecture](https://github.com/khannedy/golang-clean-architecture), disesuaikan untuk kebutuhan platform donasi.

## Stack

| Layer | Tool |
|---|---|
| HTTP framework | **Fiber v2** |
| ORM | **GORM** (Postgres driver) |
| Config | **Viper** (`config.json` + env) |
| Migration | **golang-migrate** |
| Validation | **go-playground/validator** |
| Logging | **Logrus** (structured) |
| Async messaging | **Redis + Asynq** (Kafka opsional di v2) |
| Cache / rate limit | **Redis** (go-redis) |
| Realtime | **Fiber websocket/v2** |
| Payment | **Xendit Go SDK** (primary) |
| Email | **Resend / SMTP** |
| Test | `testify`, `testcontainers-go` |

## Project Structure

```
kreatip-backend/
├── api/                            # OpenAPI spec & Postman collection
├── cmd/
│   ├── web/main.go                 # HTTP server entrypoint (Fiber)
│   └── worker/main.go              # Async worker entrypoint (Asynq)
├── db/
│   └── migrations/                 # golang-migrate SQL files
├── internal/
│   ├── config/                     # viper loader, fiber bootstrap, db, redis, logger
│   │   ├── config.go
│   │   ├── database.go
│   │   ├── redis.go
│   │   ├── fiber.go
│   │   ├── logger.go
│   │   ├── validator.go
│   │   └── route.go
│   ├── entity/                     # GORM models (tabel DB)
│   │   ├── user.go
│   │   ├── creator_profile.go
│   │   ├── donation.go
│   │   ├── payment.go
│   │   ├── ledger_entry.go
│   │   ├── wallet.go
│   │   ├── withdrawal.go
│   │   ├── bank_account.go
│   │   ├── alert_token.go
│   │   ├── webhook_event.go
│   │   └── audit_log.go
│   ├── model/                      # Request/Response DTO
│   │   ├── user_model.go
│   │   ├── donation_model.go
│   │   ├── withdrawal_model.go
│   │   ├── wallet_model.go
│   │   ├── web_response.go
│   │   └── converter/              # entity <-> model converters
│   │       ├── user_converter.go
│   │       └── donation_converter.go
│   ├── repository/                 # Data access (GORM)
│   │   ├── repository.go           # generic base repo
│   │   ├── user_repository.go
│   │   ├── creator_profile_repository.go
│   │   ├── donation_repository.go
│   │   ├── payment_repository.go
│   │   ├── ledger_repository.go
│   │   ├── wallet_repository.go
│   │   ├── withdrawal_repository.go
│   │   └── webhook_event_repository.go
│   ├── usecase/                    # Business logic
│   │   ├── user_usecase.go
│   │   ├── auth_usecase.go
│   │   ├── creator_usecase.go
│   │   ├── donation_usecase.go
│   │   ├── payment_usecase.go      # webhook handling + state machine
│   │   ├── ledger_usecase.go
│   │   ├── wallet_usecase.go
│   │   ├── withdrawal_usecase.go
│   │   └── alert_usecase.go
│   ├── delivery/
│   │   ├── http/
│   │   │   ├── route/              # router registration
│   │   │   │   └── route.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go         # JWT
│   │   │   │   ├── rate_limit.go
│   │   │   │   ├── recaptcha.go
│   │   │   │   └── request_id.go
│   │   │   ├── user_controller.go
│   │   │   ├── auth_controller.go
│   │   │   ├── creator_controller.go
│   │   │   ├── donation_controller.go
│   │   │   ├── webhook_controller.go
│   │   │   ├── wallet_controller.go
│   │   │   ├── withdrawal_controller.go
│   │   │   └── admin_controller.go
│   │   ├── websocket/
│   │   │   └── alert_handler.go    # OBS overlay push
│   │   └── messaging/
│   │       ├── donation_consumer.go
│   │       └── email_consumer.go
│   └── gateway/                    # External systems
│       ├── payment/
│       │   ├── payment_gateway.go  # interface
│       │   └── xendit_gateway.go
│       ├── disbursement/
│       │   └── xendit_disbursement.go
│       ├── email/
│       │   └── resend_gateway.go
│       ├── storage/
│       │   └── s3_gateway.go
│       └── messaging/
│           └── asynq_producer.go
├── test/                           # Integration & unit tests
│   ├── user_test.go
│   ├── donation_test.go
│   └── helper_test.go
├── config.json                     # Viper config (dev)
├── .env.example
├── docker-compose.yml              # postgres + redis + mailhog
├── Dockerfile
├── Makefile
├── go.mod
└── go.sum
```

## Data Flow

```
HTTP Request
   │
   ▼
[Middleware] → auth, rate limit, request-id, recover
   │
   ▼
[Controller]  ── validate request → build Model
   │
   ▼
[Use Case]    ── orchestrate business logic, transactions
   │            │
   │            ├──► [Repository] ── GORM ──► PostgreSQL
   │            │
   │            ├──► [Gateway]    ── HTTP ──► Xendit / Resend / S3
   │            │
   │            └──► [Messaging]  ── Asynq ──► Redis queue → Worker
   │
   ▼
[Model] (response DTO) ── JSON ──► Client
```

## Layering Rules

- `entity` tidak boleh import package lain (pure struct + GORM tag).
- `model` tidak boleh import `entity` dari request side; converter package yang jembatani.
- `repository` menerima `*gorm.DB` via constructor, return `entity.*`.
- `usecase` adalah satu-satunya tempat yang boleh memulai DB transaction (`db.Transaction(...)`).
- `delivery/http` hanya parsing, validasi, panggil usecase, serialize response.
- `gateway` dibungkus interface supaya bisa di-mock untuk test.
- Semua duit (donation/withdrawal) wajib melewati `ledger_usecase` — saldo hanya cache dari ledger.

## Config (`config.json`)

```json
{
  "app": { "name": "kreatip", "port": 8080, "env": "development" },
  "database": {
    "host": "localhost", "port": 5432, "user": "postgres",
    "password": "postgres", "name": "kreatip",
    "pool": { "idle": 5, "max": 20, "lifetime_seconds": 300 }
  },
  "redis": { "addr": "localhost:6379", "db": 0 },
  "jwt": { "secret": "change-me", "access_ttl": "15m", "refresh_ttl": "720h" },
  "xendit": { "secret_key": "", "webhook_token": "", "callback_url": "" },
  "email": { "provider": "resend", "api_key": "", "from": "noreply@kreatip.id" },
  "storage": { "endpoint": "", "bucket": "", "access_key": "", "secret_key": "" },
  "fee": { "platform_percent": 5, "platform_flat": 500, "min_donation": 1000, "min_withdrawal": 50000 }
}
```
