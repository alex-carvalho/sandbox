import json
import logging
import math
import os
import random
import time
from typing import Optional

import psycopg2
from fastapi import FastAPI
from fastapi.responses import JSONResponse
from pydantic import BaseModel

from opentelemetry import metrics, trace
from opentelemetry.exporter.otlp.proto.grpc._log_exporter import OTLPLogExporter
from opentelemetry.exporter.otlp.proto.grpc.metric_exporter import OTLPMetricExporter
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.psycopg2 import Psycopg2Instrumentor
from opentelemetry.sdk._logs import LoggerProvider, LoggingHandler
from opentelemetry.sdk._logs.export import BatchLogRecordProcessor
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.trace import SpanKind, Status, StatusCode

SERVICE_NAME = "loan-underwriting"
OTLP_ENDPOINT = os.getenv(
    "OTEL_EXPORTER_OTLP_ENDPOINT", "http://alloy.observability.svc.cluster.local:4317"
)
DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgresql://loanuser:loanpass@postgres.loan-app.svc.cluster.local:5432/loandb",
)

# ---------------------------------------------------------------------------
# OpenTelemetry Initialization
# ---------------------------------------------------------------------------
resource = Resource.create(
    {
        "service.name": SERVICE_NAME,
        "service.namespace": "loan-app",
        "service.version": "2.0.0",
    }
)

tracer_provider = TracerProvider(resource=resource)
tracer_provider.add_span_processor(
    BatchSpanProcessor(OTLPSpanExporter(endpoint=OTLP_ENDPOINT, insecure=True))
)
trace.set_tracer_provider(tracer_provider)
tracer = trace.get_tracer(SERVICE_NAME)

metric_reader = PeriodicExportingMetricReader(
    OTLPMetricExporter(endpoint=OTLP_ENDPOINT, insecure=True),
    export_interval_millis=5000,
)
meter_provider = MeterProvider(resource=resource, metric_readers=[metric_reader])
metrics.set_meter_provider(meter_provider)
meter = metrics.get_meter(SERVICE_NAME)

underwriting_eval_duration = meter.create_histogram(
    name="loan_underwriting_service_duration",
    description="Duration of loan-underwriting risk evaluation in seconds",
    unit="s",
)
db_query_duration = meter.create_histogram(
    name="loan_db_query_duration",
    description="PostgreSQL SQL query execution duration in seconds",
    unit="s",
)
underwriting_errors_counter = meter.create_counter(
    name="loan_underwriting_errors",
    description="Total underwriting errors by error_type and loan_type",
    unit="1",
)

logger_provider = LoggerProvider(resource=resource)
logger_provider.add_log_record_processor(
    BatchLogRecordProcessor(OTLPLogExporter(endpoint=OTLP_ENDPOINT, insecure=True))
)
otel_handler = LoggingHandler(level=logging.INFO, logger_provider=logger_provider)
logger = logging.getLogger(SERVICE_NAME)
logger.setLevel(logging.INFO)
logger.addHandler(otel_handler)

# Instrument psycopg2 so PostgreSQL queries appear as child spans automatically
Psycopg2Instrumentor().instrument(enable_commenter=True)


def emit_json_log(level: str, msg: str, extra_fields: dict):
    span_ctx = trace.get_current_span().get_span_context()
    trace_id = format(span_ctx.trace_id, "032x") if span_ctx.is_valid else ""
    span_id = format(span_ctx.span_id, "016x") if span_ctx.is_valid else ""

    payload = {
        "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "level": level,
        "service_name": SERVICE_NAME,
        "msg": msg,
        "trace_id": trace_id,
        "span_id": span_id,
        **extra_fields,
    }
    line = json.dumps(payload)
    print(line, flush=True)
    if level == "ERROR":
        logger.error(line)
    elif level == "WARN":
        logger.warning(line)
    else:
        logger.info(line)


class UnderwriteInput(BaseModel):
    application_id: str
    applicant_name: str
    loan_type: str
    amount: float
    term_months: int
    annual_income: float
    credit_score: int
    simulate_mode: Optional[str] = "normal"


app = FastAPI(title="Loan Underwriting & Risk Engine", version="2.0.0")
FastAPIInstrumentor.instrument_app(app)


@app.get("/healthz")
def healthz():
    return {"status": "ok", "service": SERVICE_NAME}


