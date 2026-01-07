-- Create products table with hybrid ID approach
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    base_price DECIMAL(10,2) NOT NULL CHECK (base_price >= 0),
    image_url VARCHAR(255),
    is_available BOOLEAN DEFAULT true,
    is_customizable BOOLEAN DEFAULT false,
    preparation_time INT DEFAULT 5 CHECK (preparation_time >= 0),
    display_order INT DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_products_uuid ON products(uuid);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug);
CREATE INDEX IF NOT EXISTS idx_products_available ON products(is_available);
CREATE INDEX IF NOT EXISTS idx_products_deleted ON products(deleted_at);
CREATE INDEX IF NOT EXISTS idx_products_display_order ON products(display_order);
