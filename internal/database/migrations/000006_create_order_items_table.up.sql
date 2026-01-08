-- Create order_items table with hybrid ID approach
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE SET NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL CHECK (unit_price >= 0),
    customizations JSONB,
    subtotal DECIMAL(10,2) NOT NULL CHECK (subtotal >= 0),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_order_items_uuid ON order_items(uuid);
CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product ON order_items(product_id);
CREATE INDEX IF NOT EXISTS idx_order_items_customizations ON order_items USING GIN (customizations);

-- Add comments
COMMENT ON TABLE order_items IS 'Order line items with product snapshot and customizations';
COMMENT ON COLUMN order_items.product_name IS 'Snapshot of product name at order time (for historical data)';
COMMENT ON COLUMN order_items.unit_price IS 'Base price of product at order time (before customizations)';
COMMENT ON COLUMN order_items.customizations IS 'JSONB storing selected customization options with prices';
COMMENT ON COLUMN order_items.subtotal IS 'Total item price: (unit_price + customization_prices) * quantity';
