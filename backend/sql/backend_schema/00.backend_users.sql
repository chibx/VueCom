CREATE SCHEMA IF NOT EXISTS backend;

CREATE TABLE backend.app_data (
    app_name VARCHAR(100) NOT NULL,
    admin_route VARCHAR(100) NOT NULL,
    app_logo TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE backend.countries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(5) NOT NULL UNIQUE,
    -- e.g., 'US', 'NG'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_countries_code ON backend.countries (code);

CREATE TABLE backend.states (
    id SERIAL PRIMARY KEY,
    country_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    FOREIGN KEY (country_id) REFERENCES backend.countries(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_states_name ON backend.states (name);

CREATE TABLE backend.backend_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    phone VARCHAR(20),
    role TEXT DEFAULT 'staff' -- Set to allow custom roles
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP image TEXT,
    country_id INT,
    is_email_verified BOOLEAN DEFAULT FALSE FOREIGN KEY (country_id) REFERENCES backend.countries(id)
);

CREATE TABLE backend.backend_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES backend.backend_users(id) ON DELETE CASCADE
);

CREATE TABLE backend.backend_otps (
    user_id INT NOT NULL,
    code VARCHAR(10) NOT NULL,
    expiry_date TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES backend.backend_users(id) ON DELETE CASCADE
);

CREATE TABLE backend.backend_user_activities (
    user_id INT NOT NULL,
    log_title VARCHAR(100) NOT NULL,
    -- e.g., "Login", "Password Change"
    activity TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES backend.backend_users(id) ON DELETE
    SET
        NULL
);

CREATE TABLE backend.api_keys (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    key_prefix VARCHAR(16) NOT NULL UNIQUE,
    key_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES backend.backend_users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS api_keys_prefix_idx ON backend.api_keys (key_prefix);

CREATE TABLE backend.password_reset_requests (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    reset_token VARCHAR(255) NOT NULL UNIQUE,
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES backend.backend_users(id) ON DELETE CASCADE
);