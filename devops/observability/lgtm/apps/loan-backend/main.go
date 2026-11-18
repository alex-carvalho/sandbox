package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	underwritingURL = os.Getenv("UNDERWRITING_URL")
	if underwritingURL == "" {
		underwritingURL = "http://loan-underwriting.loan-app.svc.cluster.local:8082"
	}

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
				return "HTTP " + r.Method + " " + r.URL.Path
			}),
			otelhttp.WithClientTrace(nil),
		),
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
	loanRepo = NewLoanRepository(db)

	mux := http.NewServeMux()
	mux.Handle("/underwrite", otelhttp.NewHandler(http.HandlerFunc(handleUnderwrite), "POST /underwrite"))
	mux.Handle("/loans", otelhttp.NewHandler(http.HandlerFunc(handleListLoans), "GET /loans"))
	mux.HandleFunc("/healthz", handleHealthz)

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
