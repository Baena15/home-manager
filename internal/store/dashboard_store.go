// Package store provides data access for dashboard aggregations.
package store

import (
	"context"
	"database/sql"
	"fmt"
)

// ─── MonthlySummary ─────────────────────────────────────────────────

// MonthlySummary holds aggregated income and expense totals for a month.
type MonthlySummary struct {
	IncomeTotal         float64 `json:"income_total"`
	SharedExpenseTotal  float64 `json:"shared_expense_total"`
	PrivateExpenseTotal float64 `json:"private_expense_total"`
	ExpenseTotal        float64 `json:"expense_total"`
	Balance             float64 `json:"balance"`
}

// ─── MonthData ──────────────────────────────────────────────────────

// MonthData holds income and expense totals for a single month.
type MonthData struct {
	Month               string  `json:"month"`
	IncomeTotal         float64 `json:"income_total"`
	SharedExpenseTotal  float64 `json:"shared_expense_total"`
	PrivateExpenseTotal float64 `json:"private_expense_total"`
	ExpenseTotal        float64 `json:"expense_total"`
	Balance             float64 `json:"balance"`
}

// ─── PartnerBalance ─────────────────────────────────────────────────

// PartnerBalance holds the settlement balance between two household users.
type PartnerBalance struct {
	PartnerID    string  `json:"partner_id"`
	PartnerEmail string  `json:"partner_email"`
	Amount       float64 `json:"amount"`
	YouOwe       bool    `json:"you_owe"`
	YouAreOwed   float64 `json:"you_are_owed"`
	YouOweAmount float64 `json:"you_owe_amount"`
}

// ─── DashboardStore ─────────────────────────────────────────────────

// DashboardStore handles dashboard aggregation queries.
type DashboardStore struct {
	db *DB
}

// NewDashboardStore creates a new DashboardStore.
func NewDashboardStore(db *DB) *DashboardStore {
	return &DashboardStore{db: db}
}

