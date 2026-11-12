package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

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

var (
	db                  *sql.DB
	httpClient          *http.Client
	underwritingURL     string
	tracer              trace.Tracer
	otelLogger          otellog.Logger
	underwritingDurHist metric.Float64Histogram
	dbQueryDurHist      metric.Float64Histogram
	decisionsCounter    metric.Int64Counter
	dbErrorsCounter     metric.Int64Counter
	requestedAmountHist metric.Float64Histogram
)

func initTelemetry(ctx context.Context) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "alloy.observability.svc.cluster.local:4317"
	}
	endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceNamespaceKey.String("loan-app"),
			semconv.ServiceVersionKey.String("2.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	traceExp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp, sdktrace.WithBatchTimeout(2*time.Second)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	tracer = tp.Tracer(serviceName)

	metricExp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp, sdkmetric.WithInterval(5*time.Second))),
	)
	otel.SetMeterProvider(mp)
	meter := mp.Meter(serviceName)

	underwritingDurHist, err = meter.Float64Histogram(
		"loan_underwriting_duration",
		metric.WithUnit("s"),
		metric.WithDescription("End-to-end backend orchestration & underwriting duration in seconds"),
		metric.WithExplicitBucketBoundaries(0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 0.75, 1.0, 2.5, 5.0),
	)
	if err != nil {
		return nil, err
	}

	dbQueryDurHist, err = meter.Float64Histogram(
		"loan_db_query_duration",
		metric.WithUnit("s"),
		metric.WithDescription("PostgreSQL SQL query execution duration in seconds (with Exemplars)"),
		metric.WithExplicitBucketBoundaries(0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5),
	)
	if err != nil {
		return nil, err
	}

	decisionsCounter, err = meter.Int64Counter(
		"loan_decisions",
		metric.WithDescription("Total loan underwriting decisions by decision and loan_type"),
	)
	if err != nil {
		return nil, err
	}

	dbErrorsCounter, err = meter.Int64Counter(
		"loan_db_errors",
		metric.WithDescription("Total PostgreSQL query errors in loan-backend"),
	)
	if err != nil {
		return nil, err
	}

	requestedAmountHist, err = meter.Float64Histogram(
		"loan_requested_amount_dollars",
		metric.WithDescription("Distribution of requested loan amounts in USD"),
		metric.WithExplicitBucketBoundaries(5000, 15000, 30000, 50000, 100000, 250000, 500000, 1000000),
	)
	if err != nil {
		return nil, err
	}

	logExp, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create log exporter: %w", err)
	}
	lp := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
	)
	otelLogger = lp.Logger(serviceName)

	shutdown := func(c context.Context) error {
		_ = lp.Shutdown(c)
		_ = mp.Shutdown(c)
		return tp.Shutdown(c)
	}
	return shutdown, nil
}

func emitJSONLog(ctx context.Context, level string, msg string, fields map[string]any) {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	traceID := ""
	spanID := ""
	if spanCtx.IsValid() {
		traceID = spanCtx.TraceID().String()
		spanID = spanCtx.SpanID().String()
	}

	payload := map[string]any{
		"timestamp":    time.Now().UTC().Format(time.RFC3339Nano),
		"level":        level,
		"service_name": serviceName,
		"runtime":      "go",
		"msg":          msg,
		"trace_id":     traceID,
		"span_id":      spanID,
	}
	for k, v := range fields {
		payload[k] = v
	}

	raw, _ := json.Marshal(payload)
	line := string(raw)
	fmt.Println(line)

	var rec otellog.Record
	rec.SetTimestamp(time.Now())
	rec.SetBody(otellog.StringValue(line))
	switch level {
	case "ERROR":
		rec.SetSeverity(otellog.SeverityError)
		rec.SetSeverityText("ERROR")
	case "WARN":
		rec.SetSeverity(otellog.SeverityWarn)
		rec.SetSeverityText("WARN")
	default:
		rec.SetSeverity(otellog.SeverityInfo)
		rec.SetSeverityText("INFO")
	}
	rec.AddAttributes(
		otellog.String("trace_id", traceID),
		otellog.String("span_id", spanID),
		otellog.String("service.name", serviceName),
	)
	otelLogger.Emit(ctx, rec)
}

