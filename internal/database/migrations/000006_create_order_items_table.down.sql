-- Drop indexes
DROP INDEX IF EXISTS idx_order_items_customizations;
DROP INDEX IF EXISTS idx_order_items_product;
DROP INDEX IF EXISTS idx_order_items_order;

-- Drop order_items table
DROP TABLE IF EXISTS order_items CASCADE;
