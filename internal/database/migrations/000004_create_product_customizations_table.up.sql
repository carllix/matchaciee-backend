-- Create product_customizations table with hybrid ID approach
CREATE TABLE IF NOT EXISTS product_customizations (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    customization_type VARCHAR(50) NOT NULL,
    option_name VARCHAR(100) NOT NULL,
    price_modifier DECIMAL(10,2) DEFAULT 0 CHECK (price_modifier >= -999999.99 AND price_modifier <= 999999.99),
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_customizations_uuid ON product_customizations(uuid);
CREATE INDEX IF NOT EXISTS idx_customizations_product ON product_customizations(product_id);
CREATE INDEX IF NOT EXISTS idx_customizations_type ON product_customizations(customization_type);

-- Add comment to explain the table
COMMENT ON TABLE product_customizations IS 'Stores customization options for products (e.g., Matcha Level: Strong, Milk Type: Oat Milk)';
COMMENT ON COLUMN product_customizations.customization_type IS 'Category of customization (e.g., Matcha Level, Milk Type, Sugar Level)';
COMMENT ON COLUMN product_customizations.option_name IS 'Specific option within the category (e.g., Strong, Oat Milk, Less Sugar)';
COMMENT ON COLUMN product_customizations.price_modifier IS 'Price adjustment for this option (can be positive or negative)';
