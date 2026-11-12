'use strict';

const { propagation } = require('@opentelemetry/api');
const { W3CTraceContextPropagator, W3CBaggagePropagator, CompositePropagator } = require('@opentelemetry/core');
const { NodeSDK } = require('@opentelemetry/sdk-node');
const { Resource } = require('@opentelemetry/resources');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc');
const { OTLPMetricExporter } = require('@opentelemetry/exporter-metrics-otlp-grpc');
const { OTLPLogExporter } = require('@opentelemetry/exporter-logs-otlp-grpc');
const { PeriodicExportingMetricReader } = require('@opentelemetry/sdk-metrics');
const { LoggerProvider, BatchLogRecordProcessor } = require('@opentelemetry/sdk-logs');
const { logs } = require('@opentelemetry/api-logs');
const { HttpInstrumentation } = require('@opentelemetry/instrumentation-http');
const { ExpressInstrumentation } = require('@opentelemetry/instrumentation-express');

const SERVICE_NAME = process.env.OTEL_SERVICE_NAME || 'loan-frontend';
const rawEndpoint = process.env.OTEL_EXPORTER_OTLP_ENDPOINT || 'http://alloy.observability.svc.cluster.local:4317';
const endpoint = rawEndpoint.startsWith('http') ? rawEndpoint : `http://${rawEndpoint}`;

// Register global W3C TraceContext propagator explicitly
const textMapPropagator = new CompositePropagator({
  propagators: [new W3CTraceContextPropagator(), new W3CBaggagePropagator()],
});
propagation.setGlobalPropagator(textMapPropagator);

const resource = new Resource({
  'service.name': SERVICE_NAME,
  'service.namespace': 'loan-app',
  'service.version': '2.0.0',
  'telemetry.sdk.language': 'nodejs',
});

const traceExporter = new OTLPTraceExporter({ url: endpoint });
const metricExporter = new OTLPMetricExporter({ url: endpoint });
const logExporter = new OTLPLogExporter({ url: endpoint });

const loggerProvider = new LoggerProvider({ resource });
loggerProvider.addLogRecordProcessor(new BatchLogRecordProcessor(logExporter));
logs.setGlobalLoggerProvider(loggerProvider);

const sdk = new NodeSDK({
  resource,
  traceExporter,
  textMapPropagator,
  metricReader: new PeriodicExportingMetricReader({
    exporter: metricExporter,
    exportIntervalMillis: 5000,
  }),
  instrumentations: [
    new HttpInstrumentation({
      ignoreIncomingRequestHook: (req) =>
        req.url === '/healthz' || (req.url && req.url.startsWith('/api/load')),
    }),
    new ExpressInstrumentation(),
  ],
});

sdk.start();

process.on('SIGTERM', () => {
  Promise.all([sdk.shutdown(), loggerProvider.shutdown()])
    .finally(() => process.exit(0));
});
