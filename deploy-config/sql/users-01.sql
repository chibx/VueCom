-- Users Database
\c vuecom_users;
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(50) NOT NULL,
    state VARCHAR(50) NOT NULL,
    zip_code VARCHAR(20) NOT NULL,
    country INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_email_verified BOOLEAN DEFAULT FALSE,
    image_url TEXT,
    user_name VARCHAR(50) UNIQUE,
    password_hash TEXT,
    FOREIGN KEY (country) REFERENCES countries(id)
);

CREATE TABLE customer_sessions (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

CREATE TABLE customer_otps (
    customer_id INT NOT NULL,
    code VARCHAR(10) NOT NULL,
    expiry_date TIMESTAMP NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

CREATE TABLE customer_wishlists (
    customer_id INT NOT NULL,
    product_id BIGINT NOT NULL,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    UNIQUE (customer_id, product_id)
);

CREATE TABLE customer_carts (
    customer_id INT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT DEFAULT 1,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    UNIQUE (customer_id, product_id)
);