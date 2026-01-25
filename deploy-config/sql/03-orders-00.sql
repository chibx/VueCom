-- Orders Database
\c vuecom_orders;
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    -- indirect
    user_id INTEGER NOT NULL,
    order_number TEXT NOT NULL,
    total_amount DECIMAL(12, 2) NOT NULL,
    -- currency ENUM('NGN', 'USD') DEFAULT 'NGN',
    currency TEXT DEFAULT 'NGN',
    -- status ENUM('pending', 'completed', 'cancelled') DEFAULT 'pending',
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- indirect
    billing_address_id INTEGER NOT NULL,
    -- indirect
    shipping_address_id INTEGER NOT NULL,
    -- indirect
    payment_id INTEGER NOT NULL
);

CREATE INDEX orders_user_idx ON orders(user_id);
CREATE INDEX orders_status_idx ON orders(status);
CREATE INDEX orders_total_amount_idx ON orders(total_amount);
CREATE INDEX orders_created_at_idx ON orders(created_at);
CREATE INDEX orders_updated_at_idx ON orders(updated_at);

CREATE TABLE order_returns (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE INDEX order_return_id_idx ON order_returns(order_id);
CREATE INDEX order_return_created_at_idx ON order_returns(created_at);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    -- product name
    -- I could have used the product_id to get the product name from the catalog but I
    -- decided to store by name in a case where the product is deleted
    name TEXT NOT NULL,
    price DECIMAL(12, 2) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE INDEX order_items_created_at_idx ON order_items(created_at);


\c postgres;
