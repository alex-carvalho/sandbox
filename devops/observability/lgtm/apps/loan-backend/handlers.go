package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	httpClient      *http.Client
	underwritingURL string
)

func callUnderwritingService(ctx context.Context, req LoanRequest) (*UnderwritingResult, int, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	uwReq, err := http.NewRequestWithContext(ctx, http.MethodPost, underwritingURL+"/evaluate", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	uwReq.Header.Set("Content-Type", "application/json")

	// otelhttp.NewTransport automatically creates the CLIENT span & injects W3C traceparent
	uwResp, err := httpClient.Do(uwReq)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}
	defer uwResp.Body.Close()

	respBytes, err := io.ReadAll(uwResp.Body)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}

	var uwData UnderwritingResult
	if err := json.Unmarshal(respBytes, &uwData); err != nil {
		return nil, uwResp.StatusCode, err
	}
	return &uwData, uwResp.StatusCode, nil
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

	// Business KPI metric
	requestedAmountHist.Record(ctx, req.Amount,
		metric.WithAttributes(attribute.String("loan_type", req.LoanType)),
	)

	var totalDBMs float64

	// 1. Query historical exposure in PostgreSQL (otelsql automatically traces SQL queries)
	if req.SimulateMode == "slow_db" || req.Amount > 900000 {
		qMs, err := loanRepo.AnalyzeSlowPortfolioExposure(ctx, req.LoanType)
		totalDBMs += qMs
		emitJSONLog(ctx, "WARN", "Slow PostgreSQL portfolio exposure query detected (>800ms)", map[string]any{
			"application_id": req.ApplicationID,
			"db_query_ms":    qMs,
			"loan_amount":    req.Amount,
		})
		if err != nil {
			emitJSONLog(ctx, "ERROR", "Slow query error", map[string]any{"error": err.Error()})
		}
	} else {
		qMs, _ := loanRepo.CheckApplicantHistory(ctx, req.ApplicantName)
		totalDBMs += qMs
	}

	// 2. Persist initial record (triggers constraint violation when simulate_mode == "db_error" or amount == 666)
	insertAmount := req.Amount
	if req.SimulateMode == "db_error" {
		insertAmount = 666
	}
	insMs, err := loanRepo.InsertPendingLoan(ctx, req, insertAmount, traceID)
	totalDBMs += insMs

	if err != nil {
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

	// 3. Call Python loan-underwriting service over HTTP
	uwData, uwStatusCode, err := callUnderwritingService(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		emitJSONLog(ctx, "ERROR", "Failed calling loan-underwriting service", map[string]any{
			"application_id": req.ApplicationID,
			"error":          err.Error(),
		})
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	totalDBMs += uwData.UnderwritingDBMs

	// 4. Update final underwriting decision in PostgreSQL
	updMs, _ := loanRepo.UpdateLoanDecision(ctx, req.ApplicationID, uwData, totalDBMs)
	totalDBMs += updMs

	// Business KPI metric
	decisionsCounter.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("loan_type", req.LoanType),
			attribute.String("decision", uwData.Decision),
		),
	)

	if uwStatusCode >= 400 {
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
	w.WriteHeader(uwStatusCode)
	_ = json.NewEncoder(w).Encode(finalResp)
}

func handleListLoans(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	items, qMs, err := loanRepo.ListRecentLoans(ctx, 15)
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

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","service":"loan-backend"}`))
}
