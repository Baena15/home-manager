-- ─── Finance schema for Home Manager ───────────────────────────────

-- ─── expenses ────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL,
    description VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    visibility VARCHAR(20) NOT NULL DEFAULT 'private' CHECK (visibility IN ('private', 'shared')),
    split_percentage DECIMAL(5, 2) NOT NULL DEFAULT 50.00 CHECK (split_percentage >= 0 AND split_percentage <= 100),
    expense_date DATE NOT NULL DEFAULT CURRENT_DATE,
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses(expense_date);
CREATE INDEX IF NOT EXISTS idx_expenses_visibility ON expenses(visibility);
CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses(category);

-- ─── incomes ─────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS incomes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL,
    description VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    visibility VARCHAR(20) NOT NULL DEFAULT 'private' CHECK (visibility IN ('private', 'shared')),
    income_date DATE NOT NULL DEFAULT CURRENT_DATE,
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_incomes_user_id ON incomes(user_id);
CREATE INDEX IF NOT EXISTS idx_incomes_date ON incomes(income_date);
CREATE INDEX IF NOT EXISTS idx_incomes_visibility ON incomes(visibility);
CREATE INDEX IF NOT EXISTS idx_incomes_category ON incomes(category);

-- ─── down migration ──────────────────────────────────────────────────
-- DROP TABLE IF EXISTS incomes;
-- DROP TABLE IF EXISTS expenses;
