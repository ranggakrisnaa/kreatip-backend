# Kreatip Backend Docs

Platform donasi untuk creator & streamer (mirip Saweria). Backend-only, Go + Fiber + Postgres.

## Dokumen

| File | Isi |
|---|---|
| [architecture.md](./architecture.md) | Stack, project structure (clean architecture ala khannedy), layering rules, config |
| [erd.md](./erd.md) | Skema database + Mermaid ERD + aturan ledger + state machine |
| [flows.md](./flows.md) | Sequence diagram untuk 7 flow utama (register, donation, webhook, alert, withdrawal, reconcile, auth) |
| [issues.md](./issues.md) | 65 issue siap commit, dipecah per milestone/sprint |

## Referensi arsitektur

https://github.com/khannedy/golang-clean-architecture

## Keputusan awal

- **Bahasa/framework:** Go + Fiber v2
- **ORM:** GORM + Postgres
- **Payment:** Xendit (primary)
- **Arsitektur:** Clean architecture 6 layer (entity, model, repository, usecase, delivery, gateway)
- **Fokus rilis pertama:** backend-only, API siap dipakai oleh frontend/OBS overlay

## Yang masih perlu kamu putuskan

Lihat bagian "Keputusan yang perlu kamu ambil" — hosting, platform fee %, minimum donasi/withdrawal, brand name final.
