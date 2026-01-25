-- Users Database
\c vuecom_users;

CREATE TABLE app_data (
    app_name TEXT NOT NULL,
    admin_route TEXT NOT NULL,
    app_logo TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    settings JSONB DEFAULT '{}'
);

CREATE TABLE countries (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    -- e.g., 'US', 'NG'
    code TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_countries_code ON countries (code);

CREATE TABLE states (
    id SERIAL PRIMARY KEY,
    country_id INT NOT NULL,
    name TEXT NOT NULL,
    FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_states_name ON states (name);

CREATE TABLE backend_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL,
    phone_number VARCHAR(20),
    role TEXT DEFAULT 'staff', -- Set to allow custom roles
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    image TEXT,
    country_id INT,
    is_email_verified BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (country_id) REFERENCES countries(id)
);

CREATE INDEX IF NOT EXISTS backend_user_email_idx ON backend_users USING hash (email);
CREATE INDEX IF NOT EXISTS backend_user_username_idx ON backend_users USING hash (username);
CREATE INDEX IF NOT EXISTS backend_user_role_idx ON backend_users(role);

-- CREATE TABLE backend_sessions (
--     id SERIAL PRIMARY KEY,
--     user_id INT NOT NULL,
--     token TEXT NOT NULL,
--     ip_address VARCHAR(45),
--     user_agent TEXT,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     expires_at TIMESTAMP NOT NULL,
--     FOREIGN KEY (user_id) REFERENCES backend_users(id) ON DELETE CASCADE
-- );

-- CREATE INDEX IF NOT EXISTS backend_sessions_userid_idx ON backend_sessions (user_id);
-- CREATE INDEX IF NOT EXISTS backend_sessions_expires_at_idx ON backend_sessions (expires_at);

CREATE TABLE backend_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT NOT NULL,
    refresh_token_hash TEXT NOT NULL,
    device_id UUID NOT NULL,
    last_ip TEXT NOT NULL,
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES backend_users(id) ON DELETE CASCADE
);

-- CREATE INDEX idx_backend_sessions_refresh_token ON backend_sessions(refresh_token_hash);
CREATE INDEX IF NOT EXISTS idx_backend_sessions_user_id ON backend_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_backend_sessions_device ON backend_sessions(device_id);


CREATE TABLE backend_otps (
    user_id INT NOT NULL,
    code VARCHAR(10) NOT NULL,
    expiry_date TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES backend_users(id) ON DELETE CASCADE,
    UNIQUE (code, expiry_date)
);

CREATE INDEX IF NOT EXISTS backend_otps_code_idx ON backend_otps (code);
CREATE INDEX IF NOT EXISTS backend_otps_expiry_idx ON backend_otps (expiry_date);

CREATE TABLE backend_user_activities (
    user_id INT NOT NULL,
    -- e.g., "Login", "Password Change"
    log_title TEXT NOT NULL,
    activity TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES backend_users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS backend_user_activities_idx ON backend_user_activities (user_id);
CREATE INDEX IF NOT EXISTS backend_user_activities_title_idx ON backend_user_activities USING hash(log_title);

CREATE TABLE api_keys (
    user_id INT NOT NULL,
    key_prefix VARCHAR(16) NOT NULL UNIQUE,
    key_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES backend_users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS api_keys_userid_idx ON api_keys (user_id);
CREATE INDEX IF NOT EXISTS api_keys_prefix_idx ON api_keys (key_prefix);

CREATE TABLE password_reset_requests (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    reset_token TEXT NOT NULL UNIQUE,
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES backend_users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS password_reset_requests_token_idx ON password_reset_requests USING hash(reset_token);
CREATE INDEX IF NOT EXISTS password_reset_requests_used_idx ON password_reset_requests (used);


\c postgres;
