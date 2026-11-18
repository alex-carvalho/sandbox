package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
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

var (
	tracer              trace.Tracer
	otelLogger          otellog.Logger
	decisionsCounter    metric.Int64Counter
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

	// Business KPIs (metrics that automated APM cannot infer on its own)
	decisionsCounter, err = meter.Int64Counter(
		"loan_decisions",
		metric.WithDescription("Total loan underwriting decisions by decision and loan_type"),
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
