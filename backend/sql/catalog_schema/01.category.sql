CREATE SCHEMA IF NOT EXISTS catalog;

-- Attributes
CREATE TABLE catalog.attributes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,  -- e.g., 'Size', 'Color', 'Warranty Type'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_attributes_name ON catalog.attributes (name);  -- For fast lookups by name

-- Attribute Options
CREATE TABLE catalog.category (
    id SERIAL PRIMARY KEY,
    attribute_id INTEGER NOT NULL REFERENCES catalog.attributes(id) ON DELETE CASCADE,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (attribute_id, value)  -- Prevent duplicate values per attribute
);

CREATE INDEX idx_category_value ON catalog.category (value);
CREATE INDEX idx_category_attribute_id ON catalog.category (attribute_id);


CREATE TABLE catalog.presets (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,  -- e.g., 'Clothing Preset'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE catalog.preset_attributes (
    preset_id INTEGER NOT NULL REFERENCES catalog.presets(id) ON DELETE CASCADE,
    attribute_id INTEGER NOT NULL REFERENCES catalog.attributes(id) ON DELETE CASCADE,
    PRIMARY KEY (preset_id, attribute_id)
);

CREATE INDEX idx_preset_attributes_preset_id ON catalog.preset_attributes (preset_id);
CREATE INDEX idx_preset_attributes_attribute_id ON catalog.preset_attributes (attribute_id);

-- I sense a bug here
CREATE TABLE catalog.product_category_values (
    product_id INTEGER NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES catalog.category(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

CREATE INDEX idx_product_attribute_values_product_id ON catalog.product_attribute_values (product_id);
CREATE INDEX idx_product_attribute_values_category_id ON catalog.product_attribute_values (category_id);

CREATE TABLE catalog.tags (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,  -- e.g., 'Men's Shirts'
);

CREATE TABLE catalog.product_tags (
    product_id INTEGER NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES catalog.tags(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, tag_id)
);

CREATE INDEX idx_product_tags_tag_id ON catalog.product_tags (tag_id);