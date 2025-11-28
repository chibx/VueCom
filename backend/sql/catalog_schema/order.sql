CREATE SCHEMA IF NOT EXISTS catalog;

CREATE TABLE catalog.orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    order_number VARCHAR(50) NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    currency ENUM('NGN', 'USD') DEFAULT 'NGN',
    status ENUM('pending', 'completed', 'cancelled') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    billing_address_id INTEGER NOT NULL,
    shipping_address_id INTEGER NOT NULL,
    payment_id INTEGER NOT NULL,
);


CREATE TABLE catalog.order_returns (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    reason VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);