def check_fraud_blacklist_in_db(req: UnderwriteInput) -> tuple[bool, str, float]:
    """Queries PostgreSQL fraud_blacklist table and records DB query latency."""
    db_start = time.perf_counter()
    with tracer.start_as_current_span(
        "postgres.fraud_blacklist_check",
        kind=SpanKind.CLIENT,
        attributes={
            "db.system": "postgresql",
            "db.name": "loandb",
            "db.operation": "SELECT",
            "db.sql.table": "fraud_blacklist",
            "peer.service": "postgresql",
            "loan.application_id": req.application_id,
        },
    ) as span:
        try:
            conn = psycopg2.connect(DATABASE_URL)
            with conn.cursor() as cur:
                query = """
                    SELECT pattern, reason, risk_level
                    FROM fraud_blacklist
                    WHERE LOWER(%s) LIKE '%%' || LOWER(pattern) || '%%'
                    LIMIT 1;
                """
                span.set_attribute("db.statement", query.strip())
                cur.execute(query, (req.applicant_name,))
                row = cur.fetchone()
            conn.close()
        except Exception as exc:
            span.record_exception(exc)
            emit_json_log(
                "WARN",
                f"PostgreSQL fraud check fallback: {exc}",
                {"application_id": req.application_id},
            )
            row = None

        elapsed = time.perf_counter() - db_start
        db_query_duration.record(
            elapsed,
            {
                "service": SERVICE_NAME,
                "db_operation": "SELECT",
                "table": "fraud_blacklist",
                "scenario": req.simulate_mode or "normal",
            },
        )

        is_fraud = (
            row is not None
            or req.simulate_mode == "fraud_error"
            or "fraud" in req.applicant_name.lower()
            or "blacklist" in req.applicant_name.lower()
        )
        if is_fraud:
            reason = (
                row[1]
                if row
                else "Application matched high-risk sanctions/fraud watchlist"
            )
            err = PermissionError(reason)
            span.record_exception(err)
            span.set_status(Status(StatusCode.ERROR, reason))
            return True, reason, elapsed * 1000.0

        return False, "", elapsed * 1000.0


