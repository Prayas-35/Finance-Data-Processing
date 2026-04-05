CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

CREATE INDEX IF NOT EXISTS idx_records_user_date ON financial_records(user_id, date DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_records_type_date ON financial_records(type, date DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_records_category_date ON financial_records(category, date DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_records_recent ON financial_records(created_at DESC) WHERE deleted_at IS NULL;
