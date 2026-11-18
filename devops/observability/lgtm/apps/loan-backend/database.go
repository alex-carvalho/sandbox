package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var db *sql.DB

func initDatabase() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://loanuser:loanpass@postgres.loan-app.svc.cluster.local:5432/loandb?sslmode=disable"
	}

	var err error
	for attempt := 1; attempt <= 20; attempt++ {
		// Industry-standard SQL auto-instrumentation via otelsql
		db, err = otelsql.Open("postgres", dsn,
			otelsql.WithAttributes(
				semconv.DBSystemPostgreSQL,
				attribute.String("db.name", "loandb"),
				attribute.String("peer.service", "postgresql"),
			),
			otelsql.WithSpanOptions(otelsql.SpanOptions{
				OmitConnResetSession: true,
				OmitConnPrepare:      true,
				OmitRows:             true,
				OmitConnectorConnect: true,
			}),
		)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		fmt.Printf("Waiting for PostgreSQL (attempt %d/20): %v\n", attempt, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("could not connect to PostgreSQL: %w", err)
	}

	_ = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))

	schemaSQL := `
	CREATE TABLE IF NOT EXISTS loan_applications (
		id SERIAL PRIMARY KEY,
		application_id VARCHAR(32) UNIQUE NOT NULL,
		applicant_name VARCHAR(128) NOT NULL,
		loan_type VARCHAR(32) NOT NULL,
		amount NUMERIC(12, 2) NOT NULL CONSTRAINT chk_valid_loan_amount CHECK (amount > 0 AND amount <> 666),
		term_months INT NOT NULL,
		annual_income NUMERIC(12, 2) NOT NULL,
		credit_score INT NOT NULL,
		decision VARCHAR(32) NOT NULL,
		apr NUMERIC(5, 2) DEFAULT 0,
		monthly_payment NUMERIC(12, 2) DEFAULT 0,
		dti_ratio NUMERIC(6, 2) DEFAULT 0,
		risk_score NUMERIC(5, 2) DEFAULT 0,
		db_query_ms NUMERIC(8, 2) DEFAULT 0,
		trace_id VARCHAR(64) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS fraud_blacklist (
		id SERIAL PRIMARY KEY,
		pattern VARCHAR(64) UNIQUE NOT NULL,
		reason TEXT NOT NULL,
		risk_level VARCHAR(32) NOT NULL
	);

	INSERT INTO fraud_blacklist (pattern, reason, risk_level)
	VALUES
		('fraud', 'Confirmed synthetic identity fraud ring match in OFAC/Fraud DB', 'CRITICAL'),
		('blacklist', 'Account holder listed on interbank default blacklist', 'CRITICAL')
	ON CONFLICT (pattern) DO NOTHING;
	`
	_, err = db.Exec(schemaSQL)
	return err
}