func initDatabase() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://loanuser:loanpass@postgres.loan-app.svc.cluster.local:5432/loandb?sslmode=disable"
	}

	var err error
	for attempt := 1; attempt <= 20; attempt++ {
		db, err = sql.Open("postgres", dsn)
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

// runTracedDBQuery wraps a SQL operation in an OpenTelemetry CLIENT span with db.system=postgresql
// and records its execution time in loan_db_query_duration_seconds (generating Exemplars).
func runTracedDBQuery(ctx context.Context, spanName, dbOp, table, query, scenario string, fn func(context.Context) error) (float64, error) {
	ctx, span := tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.name", "loandb"),
			attribute.String("db.operation", dbOp),
			attribute.String("db.sql.table", table),
			attribute.String("db.statement", query),
			attribute.String("peer.service", "postgresql"),
		),
	)
	defer span.End()

	start := time.Now()
	err := fn(ctx)
	elapsedSec := time.Since(start).Seconds()
	elapsedMs := float64(time.Since(start).Microseconds()) / 1000.0

	dbQueryDurHist.Record(ctx, elapsedSec,
		metric.WithAttributes(
			attribute.String("service", serviceName),
			attribute.String("db_operation", dbOp),
			attribute.String("table", table),
			attribute.String("scenario", scenario),
		),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		dbErrorsCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("db_operation", dbOp),
				attribute.String("table", table),
			),
		)
	}
	return elapsedMs, err
}

