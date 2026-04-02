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

CREATE TABLE media_folders (
    id          SERIAL PRIMARY KEY,
    parent_id   INTEGER REFERENCES media_folders(id) ON DELETE CASCADE,  -- allows deleting whole folder trees
    name        TEXT NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Prevent duplicate folder names inside the same parent (like a real file system)
ALTER TABLE media_folders ADD CONSTRAINT unique_folder_name_per_parent UNIQUE (parent_id, name);

CREATE INDEX idx_media_folders_parent ON media_folders(parent_id);

CREATE TABLE medias (
    id BIGSERIAL PRIMARY KEY,
    folder_id INTEGER REFERENCES media_folders(id) ON DELETE CASCADE,  -- move = just change this
    type TEXT NOT NULL,
    url TEXT NOT NULL,
    thumbnail_url   TEXT, -- generated thumbnail (for images)
    original_name   TEXT NOT NULL, -- what user uploaded (e.g. "product-photo.jpg")
    external_id     TEXT NOT NULL UNIQUE, -- immutable UUID.ext on disk/S3 (never changes)
    mime_type       TEXT NOT NULL,
    size_bytes      INTEGER,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_medias_folder_id ON medias(folder_id);
CREATE INDEX idx_medias_type ON medias(type);
CREATE INDEX idx_medias_external_id ON medias(external_id);

-- Products (Metadata ONLY)
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    base_price DECIMAL(15, 2) NOT NULL,
    sale_price DECIMAL(15, 2),
    discount_start TIMESTAMP,
    discount_end TIMESTAMP,
    is_new BOOLEAN DEFAULT TRUE,
    new_from TIMESTAMP,
    new_to TIMESTAMP,
    brand_id INTEGER,
    color_id INTEGER,
    slug TEXT UNIQUE NOT NULL,
    short_description TEXT DEFAULT '',
    full_description TEXT DEFAULT '',
    country_of_manufacture INTEGER, -- Country table in the backend database
    weight DECIMAL(10, 3) DEFAULT 0.00,
    enabled BOOLEAN DEFAULT TRUE,
    meta_title TEXT NOT NULL,
    meta_description TEXT NOT NULL,
    search_keywords TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    parent_id INTEGER,
    preset_id INTEGER,

    FOREIGN KEY (preset_id) REFERENCES presets(id) ON DELETE SET NULL,
    FOREIGN KEY (parent_id) REFERENCES products(id) ON DELETE SET NULL
);

CREATE INDEX idx_sku ON products (sku);
CREATE INDEX idx_slug ON products (slug);


CREATE TYPE product_relation_type AS ENUM ('related', 'upsell', 'cross_sell');
CREATE TABLE product_relations (
    source_product_id  BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    target_product_id  BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    relation_type      product_relation_type NOT NULL,
    sort_order         INTEGER DEFAULT 0,
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_relation UNIQUE (source_product_id, target_product_id, relation_type)
);

CREATE INDEX idx_product_relations_source_type ON product_relations(source_product_id, relation_type);
CREATE INDEX idx_product_relations_target ON product_relations(target_product_id);

-- This is what makes the popup "select image/video for product"
CREATE TABLE product_medias (
    product_id  BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,  -- product gone → links gone
    media_id    BIGINT NOT NULL REFERENCES medias(id) ON DELETE CASCADE,    -- media deleted → removed from all products
    sort_order  INTEGER DEFAULT 0,
    is_main     BOOLEAN DEFAULT false
);

-- One media can only be added once per product
CREATE UNIQUE INDEX idx_product_media_unique ON product_medias(product_id, media_id);
-- Fast lookups for product gallery
CREATE INDEX idx_product_media_product ON product_medias(product_id);
CREATE INDEX idx_product_media_sort ON product_medias(product_id, sort_order);

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
