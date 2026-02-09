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
    attribute_id INTEGER NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Prevent duplicate values per attribute
    UNIQUE (attribute_id, value),

    FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE
);

CREATE INDEX idx_category_value ON category (value);
CREATE INDEX idx_category_attribute_id ON category (attribute_id);

CREATE TABLE presets (
    id SERIAL PRIMARY KEY,
    -- e.g., 'Clothing Preset'
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_preset_name ON presets (name);

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
    sku TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    base_price DECIMAL(15, 2) NOT NULL,
    sale_price DECIMAL(15, 2),
    discount_period TIMESTAMP,
    slug TEXT UNIQUE NOT NULL,
    short_description TEXT DEFAULT '',
    full_description TEXT DEFAULT '',
    weight DECIMAL(10, 3) DEFAULT 0.00,
    enabled BOOLEAN DEFAULT TRUE,
    -- is_featured BOOLEAN DEFAULT FALSE,
    meta_title TEXT,
    meta_description TEXT,
    search_keywords TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    parent_id INT,
    image_url TEXT,
    preset_id INTEGER,

    FOREIGN KEY (preset_id) REFERENCES presets(id) ON DELETE SET NULL,
    FOREIGN KEY (parent_id) REFERENCES products(id) ON DELETE SET NULL
);

CREATE INDEX idx_sku ON products (sku);
CREATE INDEX idx_slug ON products (slug);

-- I sense a bug here
CREATE TABLE product_category_values (
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

-- CREATE INDEX idx_product_category_values_product_id ON product_category_values (product_id);
CREATE INDEX idx_product_category_values_category_id ON product_category_values (category_id);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    -- e.g., 'Men's Shirts'
    name TEXT NOT NULL
);

CREATE TABLE product_tags (
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, tag_id)
);

CREATE INDEX idx_product_tags_tag_id ON product_tags (tag_id);


CREATE TYPE promo_code_type AS ENUM ('percentage', 'fixed_amount', 'free_shipping');
-- Main Coupons Table
CREATE TABLE promo_codes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    code TEXT UNIQUE NOT NULL,
    type promo_code_type NOT NULL,
    discount_value DECIMAL(10, 2) NOT NULL,
    min_cart_value DECIMAL(10, 2) DEFAULT 0.00,
    start_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expiry_date TIMESTAMP WITH TIME ZONE,
    usage_limit INTEGER,
    usage_limit_per_user INTEGER DEFAULT 1,
    product_ids INTEGER[],
    -- e.g., [1, 2, 3] for specific products
    category_ids INTEGER[],
    -- e.g., [10, 20] for categories
    exclude_product_ids INTEGER[],
    exclude_category_ids INTEGER[],
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_promo_codes_code ON promo_codes USING hash(code);

CREATE TABLE promo_code_usages (
    id SERIAL PRIMARY KEY,
    code_id INTEGER NOT NULL,
    customer_id INTEGER NOT NULL,
    order_id INTEGER NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (code_id) REFERENCES promo_codes(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_promo_code_usages_code_id ON promo_code_usages(code_id);
CREATE INDEX IF NOT EXISTS idx_promo_code_usages_customer_id ON promo_code_usages(customer_id);
CREATE INDEX IF NOT EXISTS idx_promo_code_usages_order_id ON promo_code_usages(order_id);
CREATE INDEX IF NOT EXISTS idx_promo_code_usages_used_at ON promo_code_usages(used_at);

\c postgres;
