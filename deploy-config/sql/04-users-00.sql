-- Users Database
\c vuecom_users;

CREATE TABLE app_data (
    app_name TEXT NOT NULL,
    admin_route TEXT,
    app_logo TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    settings JSONB DEFAULT '{}'
);

CREATE TABLE continents (
    id INT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE countries (
    id INT PRIMARY KEY,
    name TEXT NOT NULL,
    -- e.g., 'US', 'NG'
    code TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL,
    currency TEXT NOT NULL,
    continent_id INT NOT NULL,
    FOREIGN KEY (continent_id) REFERENCES continents(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_countries_code ON countries (code);

CREATE TABLE states (
    id INT PRIMARY KEY,
    country_id INT NOT NULL,
    name TEXT NOT NULL,
    FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_states_country_id ON states (country_id);

CREATE TABLE cities (
    id INT PRIMARY KEY,
    state_id INT NOT NULL,
    name TEXT NOT NULL,
    FOREIGN KEY (state_id) REFERENCES states(id) ON DELETE CASCADE
);


CREATE TABLE backend_users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL,
    phone_number TEXT,
    role TEXT DEFAULT 'staff', -- Set to allow custom roles
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT,
    image_url TEXT,
    country_id INT,
    is_2fa_enabled BOOLEAN DEFAULT FALSE,
    is_email_verified BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (created_by) REFERENCES backend_users(id),
    FOREIGN KEY (country_id) REFERENCES countries(id)
);

CREATE INDEX IF NOT EXISTS backend_user_email_idx ON backend_users USING hash (email);
CREATE INDEX IF NOT EXISTS backend_user_username_idx ON backend_users USING hash (username);
CREATE INDEX IF NOT EXISTS backend_user_role_idx ON backend_users(role);

CREATE TABLE backend_signup_data (
    id UUID NOT NULL,
    token TEXT NOT NULL
);

CREATE TABLE backend_2fa_tokens (
    user_id INT NOT NULL,
    encrypted_token TEXT NOT NULL,
    hashed_backup_token TEXT NOT NULL,
    UNIQUE (user_id),
    FOREIGN KEY (user_id) REFERENCES backend_users(id) ON DELETE CASCADE
);

CREATE TABLE backend_sessions (
    user_id INT NOT NULL,
    refresh_token_hash TEXT NOT NULL,
    device_id UUID NOT NULL,
    last_ip TEXT NOT NULL,
    user_agent TEXT,
    used BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES backend_users(id) ON DELETE CASCADE
);

-- CREATE INDEX idx_backend_sessions_refresh_token ON backend_sessions(refresh_token_hash);
CREATE INDEX IF NOT EXISTS idx_backend_sessions_user_id ON backend_sessions(user_id);


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


-- FUNCTIONS
create or replace function check_if_creator(user_id int, assumed_creator int)
returns bool
language plpgsql
as $$
declare
    tmp_id int;
begin
    select created_by into tmp_id from backend_users where id = user_id;

    if found then
        loop
            if tmp_id = assumed_creator then
                return TRUE;
            else
                select created_by into tmp_id from backend_users where id = tmp_id;

                if tmp_id = NULL then
                    return FALSE;
                end if;
            end if;
        end loop;
    end if;
end$$;

\c postgres;
