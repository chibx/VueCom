CREATE SCHEMA IF NOT EXISTS catalog;

-- Products (Metadata ONLY)
CREATE TABLE catalog.products (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    short_description TEXT,
    full_description TEXT,
    category_id INT,
    base_price DECIMAL(12,2) NOT NULL,
    sale_price DECIMAL(12,2),
    weight_grams INT,
    dimensions VARCHAR(50), -- e.g., "10x5x3 cm"
    is_active BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,
    meta_title VARCHAR(255),
    meta_description TEXT,
    search_keywords TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    parent_id INT NULL,
    image_url VARCHAR(255),
    preset_id INTEGER REFERENCES catalog.presets(id) ON DELETE SET NULL,  -- Optional

    FOREIGN KEY (category_id) REFERENCES catalog.categories(id) ON DELETE SET NULL,
    FOREIGN KEY (parent_id) REFERENCES catalog.products(id) ON DELETE SET NULL,

    INDEX idx_sku (sku),
    INDEX idx_slug (slug),
    INDEX idx_category (category_id),
    INDEX idx_brand (brand_id),
    INDEX idx_active (is_active),
    INDEX ft_search (name, short_description, search_keywords)
);

-- -- Categories
-- CREATE TABLE categories (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(100) NOT NULL,
--     slug VARCHAR(100) UNIQUE NOT NULL,
--     parent_id INT NULL,
--     image_url VARCHAR(255),
--     is_active BOOLEAN DEFAULT TRUE,
--     FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );