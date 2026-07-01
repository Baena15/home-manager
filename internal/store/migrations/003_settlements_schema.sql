-- ─── Settlements schema for Home Manager ───────────────────────────

-- ─── expenses: add settlement tracking ─────────────────────────────
ALTER TABLE expenses
    ADD COLUMN IF NOT EXISTS settled_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS settled_by UUID REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_expenses_settled_at ON expenses(settled_at);

-- ─── settlements ───────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS settlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    description VARCHAR(255),
    settlement_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_settlements_from_user_id ON settlements(from_user_id);
CREATE INDEX IF NOT EXISTS idx_settlements_to_user_id ON settlements(to_user_id);
CREATE INDEX IF NOT EXISTS idx_settlements_date ON settlements(settlement_date);

-- ─── down migration ────────────────────────────────────────────────
-- DROP TABLE IF EXISTS settlements;
-- ALTER TABLE expenses DROP COLUMN IF EXISTS settled_at;
-- ALTER TABLE expenses DROP COLUMN IF EXISTS settled_by;