@app.post("/evaluate")
def evaluate_loan(req: UnderwriteInput):
    start_time = time.perf_counter()
    current_span = trace.get_current_span()
    trace_id = (
        format(current_span.get_span_context().trace_id, "032x")
        if current_span.get_span_context().is_valid
        else ""
    )

    current_span.set_attributes(
        {
            "loan.application_id": req.application_id,
            "loan.type": req.loan_type,
            "loan.amount": req.amount,
            "loan.credit_score": req.credit_score,
            "loan.simulate_mode": req.simulate_mode or "normal",
        }
    )

    emit_json_log(
        "INFO",
        "Underwriting service starting fraud and credit evaluation",
        {
            "application_id": req.application_id,
            "loan_type": req.loan_type,
            "amount": req.amount,
            "credit_score": req.credit_score,
            "simulate_mode": req.simulate_mode,
        },
    )

    # 1. PostgreSQL Fraud Blacklist Check
    is_fraud, fraud_reason, fraud_db_ms = check_fraud_blacklist_in_db(req)
    if is_fraud:
        duration = time.perf_counter() - start_time
        underwriting_eval_duration.record(
            duration,
            {"loan_type": req.loan_type, "decision": "FRAUD_BLOCKED", "scenario": "fraud_error"},
        )
        underwriting_errors_counter.add(
            1, {"error_type": "fraud_blacklist", "loan_type": req.loan_type}
        )
        current_span.set_status(Status(StatusCode.ERROR, fraud_reason))
        emit_json_log(
            "ERROR",
            "SECURITY ALERT: Loan application blocked by fraud detection engine",
            {
                "application_id": req.application_id,
                "reason": fraud_reason,
                "db_query_ms": round(fraud_db_ms, 2),
            },
        )
        return JSONResponse(
            status_code=403,
            content={
                "application_id": req.application_id,
                "decision": "FRAUD_BLOCKED",
                "apr": 0.0,
                "monthly_payment": 0.0,
                "dti_ratio": 0.0,
                "risk_score": 99.9,
                "reason": fraud_reason,
                "trace_id": trace_id,
                "underwriting_db_ms": round(fraud_db_ms, 2),
            },
        )

    # 2. External Credit Bureau Check
    with tracer.start_as_current_span(
        "credit-bureau.external_lookup",
        kind=SpanKind.CLIENT,
        attributes={
            "peer.service": "experian-credit-bureau",
            "loan.credit_score": req.credit_score,
            "loan.simulate_mode": req.simulate_mode or "normal",
        },
    ) as bureau_span:
        if req.simulate_mode == "bureau_timeout" or req.credit_score < 450:
            delay_s = random.uniform(0.55, 0.95)
            time.sleep(delay_s)
            err_msg = (
                f"Upstream Credit Bureau 504 Gateway Timeout after {int(delay_s*1000)}ms "
                f"(credit_score={req.credit_score})"
            )
            bureau_span.record_exception(TimeoutError(err_msg))
            bureau_span.set_status(Status(StatusCode.ERROR, err_msg))
            current_span.set_status(Status(StatusCode.ERROR, err_msg))

            duration = time.perf_counter() - start_time
            underwriting_eval_duration.record(
                duration,
                {
                    "loan_type": req.loan_type,
                    "decision": "BUREAU_TIMEOUT",
                    "scenario": "bureau_timeout",
                },
            )
            underwriting_errors_counter.add(
                1, {"error_type": "bureau_timeout", "loan_type": req.loan_type}
            )
            emit_json_log(
                "ERROR",
                "Credit bureau verification failed with gateway timeout",
                {
                    "application_id": req.application_id,
                    "credit_score": req.credit_score,
                    "bureau_latency_ms": round(delay_s * 1000, 2),
                    "error": err_msg,
                },
            )
            return JSONResponse(
                status_code=504,
                content={
                    "application_id": req.application_id,
                    "decision": "BUREAU_TIMEOUT",
                    "apr": 0.0,
                    "monthly_payment": 0.0,
                    "dti_ratio": 0.0,
                    "risk_score": 0.0,
                    "reason": err_msg,
                    "trace_id": trace_id,
                    "underwriting_db_ms": round(fraud_db_ms, 2),
                },
            )

        bureau_delay = random.uniform(0.03, 0.14)
        time.sleep(bureau_delay)
        bureau_span.set_attribute("bureau.response_ms", round(bureau_delay * 1000, 2))

    # 3. Actuarial Risk Scoring & Pricing Calculation
    with tracer.start_as_current_span("underwriting.actuarial_pricing") as price_span:
        time.sleep(random.uniform(0.015, 0.045))
        base_apr_map = {"home": 5.25, "auto": 6.10, "business": 7.75, "personal": 9.50}
        base_apr = base_apr_map.get(req.loan_type, 7.0)

        if req.credit_score >= 760:
            base_apr -= 0.75
        elif req.credit_score >= 700:
            base_apr += 0.25
        elif req.credit_score >= 640:
            base_apr += 1.85
        else:
            base_apr += 4.50

        monthly_rate = (base_apr / 100.0) / 12.0
        term = max(req.term_months, 12)
        monthly_payment = (req.amount * monthly_rate) / (
            1.0 - math.pow(1.0 + monthly_rate, -term)
        )
        monthly_income = max(req.annual_income / 12.0, 1.0)
        dti = (monthly_payment / monthly_income) * 100.0

        risk_score = max(
            1.0,
            min(
                99.0,
                round(
                    ((850 - req.credit_score) / 5.5) * 0.65 + min(dti, 100.0) * 0.35, 1
                ),
            ),
        )

        decision = "APPROVED"
        reason = "Meets automated credit score, fraud check, and DTI underwriting thresholds"
        if req.credit_score < 580 or dti > 50.0:
            decision = "REJECTED"
            reason = f"High actuarial risk profile (CreditScore={req.credit_score}, DTI={dti:.1f}%, RiskScore={risk_score})"
        elif req.credit_score < 650 or dti > 36.0 or req.amount > 500000:
            decision = "MANUAL_REVIEW"
            reason = f"Routed to senior underwriter queue (CreditScore={req.credit_score}, DTI={dti:.1f}%, Amount=${req.amount:,.0f})"

        price_span.set_attributes(
            {
                "loan.decision": decision,
                "loan.apr": round(base_apr, 2),
                "loan.dti_ratio": round(dti, 2),
                "loan.risk_score": risk_score,
            }
        )

    duration = time.perf_counter() - start_time
    underwriting_eval_duration.record(
        duration,
        {
            "loan_type": req.loan_type,
            "decision": decision,
            "scenario": req.simulate_mode or "normal",
        },
    )

    emit_json_log(
        "INFO",
        "Underwriting evaluation completed",
        {
            "application_id": req.application_id,
            "loan_type": req.loan_type,
            "decision": decision,
            "apr": round(base_apr, 2),
            "dti_ratio": round(dti, 2),
            "risk_score": risk_score,
            "duration_ms": round(duration * 1000, 2),
        },
    )

    return {
        "application_id": req.application_id,
        "decision": decision,
        "apr": round(base_apr, 2),
        "monthly_payment": round(monthly_payment, 2),
        "dti_ratio": round(dti, 2),
        "risk_score": risk_score,
        "reason": reason,
        "trace_id": trace_id,
        "underwriting_db_ms": round(fraud_db_ms, 2),
    }
