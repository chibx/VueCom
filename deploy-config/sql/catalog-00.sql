-- Catalog Database
\c catalog;
-- Products (Metadata ONLY)
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    short_description TEXT,
    full_description TEXT,
    category_id INT,
    base_price DECIMAL(12, 2) NOT NULL,
    sale_price DECIMAL(12, 2),
    weight_grams INT,
    dimensions VARCHAR(50),
    -- e.g., "10x5x3 cm"
    is_active BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,
    meta_title VARCHAR(255),
    meta_description TEXT,
    search_keywords TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    parent_id INT NULL,
    image_url VARCHAR(255),
    preset_id INTEGER REFERENCES presets(id) ON DELETE
    SET
        NULL,
        -- Optional
        FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE
    SET
        NULL,
        FOREIGN KEY (parent_id) REFERENCES products(id) ON DELETE
    SET
        NULL,
        INDEX idx_sku (sku),
        INDEX idx_slug (slug),
        INDEX idx_category (category_id),
        INDEX idx_brand (brand_id),
        INDEX idx_active (is_active),
        INDEX ft_search (name, short_description, search_keywords)
);

-- Main Coupons Table
CREATE TABLE promo_codes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    TYPE TEXT NOT NULL CHECK (
        TYPE IN (
            'percentage',
            'fixed_amount',
            'fixed_product',
            'bogo',
            'free_shipping'
        )
    ),
    -- Or use ENUM if preferred: CREATE TYPE coupon_type AS ENUM (...);
    discount_value DECIMAL(10, 2) NOT NULL,
    min_cart_value DECIMAL(10, 2) DEFAULT 0.00,
    expiry_date TIMESTAMP WITH TIME ZONE,
    start_date TIMESTAMP WITH TIME ZONE,
    usage_limit INTEGER,
    usage_limit_per_user INTEGER DEFAULT 1,
    product_ids JSONB,
    -- e.g., [1, 2, 3] for specific products
    category_ids JSONB,
    -- e.g., [10, 20] for categories
    exclude_product_ids JSONB,
    exclude_category_ids JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- -- Trigger for updating updated_at
-- CREATE OR REPLACE FUNCTION update_updated_at()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     NEW.updated_at = CURRENT_TIMESTAMP;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;
-- CREATE TRIGGER trigger_promo_codes_updated_at
-- BEFORE UPDATE ON promo_codes
-- FOR EACH ROW
-- EXECUTE FUNCTION update_updated_at();
-- Coupon Usages Table (for tracking redemptions)
CREATE TABLE promo_code_usages (
    id SERIAL PRIMARY KEY,
    code_id INTEGER NOT NULL REFERENCES promo_codes(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_promo_codes_code ON promo_codes(code);

CREATE INDEX IF NOT EXISTS idx_promo_code_usages_code_id ON promo_code_usages(code_id);

CREATE INDEX IF NOT EXISTS idx_promo_code_usages_user_id ON promo_code_usages(user_id);