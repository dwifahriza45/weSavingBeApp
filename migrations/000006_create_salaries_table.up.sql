CREATE TABLE IF NOT EXISTS salaries (
    id BIGSERIAL PRIMARY KEY,
    salary_id VARCHAR(25) UNIQUE NOT NULL,
    user_id VARCHAR(25) NOT NULL,
    amount BIGINT NOT NULL,
    source VARCHAR(50),
    description VARCHAR(255),
    received_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);