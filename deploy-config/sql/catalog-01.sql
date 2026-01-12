-- Catalog Database
\c catalog

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

-- I sense a bug here
CREATE TABLE product_category_values (
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

CREATE INDEX idx_product_attribute_values_product_id ON product_attribute_values (product_id);

CREATE INDEX idx_product_attribute_values_category_id ON product_attribute_values (category_id);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    -- e.g., 'Men's Shirts'
);

CREATE TABLE product_tags (
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, tag_id)
);

CREATE INDEX idx_product_tags_tag_id ON product_tags (tag_id);