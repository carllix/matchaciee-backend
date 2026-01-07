-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    midtrans_order_id VARCHAR(100) UNIQUE NOT NULL,
    payment_type VARCHAR(50),
    gross_amount DECIMAL(10,2) NOT NULL CHECK (gross_amount >= 0),
    transaction_status VARCHAR(50),
    transaction_id VARCHAR(100),
    transaction_time TIMESTAMP,
    settlement_time TIMESTAMP,
    fraud_status VARCHAR(50),
    status_message TEXT,
    payment_metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_payments_order ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_midtrans_order ON payments(midtrans_order_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(transaction_status);
CREATE INDEX IF NOT EXISTS idx_payments_transaction_id ON payments(transaction_id);
CREATE INDEX IF NOT EXISTS idx_payments_metadata ON payments USING GIN (payment_metadata);

-- Add comments
COMMENT ON TABLE payments IS 'Payment tracking table for Midtrans integration';
COMMENT ON COLUMN payments.midtrans_order_id IS 'Unique order ID sent to Midtrans (format: MC-250107-001-timestamp)';
COMMENT ON COLUMN payments.transaction_status IS 'Midtrans status: pending, settlement, expire, cancel, deny, refund';
COMMENT ON COLUMN payments.payment_metadata IS 'Full webhook payload from Midtrans for debugging/audit';
