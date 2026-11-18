package main

const serviceName = "loan-backend"

type LoanRequest struct {
	ApplicationID string  `json:"application_id"`
	ApplicantName string  `json:"applicant_name"`
	LoanType      string  `json:"loan_type"`
	Amount        float64 `json:"amount"`
	TermMonths    int     `json:"term_months"`
	AnnualIncome  float64 `json:"annual_income"`
	CreditScore   int     `json:"credit_score"`
	SimulateMode  string  `json:"simulate_mode"` // normal, slow_db, fraud_error, bureau_timeout, db_error
}

type LoanResponse struct {
	ApplicationID  string  `json:"application_id"`
	ApplicantName  string  `json:"applicant_name,omitempty"`
	LoanType       string  `json:"loan_type,omitempty"`
	Amount         float64 `json:"amount,omitempty"`
	CreditScore    int     `json:"credit_score,omitempty"`
	Decision       string  `json:"decision"`
	APR            float64 `json:"apr"`
	MonthlyPayment float64 `json:"monthly_payment"`
	DTIRatio       float64 `json:"dti_ratio"`
	RiskScore      float64 `json:"risk_score"`
	Reason         string  `json:"reason"`
	TraceID        string  `json:"trace_id"`
	DBQueryMs      float64 `json:"db_query_ms"`
	ProcessingMs   int64   `json:"processing_ms"`
	CreatedAt      string  `json:"created_at,omitempty"`
}

type UnderwritingResult struct {
	ApplicationID    string  `json:"application_id"`
	Decision         string  `json:"decision"`
	APR              float64 `json:"apr"`
	MonthlyPayment   float64 `json:"monthly_payment"`
	DTIRatio         float64 `json:"dti_ratio"`
	RiskScore        float64 `json:"risk_score"`
	Reason           string  `json:"reason"`
	UnderwritingDBMs float64 `json:"underwriting_db_ms"`
}
