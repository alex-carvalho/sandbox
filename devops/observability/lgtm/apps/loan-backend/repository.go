package main

import (
	"context"
	"database/sql"
	"time"
)

type LoanRepository struct {
	db *sql.DB
}

var loanRepo *LoanRepository

func NewLoanRepository(db *sql.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func elapsedMsSince(t time.Time) float64 {
	return float64(time.Since(t).Microseconds()) / 1000.0
}

// AnalyzeSlowPortfolioExposure runs an intentionally slow analytical query (pg_sleep)
// to demonstrate SQL bottleneck detection in Tempo and Mimir.
func (r *LoanRepository) AnalyzeSlowPortfolioExposure(ctx context.Context, loanType string) (float64, error) {
	const slowSQL = `SELECT pg_sleep(0.85), COUNT(*), COALESCE(SUM(amount), 0) FROM loan_applications WHERE loan_type = $1`
	start := time.Now()
	var sleepDummy string
	var count int
	var totalExposure float64
	err := r.db.QueryRowContext(ctx, slowSQL, loanType).Scan(&sleepDummy, &count, &totalExposure)
	return elapsedMsSince(start), err
}

func (r *LoanRepository) CheckApplicantHistory(ctx context.Context, applicantName string) (float64, error) {
	const fastSQL = `SELECT COUNT(*), COALESCE(AVG(apr), 0) FROM loan_applications WHERE applicant_name = $1`
	start := time.Now()
	var count int
	var avgAPR float64
	err := r.db.QueryRowContext(ctx, fastSQL, applicantName).Scan(&count, &avgAPR)
	return elapsedMsSince(start), err
}

func (r *LoanRepository) InsertPendingLoan(ctx context.Context, req LoanRequest, amount float64, traceID string) (float64, error) {
	const insertSQL = `
		INSERT INTO loan_applications (
			application_id, applicant_name, loan_type, amount,
			term_months, annual_income, credit_score, decision, trace_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, 'PENDING', $8)`
	start := time.Now()
	_, err := r.db.ExecContext(ctx, insertSQL,
		req.ApplicationID, req.ApplicantName, req.LoanType, amount,
		req.TermMonths, req.AnnualIncome, req.CreditScore, traceID,
	)
	return elapsedMsSince(start), err
}

func (r *LoanRepository) UpdateLoanDecision(ctx context.Context, applicationID string, uw *UnderwritingResult, totalDBMs float64) (float64, error) {
	const updateSQL = `
		UPDATE loan_applications
		SET decision = $1, apr = $2, monthly_payment = $3, dti_ratio = $4, risk_score = $5, db_query_ms = $6
		WHERE application_id = $7`
	start := time.Now()
	_, err := r.db.ExecContext(ctx, updateSQL,
		uw.Decision, uw.APR, uw.MonthlyPayment, uw.DTIRatio, uw.RiskScore, totalDBMs, applicationID,
	)
	return elapsedMsSince(start), err
}

func (r *LoanRepository) ListRecentLoans(ctx context.Context, limit int) ([]LoanResponse, float64, error) {
	const selectSQL = `
		SELECT application_id, applicant_name, loan_type, amount, credit_score,
		       decision, apr, db_query_ms, trace_id, created_at
		FROM loan_applications
		ORDER BY id DESC
		LIMIT $1`
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, selectSQL, limit)
	if err != nil {
		return nil, elapsedMsSince(start), err
	}
	defer rows.Close()

	var items []LoanResponse
	for rows.Next() {
		var item LoanResponse
		var createdAt time.Time
		if scanErr := rows.Scan(
			&item.ApplicationID,
			&item.ApplicantName,
			&item.LoanType,
			&item.Amount,
			&item.CreditScore,
			&item.Decision,
			&item.APR,
			&item.DBQueryMs,
			&item.TraceID,
			&createdAt,
		); scanErr != nil {
			return nil, elapsedMsSince(start), scanErr
		}
		item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		items = append(items, item)
	}
	return items, elapsedMsSince(start), nil
}
