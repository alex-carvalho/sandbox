'use strict';

const { trace, context, propagation, SpanKind, SpanStatusCode } = require('@opentelemetry/api');
const { W3CTraceContextPropagator, W3CBaggagePropagator, CompositePropagator } = require('@opentelemetry/core');
const { NodeSDK } = require('@opentelemetry/sdk-node');
const { resourceFromAttributes } = require('@opentelemetry/resources');
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

const resource = resourceFromAttributes({
  'service.name': SERVICE_NAME,
  'service.namespace': 'loan-app',
  'service.version': '2.0.0',
  'telemetry.sdk.language': 'nodejs',
});

const traceExporter = new OTLPTraceExporter({ url: endpoint });
const metricExporter = new OTLPMetricExporter({ url: endpoint });
const logExporter = new OTLPLogExporter({ url: endpoint });

const loggerProvider = new LoggerProvider({
  resource,
  processors: [new BatchLogRecordProcessor(logExporter)],
});
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

// Automatically instrument global fetch() calls for W3C TraceContext propagation
// and outbound HTTP CLIENT spans without manual fetch wrappers in application code.
const originalFetch = globalThis.fetch;
if (typeof originalFetch === 'function') {
  const tracer = trace.getTracer(SERVICE_NAME);
  globalThis.fetch = async function instrumentedFetch(input, options = {}) {
    const rawUrl = typeof input === 'string' ? input : (input && input.url ? input.url : String(input));
    if (rawUrl.includes('/api/load') || rawUrl.includes('/healthz')) {
      return originalFetch.call(this, input, options);
    }

    let urlObj;
    try {
      urlObj = new URL(rawUrl, 'http://localhost');
    } catch (_err) {
      return originalFetch.call(this, input, options);
    }

    const method = (options && options.method) || (typeof input === 'object' && input && input.method) || 'GET';
    const hostParts = urlObj.hostname.split('.');
    const peerService = hostParts[0] || urlObj.hostname;

    return tracer.startActiveSpan(
      `HTTP ${method.toUpperCase()} ${urlObj.pathname}`,
      {
        kind: SpanKind.CLIENT,
        attributes: {
          'http.method': method.toUpperCase(),
          'http.url': rawUrl,
          'peer.service': peerService,
          'server.address': peerService,
        },
      },
      async (clientSpan) => {
        const headers = new Headers(
          (options && options.headers) || (typeof input === 'object' && input && input.headers) || {}
        );
        const carrier = {};
        propagation.inject(context.active(), carrier);
        for (const [k, v] of Object.entries(carrier)) {
          headers.set(k, v);
        }

        try {
          const response = await originalFetch.call(this, input, {
            ...options,
            headers,
          });
          clientSpan.setAttribute('http.status_code', response.status);
          if (response.status >= 400) {
            clientSpan.setStatus({
              code: SpanStatusCode.ERROR,
              message: `HTTP ${response.status}`,
            });
          }
          return response;
        } catch (err) {
          clientSpan.recordException(err);
          clientSpan.setStatus({ code: SpanStatusCode.ERROR, message: err.message });
          throw err;
        } finally {
          clientSpan.end();
        }
      }
    );
  };
}

process.on('SIGTERM', () => {
  Promise.all([sdk.shutdown(), loggerProvider.shutdown()])
    .finally(() => process.exit(0));
});
