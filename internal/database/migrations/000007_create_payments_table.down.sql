-- Drop indexes
DROP INDEX IF EXISTS idx_payments_metadata;
DROP INDEX IF EXISTS idx_payments_transaction_id;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_midtrans_order;
DROP INDEX IF EXISTS idx_payments_order;

-- Drop payments table
DROP TABLE IF EXISTS payments CASCADE;