func handleUnderwrite(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid request payload")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.SimulateMode == "" {
		req.SimulateMode = "normal"
	}

	span.SetAttributes(
		attribute.String("loan.application_id", req.ApplicationID),
		attribute.String("loan.type", req.LoanType),
		attribute.Float64("loan.amount", req.Amount),
		attribute.Int("loan.credit_score", req.CreditScore),
		attribute.String("loan.simulate_mode", req.SimulateMode),
	)

	emitJSONLog(ctx, "INFO", "Go backend orchestrating loan application", map[string]any{
		"application_id": req.ApplicationID,
		"loan_type":      req.LoanType,
		"amount":         req.Amount,
		"credit_score":   req.CreditScore,
		"simulate_mode":  req.SimulateMode,
	})

	requestedAmountHist.Record(ctx, req.Amount,
		metric.WithAttributes(attribute.String("loan_type", req.LoanType)),
	)

	var totalDBMs float64

	// 1. Query historical exposure in PostgreSQL (fast or slow analytical query if slow_db / amount > 900k)
	if req.SimulateMode == "slow_db" || req.Amount > 900000 {
		slowSQL := `SELECT pg_sleep(0.85), COUNT(*), COALESCE(SUM(amount), 0) FROM loan_applications WHERE loan_type = $1`
		qMs, err := runTracedDBQuery(ctx, "postgres.slow_portfolio_exposure_analysis", "SELECT", "loan_applications", slowSQL, "slow_db", func(qCtx context.Context) error {
			var sleepDummy string
			var count int
			var totalExposure float64
			return db.QueryRowContext(qCtx, slowSQL, req.LoanType).Scan(&sleepDummy, &count, &totalExposure)
		})
		totalDBMs += qMs
		emitJSONLog(ctx, "WARN", "Slow PostgreSQL portfolio exposure query detected (>800ms)", map[string]any{
			"application_id": req.ApplicationID,
			"db_query_ms":    qMs,
			"sql":            slowSQL,
			"loan_amount":    req.Amount,
		})
		if err != nil {
			emitJSONLog(ctx, "ERROR", "Slow query error", map[string]any{"error": err.Error()})
		}
	} else {
		fastSQL := `SELECT COUNT(*), COALESCE(AVG(apr), 0) FROM loan_applications WHERE applicant_name = $1`
		qMs, _ := runTracedDBQuery(ctx, "postgres.select_applicant_history", "SELECT", "loan_applications", fastSQL, req.SimulateMode, func(qCtx context.Context) error {
			var count int
			var avgAPR float64
			return db.QueryRowContext(qCtx, fastSQL, req.ApplicantName).Scan(&count, &avgAPR)
		})
		totalDBMs += qMs
	}

	// 2. Insert initial record into PostgreSQL (triggers constraint violation when simulate_mode == "db_error" or amount == 666)
	insertAmount := req.Amount
	if req.SimulateMode == "db_error" {
		insertAmount = 666 // Violates CHECK (amount > 0 AND amount <> 666)
	}
	insertSQL := `INSERT INTO loan_applications (application_id, applicant_name, loan_type, amount, term_months, annual_income, credit_score, decision, trace_id) VALUES ($1, $2, $3, $4, $5, $6, $7, 'PENDING', $8)`
	insMs, err := runTracedDBQuery(ctx, "postgres.insert_loan_application", "INSERT", "loan_applications", insertSQL, req.SimulateMode, func(qCtx context.Context) error {
		_, execErr := db.ExecContext(qCtx, insertSQL,
			req.ApplicationID, req.ApplicantName, req.LoanType, insertAmount,
			req.TermMonths, req.AnnualIncome, req.CreditScore, traceID,
		)
		return execErr
	})
	totalDBMs += insMs

	if err != nil {
		durationSec := time.Since(start).Seconds()
		underwritingDurHist.Record(ctx, durationSec,
			metric.WithAttributes(
				attribute.String("loan_type", req.LoanType),
				attribute.String("decision", "DB_ERROR"),
			),
		)
		decisionsCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("loan_type", req.LoanType),
				attribute.String("decision", "DB_ERROR"),
			),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		emitJSONLog(ctx, "ERROR", "PostgreSQL constraint violation while persisting loan application", map[string]any{
			"application_id": req.ApplicationID,
			"amount":         insertAmount,
			"db_query_ms":    insMs,
			"sql_error":      err.Error(),
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(LoanResponse{
			ApplicationID: req.ApplicationID,
			Decision:      "DB_ERROR",
			Reason:        fmt.Sprintf("PostgreSQL Error: %v", err),
			TraceID:       traceID,
			DBQueryMs:     totalDBMs,
			ProcessingMs:  time.Since(start).Milliseconds(),
		})
		return
	}

	// 3. Call Python loan-underwriting service over HTTP with W3C TraceContext propagation
	bodyBytes, _ := json.Marshal(req)
	uwCtx, uwSpan := tracer.Start(ctx, "HTTP POST loan-underwriting/evaluate",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("peer.service", "loan-underwriting"),
			attribute.String("server.address", "loan-underwriting"),
			attribute.String("http.method", "POST"),
			attribute.String("http.url", underwritingURL+"/evaluate"),
		),
	)
	uwReq, err := http.NewRequestWithContext(uwCtx, http.MethodPost, underwritingURL+"/evaluate", bytes.NewReader(bodyBytes))
	if err != nil {
		uwSpan.RecordError(err)
		uwSpan.End()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uwReq.Header.Set("Content-Type", "application/json")
	otel.GetTextMapPropagator().Inject(uwCtx, propagation.HeaderCarrier(uwReq.Header))

	uwResp, err := httpClient.Do(uwReq)
	if err != nil {
		uwSpan.RecordError(err)
		uwSpan.SetStatus(codes.Error, err.Error())
		uwSpan.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		emitJSONLog(ctx, "ERROR", "Failed calling loan-underwriting service", map[string]any{
			"application_id": req.ApplicationID,
			"error":          err.Error(),
		})
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	uwSpan.SetAttributes(attribute.Int("http.status_code", uwResp.StatusCode))
	if uwResp.StatusCode >= 400 {
		uwSpan.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", uwResp.StatusCode))
	}
	uwSpan.End()
	defer uwResp.Body.Close()

	respBytes, _ := io.ReadAll(uwResp.Body)
	var uwData struct {
		ApplicationID    string  `json:"application_id"`
		Decision         string  `json:"decision"`
		APR              float64 `json:"apr"`
		MonthlyPayment   float64 `json:"monthly_payment"`
		DTIRatio         float64 `json:"dti_ratio"`
		RiskScore        float64 `json:"risk_score"`
		Reason           string  `json:"reason"`
		UnderwritingDBMs float64 `json:"underwriting_db_ms"`
	}
	_ = json.Unmarshal(respBytes, &uwData)
	totalDBMs += uwData.UnderwritingDBMs

	// 4. Update final underwriting decision in PostgreSQL
	updateSQL := `UPDATE loan_applications SET decision = $1, apr = $2, monthly_payment = $3, dti_ratio = $4, risk_score = $5, db_query_ms = $6 WHERE application_id = $7`
	updMs, _ := runTracedDBQuery(ctx, "postgres.update_loan_decision", "UPDATE", "loan_applications", updateSQL, req.SimulateMode, func(qCtx context.Context) error {
		_, execErr := db.ExecContext(qCtx, updateSQL,
			uwData.Decision, uwData.APR, uwData.MonthlyPayment, uwData.DTIRatio, uwData.RiskScore, totalDBMs, req.ApplicationID,
		)
		return execErr
	})
	totalDBMs += updMs

	durationSec := time.Since(start).Seconds()
	underwritingDurHist.Record(ctx, durationSec,
		metric.WithAttributes(
			attribute.String("loan_type", req.LoanType),
			attribute.String("decision", uwData.Decision),
		),
	)
	decisionsCounter.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("loan_type", req.LoanType),
			attribute.String("decision", uwData.Decision),
		),
	)

	if uwResp.StatusCode >= 400 {
		span.SetStatus(codes.Error, uwData.Reason)
	}

	finalResp := LoanResponse{
		ApplicationID:  req.ApplicationID,
		ApplicantName:  req.ApplicantName,
		LoanType:       req.LoanType,
		Amount:         req.Amount,
		CreditScore:    req.CreditScore,
		Decision:       uwData.Decision,
		APR:            uwData.APR,
		MonthlyPayment: uwData.MonthlyPayment,
		DTIRatio:       uwData.DTIRatio,
		RiskScore:      uwData.RiskScore,
		Reason:         uwData.Reason,
		TraceID:        traceID,
		DBQueryMs:      float64(int(totalDBMs*100)) / 100,
		ProcessingMs:   time.Since(start).Milliseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(uwResp.StatusCode)
	_ = json.NewEncoder(w).Encode(finalResp)
}

