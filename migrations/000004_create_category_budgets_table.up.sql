CREATE TABLE IF NOT EXISTS category_budgets (
    id BIGSERIAL PRIMARY KEY,
    budget_id VARCHAR(25) UNIQUE NOT NULL,
    user_id VARCHAR(25) NOT NULL,
    category_id VARCHAR(25) NOT NULL,
    allocated_amount BIGINT NOT NULL,
    used_amount BIGINT NOT NULL DEFAULT 0,
    period VARCHAR(10) NOT NULL DEFAULT 'monthly',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
);