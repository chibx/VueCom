CREATE SCHEMA IF NOT EXISTS inventory;
-- Database: inventory_service
-- CREATE DATABASE inventory_service;

-- Warehouses (Lagos, Abuja, PH, etc.)
CREATE TABLE inventory.warehouses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) UNIQUE NOT NULL, -- e.g., LOS1, ABJ1
    name VARCHAR(100) NOT NULL,
    address TEXT,
    city VARCHAR(50),
    state INT,
    country INT,
    is_active BOOLEAN DEFAULT TRUE,
    capacity INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

    INDEX idx_warehouse_code (code),
    FOREIGN KEY (state) REFERENCES backend.states(id),
    FOREIGN KEY (country) REFERENCES backend.countries(id),
);


-- Inventory (Real-time Stock)
CREATE TABLE inventory.inventory (
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

    FOREIGN KEY (warehouse_id) REFERENCES inventory.warehouses(id) ON DELETE RESTRICT,

    UNIQUE (sku, warehouse_id),
    INDEX idx_sku (sku),
    INDEX idx_product (product_id),
    INDEX idx_warehouse (warehouse_id),
    INDEX idx_available (available_qty),
    INDEX idx_updated (updated_at)
);


-- Stock Movement Log (Audit Trail)
CREATE TABLE inventory.stock_movements (
    id BIGERIAL PRIMARY KEY,
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

    FOREIGN KEY (inventory_id) REFERENCES inventory.inventory(id) ON DELETE RESTRICT,
    INDEX idx_sku (sku),
    INDEX idx_type (movement_type),
    INDEX idx_date (created_at)
);