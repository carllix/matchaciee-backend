-- Drop indexes
DROP INDEX IF EXISTS idx_categories_display_order;
DROP INDEX IF EXISTS idx_categories_active;
DROP INDEX IF EXISTS idx_categories_slug;

-- Drop categories table
DROP TABLE IF EXISTS categories CASCADE;