// MonthlySummary returns aggregated totals for the given user and month.
func (s *DashboardStore) MonthlySummary(ctx context.Context, userID, month string) (*MonthlySummary, error) {
	incomeQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM incomes
		WHERE user_id = $1
		  AND TO_CHAR(income_date, 'YYYY-MM') = $2
	`
	expenseQuery := `
		SELECT
			COALESCE(SUM(CASE
				WHEN visibility = 'shared' AND user_id = $1 THEN amount * split_percentage / 100
				WHEN visibility = 'shared' AND user_id != $1 THEN amount * (100 - split_percentage) / 100
				WHEN visibility = 'private' AND user_id = $1 THEN amount
				ELSE 0
			END), 0) AS expense_total,
			COALESCE(SUM(CASE
				WHEN visibility = 'shared' AND user_id = $1 THEN amount * split_percentage / 100
				WHEN visibility = 'shared' AND user_id != $1 THEN amount * (100 - split_percentage) / 100
				ELSE 0
			END), 0) AS shared_expense_total,
			COALESCE(SUM(CASE
				WHEN visibility = 'private' AND user_id = $1 THEN amount
				ELSE 0
			END), 0) AS private_expense_total
		FROM expenses
		WHERE (user_id = $1 OR visibility = 'shared')
		  AND TO_CHAR(expense_date, 'YYYY-MM') = $2
	`

	summary := &MonthlySummary{}
	if err := s.db.QueryRowContext(ctx, incomeQuery, userID, month).Scan(&summary.IncomeTotal); err != nil {
		return nil, fmt.Errorf("failed to calculate income total: %w", err)
	}
	if err := s.db.QueryRowContext(ctx, expenseQuery, userID, month).Scan(
		&summary.ExpenseTotal,
		&summary.SharedExpenseTotal,
		&summary.PrivateExpenseTotal,
	); err != nil {
		return nil, fmt.Errorf("failed to calculate expense total: %w", err)
	}

	summary.Balance = summary.IncomeTotal - summary.ExpenseTotal
	return summary, nil
}

// PartnerBalance calculates the net settlement balance between the given user
// and the other household user across all shared expenses and recorded settlements.
// YouAreOwed is the partner's share of shared expenses paid by the user.
// YouOweAmount is the user's share of shared expenses paid by the partner.
// Amount is the net balance: YouAreOwed minus YouOweAmount minus settlements.
func (s *DashboardStore) PartnerBalance(ctx context.Context, userID string) (*PartnerBalance, error) {
	partnerQuery := `
		SELECT id, email
		FROM users
		WHERE id != $1
		LIMIT 1
	`

	balance := &PartnerBalance{}
	if err := s.db.QueryRowContext(ctx, partnerQuery, userID).Scan(&balance.PartnerID, &balance.PartnerEmail); err != nil {
		if err == sql.ErrNoRows {
			return &PartnerBalance{}, nil
		}
		return nil, fmt.Errorf("failed to get partner: %w", err)
	}

	expenseOwedQuery := `
		SELECT COALESCE(SUM(amount * (100 - split_percentage) / 100), 0)
		FROM expenses
		WHERE visibility = 'shared'
		  AND user_id = $1
		  AND settled_at IS NULL
	`

	expenseOweQuery := `
		SELECT COALESCE(SUM(amount * (100 - split_percentage) / 100), 0)
		FROM expenses
		WHERE visibility = 'shared'
		  AND user_id != $1
		  AND settled_at IS NULL
	`

	settlementsReceivedQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM settlements
		WHERE to_user_id = $1 AND from_user_id != $1
	`

	settlementsSentQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM settlements
		WHERE from_user_id = $1 AND to_user_id != $1
	`

	if err := s.db.QueryRowContext(ctx, expenseOwedQuery, userID).Scan(&balance.YouAreOwed); err != nil {
		return nil, fmt.Errorf("failed to calculate amount owed to user: %w", err)
	}
	if err := s.db.QueryRowContext(ctx, expenseOweQuery, userID).Scan(&balance.YouOweAmount); err != nil {
		return nil, fmt.Errorf("failed to calculate amount user owes: %w", err)
	}

	var settlementsReceived float64
	var settlementsSent float64
	if err := s.db.QueryRowContext(ctx, settlementsReceivedQuery, userID).Scan(&settlementsReceived); err != nil {
		return nil, fmt.Errorf("failed to calculate settlements received: %w", err)
	}
	if err := s.db.QueryRowContext(ctx, settlementsSentQuery, userID).Scan(&settlementsSent); err != nil {
		return nil, fmt.Errorf("failed to calculate settlements sent: %w", err)
	}

	net := (balance.YouAreOwed - settlementsReceived) - (balance.YouOweAmount - settlementsSent)
	if net >= 0 {
		balance.Amount = net
		balance.YouOwe = false
	} else {
		balance.Amount = -net
		balance.YouOwe = true
	}

	return balance, nil
}

// MonthlyTotals returns income and expense totals for each month of the year.
func (s *DashboardStore) MonthlyTotals(ctx context.Context, userID, year string) ([]MonthData, error) {
	query := `
		SELECT
			months.month,
			COALESCE(income_totals.total, 0) AS income_total,
			COALESCE(expense_totals.shared_total, 0) AS shared_expense_total,
			COALESCE(expense_totals.private_total, 0) AS private_expense_total,
			COALESCE(expense_totals.shared_total, 0) + COALESCE(expense_totals.private_total, 0) AS expense_total
		FROM generate_series(1, 12) AS months(month)
		LEFT JOIN (
			SELECT EXTRACT(MONTH FROM income_date)::int AS month, COALESCE(SUM(amount), 0) AS total
			FROM incomes
			WHERE user_id = $1
			  AND EXTRACT(YEAR FROM income_date)::text = $2
			GROUP BY month
		) AS income_totals ON months.month = income_totals.month
		LEFT JOIN (
			SELECT
				EXTRACT(MONTH FROM expense_date)::int AS month,
				COALESCE(SUM(CASE
					WHEN visibility = 'shared' AND user_id = $1 THEN amount * split_percentage / 100
					WHEN visibility = 'shared' AND user_id != $1 THEN amount * (100 - split_percentage) / 100
					ELSE 0
				END), 0) AS shared_total,
				COALESCE(SUM(CASE
					WHEN visibility = 'private' AND user_id = $1 THEN amount
					ELSE 0
				END), 0) AS private_total
			FROM expenses
			WHERE (user_id = $1 OR visibility = 'shared')
			  AND EXTRACT(YEAR FROM expense_date)::text = $2
			GROUP BY month
		) AS expense_totals ON months.month = expense_totals.month
		ORDER BY months.month
	`

	rows, err := s.db.QueryContext(ctx, query, userID, year)
	if err != nil {
		return nil, fmt.Errorf("failed to list monthly totals: %w", err)
	}
	defer rows.Close()

	var data []MonthData
	for rows.Next() {
		var month int
		var item MonthData
		if err := rows.Scan(&month, &item.IncomeTotal, &item.SharedExpenseTotal, &item.PrivateExpenseTotal, &item.ExpenseTotal); err != nil {
			return nil, fmt.Errorf("failed to scan monthly total: %w", err)
		}
		item.Month = fmt.Sprintf("%s-%02d", year, month)
		item.Balance = item.IncomeTotal - item.ExpenseTotal
		data = append(data, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating monthly totals: %w", err)
	}

	return data, nil
}
