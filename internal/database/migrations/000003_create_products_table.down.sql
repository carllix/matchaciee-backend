-- Drop indexes
DROP INDEX IF EXISTS idx_products_display_order;
DROP INDEX IF EXISTS idx_products_deleted;
DROP INDEX IF EXISTS idx_products_available;
DROP INDEX IF EXISTS idx_products_slug;
DROP INDEX IF EXISTS idx_products_category;

-- Drop products table
DROP TABLE IF EXISTS products CASCADE;
