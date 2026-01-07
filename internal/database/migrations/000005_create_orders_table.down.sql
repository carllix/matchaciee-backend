-- Drop indexes
DROP INDEX IF EXISTS idx_orders_queue;
DROP INDEX IF EXISTS idx_orders_created;
DROP INDEX IF EXISTS idx_orders_source;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_user;
DROP INDEX IF EXISTS idx_orders_number;

-- Drop orders table
DROP TABLE IF EXISTS orders CASCADE;