func handleListLoans(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	selectSQL := `SELECT application_id, applicant_name, loan_type, amount, credit_score, decision, apr, db_query_ms, trace_id, created_at FROM loan_applications ORDER BY id DESC LIMIT 15`

	var items []LoanResponse
	qMs, err := runTracedDBQuery(ctx, "postgres.list_recent_loans", "SELECT", "loan_applications", selectSQL, "list", func(qCtx context.Context) error {
		rows, qErr := db.QueryContext(qCtx, selectSQL)
		if qErr != nil {
			return qErr
		}
		defer rows.Close()
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
				return scanErr
			}
			item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
			items = append(items, item)
		}
		return nil
	})

	if err != nil {
		emitJSONLog(ctx, "ERROR", "Failed listing loan applications from PostgreSQL", map[string]any{"error": err.Error()})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items":       items,
		"db_query_ms": qMs,
	})
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	underwritingURL = os.Getenv("UNDERWRITING_URL")
	if underwritingURL == "" {
		underwritingURL = "http://loan-underwriting.loan-app.svc.cluster.local:8082"
	}

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	shutdown, err := initTelemetry(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize OpenTelemetry: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(shutCtx)
	}()

	if err := initDatabase(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize PostgreSQL: %v\n", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/underwrite", otelhttp.NewHandler(http.HandlerFunc(handleUnderwrite), "POST /underwrite"))
	mux.Handle("/loans", otelhttp.NewHandler(http.HandlerFunc(handleListLoans), "GET /loans"))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"loan-backend"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		fmt.Printf("loan-backend listening on :%s\n", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "HTTP server error: %v\n", err)
		}
	}()

	<-ctx.Done()
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutCtx)
}
