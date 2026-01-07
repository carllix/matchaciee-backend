-- Drop indexes
DROP INDEX IF EXISTS idx_customizations_type;
DROP INDEX IF EXISTS idx_customizations_product;

-- Drop product_customizations table
DROP TABLE IF EXISTS product_customizations CASCADE;
