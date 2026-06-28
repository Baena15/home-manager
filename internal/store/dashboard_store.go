// Package store provides data access for dashboard aggregations.
package store

import (
	"context"
	"fmt"
)

// ─── MonthlySummary ─────────────────────────────────────────────────

// MonthlySummary holds aggregated income and expense totals for a month.
type MonthlySummary struct {
	IncomeTotal  float64 `json:"income_total"`
	ExpenseTotal float64 `json:"expense_total"`
	Balance      float64 `json:"balance"`
}

// ─── MonthData ──────────────────────────────────────────────────────

// MonthData holds income and expense totals for a single month.
type MonthData struct {
	Month        string  `json:"month"`
	IncomeTotal  float64 `json:"income_total"`
	ExpenseTotal float64 `json:"expense_total"`
	Balance      float64 `json:"balance"`
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
		WHERE (user_id = $1 OR visibility = 'shared')
		  AND TO_CHAR(income_date, 'YYYY-MM') = $2
	`
	expenseQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE (user_id = $1 OR visibility = 'shared')
		  AND TO_CHAR(expense_date, 'YYYY-MM') = $2
	`

	summary := &MonthlySummary{}
	if err := s.db.QueryRowContext(ctx, incomeQuery, userID, month).Scan(&summary.IncomeTotal); err != nil {
		return nil, fmt.Errorf("failed to calculate income total: %w", err)
	}
	if err := s.db.QueryRowContext(ctx, expenseQuery, userID, month).Scan(&summary.ExpenseTotal); err != nil {
		return nil, fmt.Errorf("failed to calculate expense total: %w", err)
	}

	summary.Balance = summary.IncomeTotal - summary.ExpenseTotal
	return summary, nil
}

// MonthlyTotals returns income and expense totals for each month of the year.
func (s *DashboardStore) MonthlyTotals(ctx context.Context, userID, year string) ([]MonthData, error) {
	query := `
		SELECT months.month, COALESCE(income_totals.total, 0) AS income_total, COALESCE(expense_totals.total, 0) AS expense_total
		FROM generate_series(1, 12) AS months(month)
		LEFT JOIN (
			SELECT EXTRACT(MONTH FROM income_date)::int AS month, COALESCE(SUM(amount), 0) AS total
			FROM incomes
			WHERE (user_id = $1 OR visibility = 'shared')
			  AND EXTRACT(YEAR FROM income_date)::text = $2
			GROUP BY month
		) AS income_totals ON months.month = income_totals.month
		LEFT JOIN (
			SELECT EXTRACT(MONTH FROM expense_date)::int AS month, COALESCE(SUM(amount), 0) AS total
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
		if err := rows.Scan(&month, &item.IncomeTotal, &item.ExpenseTotal); err != nil {
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
