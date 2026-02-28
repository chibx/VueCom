-- Users Database
\c vuecom_users;
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone_number TEXT NOT NULL,
    country_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_email_verified BOOLEAN DEFAULT FALSE,
    image_url TEXT,
    password_hash TEXT,
    FOREIGN KEY (country_id) REFERENCES countries(id)
);

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

CREATE INDEX IF NOT EXISTS customer_address_id_idx ON customer_addresses (customer_id);

CREATE TABLE customer_sessions (
    customer_id INT NOT NULL,
    refresh_token_hash TEXT NOT NULL,
    device_id UUID NOT NULL,
    last_ip TEXT NOT NULL,
    user_agent TEXT,
    used BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

-- CREATE INDEX idx_backend_sessions_refresh_token ON backend_sessions(refresh_token_hash);
CREATE INDEX IF NOT EXISTS idx_customer_sessions_user_id ON customer_sessions(customer_id);


CREATE TABLE customer_otps (
    customer_id INT NOT NULL,
    code VARCHAR(10) NOT NULL,
    expiry_date TIMESTAMP NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    UNIQUE(customer_id, code)
);

CREATE INDEX IF NOT EXISTS customer_otps_code_idx ON customer_otps (code);
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
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    UNIQUE (customer_id, product_id)
);

CREATE INDEX IF NOT EXISTS customer_cart_id_idx ON customer_carts (customer_id);
CREATE INDEX IF NOT EXISTS customer_cart_product_idx ON customer_carts (product_id);


\c postgres;
