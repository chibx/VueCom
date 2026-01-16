-- Catalog Database
\c vuecom_catalog;

CREATE TABLE attributes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    -- e.g., 'Size', 'Color', 'Warranty Type'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attributes_name ON attributes (name);

-- For fast lookups by name
-- Attribute Options
CREATE TABLE category (
    id SERIAL PRIMARY KEY,
    attribute_id INTEGER NOT NULL REFERENCES attributes(id) ON DELETE CASCADE,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (attribute_id, value) -- Prevent duplicate values per attribute
);

CREATE INDEX idx_category_value ON category (value);

CREATE INDEX idx_category_attribute_id ON category (attribute_id);

CREATE TABLE presets (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    -- e.g., 'Clothing Preset'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE preset_attributes (
    preset_id INTEGER NOT NULL REFERENCES presets(id) ON DELETE CASCADE,
    attribute_id INTEGER NOT NULL REFERENCES attributes(id) ON DELETE CASCADE,
    PRIMARY KEY (preset_id, attribute_id)
);

CREATE INDEX idx_preset_attributes_preset_id ON preset_attributes (preset_id);

CREATE INDEX idx_preset_attributes_attribute_id ON preset_attributes (attribute_id);

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
    preset_id INTEGER REFERENCES presets(id) ON DELETE SET NULL,
    -- Optional
    FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE SET NULL,
    FOREIGN KEY (parent_id) REFERENCES products(id) ON DELETE SET NULL
);

CREATE INDEX idx_sku ON products (sku);
CREATE INDEX idx_slug ON products (slug);
CREATE INDEX idx_category ON products (category_id);
CREATE INDEX idx_active ON products (is_active);
CREATE INDEX ft_search ON products (name, short_description, search_keywords);

-- I sense a bug here
CREATE TABLE product_category_values (
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

CREATE INDEX idx_product_category_values_product_id ON product_category_values (product_id);

CREATE INDEX idx_product_category_values_category_id ON product_category_values (category_id);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
    -- e.g., 'Men's Shirts'
);

CREATE TABLE product_tags (
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, tag_id)
);

CREATE INDEX idx_product_tags_tag_id ON product_tags (tag_id);

-- Main Coupons Table
CREATE TABLE promo_codes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    TYPE TEXT NOT NULL CHECK (TYPE IN ('percentage', 'fixed_amount', 'fixed_product', 'bogo', 'free_shipping')),
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
    user_id INTEGER,
    order_id INTEGER,
    used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_promo_codes_code ON promo_codes(code);

CREATE INDEX IF NOT EXISTS idx_promo_code_usages_code_id ON promo_code_usages(code_id);

CREATE INDEX IF NOT EXISTS idx_promo_code_usages_user_id ON promo_code_usages(user_id);