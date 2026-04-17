-- =============================================================
-- 000001_init_schema.down.sql
-- Drop all tables in reverse FK dependency order.
-- =============================================================

DROP TABLE IF EXISTS audit_logs       CASCADE;
DROP TABLE IF EXISTS webhook_events   CASCADE;
DROP TABLE IF EXISTS withdrawals      CASCADE;
DROP TABLE IF EXISTS ledger_entries   CASCADE;
DROP TABLE IF EXISTS payments         CASCADE;
DROP TABLE IF EXISTS donations        CASCADE;
DROP TABLE IF EXISTS bank_accounts    CASCADE;
DROP TABLE IF EXISTS alert_tokens     CASCADE;
DROP TABLE IF EXISTS refresh_tokens   CASCADE;
DROP TABLE IF EXISTS wallets          CASCADE;
DROP TABLE IF EXISTS creator_profiles CASCADE;
DROP TABLE IF EXISTS users            CASCADE;
