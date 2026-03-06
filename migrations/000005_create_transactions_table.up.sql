CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    transaction_id VARCHAR(25) UNIQUE NOT NULL,
    user_id VARCHAR(25) NOT NULL,
    category_id VARCHAR(25) NOT NULL,
    amount BIGINT NOT NULL,
    type VARCHAR(10) NOT NULL,
    description VARCHAR(255),
    transaction_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
);