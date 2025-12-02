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
    -- updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE catalog.order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    -- product name
    -- I could have used the product_id to get the product name from the catalog but I 
    -- decided to store by name in a case where the product is deleted
    name VARCHAR(255) NOT NULL,
    price DECIMAL(12,2) NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

    FOREIGN KEY (order_id) REFERENCES catalog.orders(id) ON DELETE CASCADE,
);