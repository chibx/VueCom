-- Inventory Database
\c vuecom_inventory;

CREATE TABLE countries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(5) NOT NULL UNIQUE,
    -- e.g., 'US', 'NG'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_countries_code ON countries (code);

CREATE TABLE states (
    id SERIAL PRIMARY KEY,
    country_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_states_name ON states (name);

-- Warehouses (Lagos, Abuja, PH, etc.)
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) UNIQUE NOT NULL, -- e.g., LOS1, ABJ1
    name VARCHAR(100) NOT NULL,
    address TEXT,
    city VARCHAR(50),
    state INT,
    country INT,
    is_active BOOLEAN DEFAULT TRUE,
    capacity INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (state) REFERENCES states(id),
    FOREIGN KEY (country) REFERENCES countries(id)
);

CREATE INDEX idx_warehouse_code ON warehouses(code);


-- Inventory (Real-time Stock)
CREATE TABLE inventory (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(60) NOT NULL,
    product_id BIGINT NOT NULL, -- Reference to catalog
    warehouse_id INT NOT NULL,

    available_qty INT DEFAULT 0,
    reserved_qty INT DEFAULT 0,
    on_hold_qty INT DEFAULT 0,
    total_qty INT GENERATED ALWAYS AS (available_qty + reserved_qty + on_hold_qty) STORED,

    safety_stock INT DEFAULT 0,
    reorder_level INT DEFAULT 0,

    last_restocked_at TIMESTAMP NULL,
    last_sold_at TIMESTAMP NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE RESTRICT,

    UNIQUE (sku, warehouse_id)
);

CREATE INDEX idx_sku ON inventory(sku);
CREATE INDEX idx_product ON inventory(product_id);
CREATE INDEX idx_warehouse ON inventory(warehouse_id);
CREATE INDEX idx_available ON inventory(available_qty);
CREATE INDEX idx_updated ON inventory(updated_at);


-- Stock Movement Log (Audit Trail)
CREATE TABLE stock_movements (
    id BIGSERIAL PRIMARY KEY,
    inventory_id BIGINT NOT NULL,
    sku VARCHAR(60) NOT NULL,
    warehouse_id INT NOT NULL,
    -- movement_type ENUM('RESTOCK', 'SALE', 'RETURN', 'ADJUSTMENT', 'TRANSFER') NOT NULL,
    movement_type VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    reference VARCHAR(100), -- e.g., order_id, transfer_id
    notes TEXT,
    created_by VARCHAR(50), -- admin, system, supplier
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (inventory_id) REFERENCES inventory(id) ON DELETE RESTRICT
);

CREATE INDEX idx_sku_stock_movements ON stock_movements(sku);
CREATE INDEX idx_type_stock_movements ON stock_movements(movement_type);
CREATE INDEX idx_date_stock_movements ON stock_movements(created_at);