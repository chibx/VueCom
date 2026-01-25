-- Users Database
\c vuecom_users;
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone_number VARCHAR(20) NOT NULL,
    country_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_email_verified BOOLEAN DEFAULT FALSE,
    image_url TEXT,
    password_hash TEXT,
    FOREIGN KEY (country) REFERENCES countries(id)
);

CREATE INDEX IF NOT EXISTS customers_email_idx ON customers USING hash (email);

CREATE TABLE customer_addresses (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    street_address TEXT NOT NULL,
    zip_code TEXT,
    state_id INT,
    country_id INT,
    FOREIGN KEY (country_id) REFERENCES countries(id),
    FOREIGN KEY (state_id) REFERENCES states(id)
);

CREATE INDEX IF NOT EXISTS customer_address_id_idx ON customer_addresses USING hash (customer_id);


-- CREATE TABLE customer_sessions (
--     id SERIAL PRIMARY KEY,
--     customer_id INT NOT NULL,
--     token VARCHAR(255) NOT NULL UNIQUE,
--     ip_address VARCHAR(45),
--     user_agent TEXT,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     expires_at TIMESTAMP NOT NULL,
--     FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
-- );

-- CREATE INDEX IF NOT EXISTS customer_sessions_id_idx ON customer_sessions USING hash (customer_id);
-- CREATE INDEX IF NOT EXISTS customer_sessions_expires_at_idx ON customer_sessions (expires_at);

CREATE TABLE customer_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id INT NOT NULL,
    refresh_token_hash TEXT NOT NULL,
    device_id UUID NOT NULL,
    last_ip TEXT NOT NULL,
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES customers(id) ON DELETE CASCADE
);

-- CREATE INDEX idx_backend_sessions_refresh_token ON backend_sessions(refresh_token_hash);
CREATE INDEX IF NOT EXISTS idx_customer_sessions_user_id ON customer_sessions(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_sessions_device ON customer_sessions(device_id);



CREATE TABLE customer_otps (
    customer_id INT NOT NULL,
    code VARCHAR(10) NOT NULL,
    expiry_date TIMESTAMP NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    UNIQUE(customer_id, code)
);

CREATE INDEX IF NOT EXISTS customer_otps_expiry_idx ON customer_otps (expiry_date);

CREATE TABLE customer_wishlists (
    customer_id INT NOT NULL,
    product_id BIGINT NOT NULL,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    UNIQUE (customer_id, product_id)
);

CREATE INDEX IF NOT EXISTS customer_wishlists_id_idx ON customer_wishlists (customer_id);
CREATE INDEX IF NOT EXISTS customer_wishlists_product_idx ON customer_wishlists (product_id);

CREATE TABLE customer_carts (
    customer_id INT NOT NULL,
    -- indirect
    product_id BIGINT NOT NULL,
    quantity INT DEFAULT 1,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    UNIQUE (customer_id, product_id)
);

CREATE INDEX IF NOT EXISTS customer_cart_id_idx ON customer_cart (customer_id);
CREATE INDEX IF NOT EXISTS customer_cart_product_idx ON customer_cart (product_id);


\c postgres;
