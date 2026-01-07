-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(20) UNIQUE NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    customer_name VARCHAR(255) NOT NULL,
    order_source VARCHAR(20) NOT NULL CHECK (order_source IN ('guest', 'member', 'kiosk')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'preparing', 'ready', 'completed', 'cancelled')),
    subtotal DECIMAL(10,2) NOT NULL CHECK (subtotal >= 0),
    tax DECIMAL(10,2) DEFAULT 0 CHECK (tax >= 0),
    total DECIMAL(10,2) NOT NULL CHECK (total >= 0),
    queue_number INT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_orders_number ON orders(order_number);
CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_source ON orders(order_source);
CREATE INDEX IF NOT EXISTS idx_orders_created ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_orders_queue ON orders(queue_number);

-- Add comments
COMMENT ON TABLE orders IS 'Main orders table for all order types (guest, member, kiosk)';
COMMENT ON COLUMN orders.order_number IS 'Unique order number format: MC-YYMMDD-XXX';
COMMENT ON COLUMN orders.user_id IS 'NULL for guest/kiosk orders, filled for member orders';
COMMENT ON COLUMN orders.customer_name IS 'Customer name for pickup (can differ from user.full_name)';
COMMENT ON COLUMN orders.queue_number IS 'Daily auto-increment queue number for pickup station display';
