'use strict';

const crypto = require('crypto');
const express = require('express');
const {
  trace,
  context,
  propagation,
  SpanKind,
  SpanStatusCode,
  metrics,
} = require('@opentelemetry/api');
const { logs, SeverityNumber } = require('@opentelemetry/api-logs');

const SERVICE_NAME = 'loan-frontend';
const PORT = process.env.PORT || 8080;
const BACKEND_URL = process.env.LOAN_BACKEND_URL || 'http://loan-backend.loan-app.svc.cluster.local:8081';

const tracer = trace.getTracer(SERVICE_NAME);
const meter = metrics.getMeter(SERVICE_NAME);
const otelLogger = logs.getLogger(SERVICE_NAME);

// Traffic Generator State controlled from the Frontend UI
const loadState = {
  enabled: true,
  interval_ms: 2000,
  total_sent: 0,
  last_scenario: 'normal',
  last_decision: 'APPROVED',
  updated_at: new Date().toISOString(),
};

// OpenTelemetry Metrics
const requestsCounter = meter.createCounter('loan_frontend_requests', {
  description: 'Total loan applications received by Node.js loan-frontend',
});

const requestDurationHist = meter.createHistogram('loan_frontend_request_duration', {
  description: 'End-to-end loan application request duration on Node.js frontend (seconds)',
  unit: 's',
  advice: {
    explicitBucketBoundaries: [0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 0.75, 1.0, 2.5, 5.0],
  },
});

// Helper that wraps Node global fetch() in an OpenTelemetry CLIENT span
// and explicitly injects W3C TraceContext (traceparent / tracestate) headers
// so loan-backend and loan-underwriting continue the exact same distributed trace.
async function tracedFetch(spanName, url, options = {}) {
  return tracer.startActiveSpan(
    spanName,
    {
      kind: SpanKind.CLIENT,
      attributes: {
        'http.method': options.method || 'GET',
        'http.url': url,
        'peer.service': 'loan-backend',
        'server.address': 'loan-backend',
      },
    },
    async (clientSpan) => {
      const headers = { ...(options.headers || {}) };
      propagation.inject(context.active(), headers);

      try {
        const response = await fetch(url, {
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
}

function emitJsonLog(level, msg, extraFields = {}) {
  const activeSpan = trace.getSpan(context.active());
  const spanCtx = activeSpan ? activeSpan.spanContext() : null;
  const traceId = spanCtx && spanCtx.traceId ? spanCtx.traceId : '';
  const spanId = spanCtx && spanCtx.spanId ? spanCtx.spanId : '';

  const payload = {
    timestamp: new Date().toISOString(),
    level,
    service_name: SERVICE_NAME,
    runtime: 'nodejs',
    msg,
    trace_id: traceId,
    span_id: spanId,
    ...extraFields,
  };

  const line = JSON.stringify(payload);
  console.log(line);

  const severityMap = {
    ERROR: SeverityNumber.ERROR,
    WARN: SeverityNumber.WARN,
    INFO: SeverityNumber.INFO,
  };

  otelLogger.emit({
    severityNumber: severityMap[level] || SeverityNumber.INFO,
    severityText: level,
    body: line,
    attributes: {
      trace_id: traceId,
      span_id: spanId,
      'service.name': SERVICE_NAME,
    },
    context: context.active(),
  });
}

const app = express();
app.use(express.json());

app.get('/healthz', (_req, res) => {
  res.json({ status: 'ok', service: SERVICE_NAME, runtime: 'nodejs' });
});

// ---------------------------------------------------------------------------
// Load Generator Control Endpoints (untraced control plane)
// ---------------------------------------------------------------------------
app.get('/api/load/status', (_req, res) => {
  res.json(loadState);
});

app.post('/api/load/control', (req, res) => {
  const { enabled, interval_ms, record_sent, last_scenario, last_decision } = req.body || {};
  if (typeof enabled === 'boolean') {
    loadState.enabled = enabled;
    loadState.updated_at = new Date().toISOString();
  }
  if (typeof interval_ms === 'number' && interval_ms >= 200 && interval_ms <= 30000) {
    loadState.interval_ms = interval_ms;
    loadState.updated_at = new Date().toISOString();
  }
  if (record_sent) {
    loadState.total_sent += 1;
    if (last_scenario) loadState.last_scenario = last_scenario;
    if (last_decision) loadState.last_decision = last_decision;
  }
  res.json(loadState);
});

// Proxy recent loans history from Go backend + PostgreSQL
app.get('/api/loans/history', async (_req, res) => {
  const start = performance.now();
  const activeSpan = trace.getSpan(context.active());
  try {
    const response = await tracedFetch(
      'HTTP GET loan-backend/loans',
      `${BACKEND_URL}/loans?limit=15`,
      { method: 'GET' }
    );
    const data = await response.json();
    emitJsonLog('INFO', 'Fetched recent loan applications from backend DB', {
      count: Array.isArray(data.items) ? data.items.length : 0,
      duration_ms: Math.round(performance.now() - start),
    });
    res.status(response.status).json(data);
  } catch (err) {
    if (activeSpan) {
      activeSpan.recordException(err);
      activeSpan.setStatus({ code: SpanStatusCode.ERROR, message: err.message });
    }
    emitJsonLog('ERROR', 'Failed to fetch loan history from backend', { error: err.message });
    res.status(502).json({ error: err.message });
  }
});

// Submit loan application -> Go Backend -> Python Underwriting + PostgreSQL
app.post('/api/loans/apply', async (req, res) => {
  const startTime = performance.now();
  const appId = `LN-${crypto.randomBytes(4).toString('hex').toUpperCase()}`;
  const activeSpan = trace.getSpan(context.active());
  const traceId = activeSpan ? activeSpan.spanContext().traceId : '';

  const {
    applicant_name = 'Alex Rivera',
    loan_type = 'home',
    amount = 250000,
    term_months = 360,
    annual_income = 145000,
    credit_score = 760,
    simulate_mode = 'normal',
  } = req.body || {};

  // Allow user input data to automatically trigger presentation scenarios
  let effectiveMode = simulate_mode;
  const lowerName = String(applicant_name).toLowerCase();
  if (effectiveMode === 'normal') {
    if (lowerName.includes('fraud') || lowerName.includes('blacklist')) {
      effectiveMode = 'fraud_error';
    } else if (Number(amount) === 666 || Number(amount) < 0) {
      effectiveMode = 'db_error';
    } else if (Number(amount) > 900000 || lowerName.includes('slow')) {
      effectiveMode = 'slow_db';
    } else if (Number(credit_score) < 450) {
      effectiveMode = 'bureau_timeout';
    }
  }

  if (activeSpan) {
    activeSpan.setAttributes({
      'loan.application_id': appId,
      'loan.type': loan_type,
      'loan.amount': Number(amount),
      'loan.credit_score': Number(credit_score),
      'loan.simulate_mode': effectiveMode,
    });
  }

  // Child span on frontend for input validation
  await tracer.startActiveSpan('frontend.validate_input', async (valSpan) => {
    valSpan.setAttributes({
      'loan.application_id': appId,
      'loan.effective_mode': effectiveMode,
    });
    await new Promise((r) => setTimeout(r, 8));
    valSpan.end();
  });

  emitJsonLog('INFO', 'Received loan application on Node.js frontend portal', {
    application_id: appId,
    loan_type,
    amount: Number(amount),
    credit_score: Number(credit_score),
    simulate_mode: effectiveMode,
  });

  const payload = {
    application_id: appId,
    applicant_name,
    loan_type,
    amount: Number(amount),
    term_months: Number(term_months),
    annual_income: Number(annual_income),
    credit_score: Number(credit_score),
    simulate_mode: effectiveMode,
  };

  try {
    const response = await tracedFetch(
      'HTTP POST loan-backend/underwrite',
      `${BACKEND_URL}/underwrite`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      }
    );

    const data = await response.json();
    data.trace_id = traceId || data.trace_id;

    const durationSec = (performance.now() - startTime) / 1000.0;
    const isError = response.status >= 400;
    const statusLabel = isError ? 'error' : 'success';

    requestsCounter.add(1, {
      loan_type,
      status: statusLabel,
      decision: data.decision || 'ERROR',
      scenario: effectiveMode,
    });
    requestDurationHist.record(durationSec, {
      loan_type,
      status: statusLabel,
      scenario: effectiveMode,
    });

    if (isError) {
      if (activeSpan) {
        activeSpan.setStatus({
          code: SpanStatusCode.ERROR,
          message: data.reason || `HTTP ${response.status}`,
        });
      }
      emitJsonLog('ERROR', 'Loan application failed during downstream processing', {
        application_id: appId,
        loan_type,
        scenario: effectiveMode,
        status_code: response.status,
        decision: data.decision,
        reason: data.reason,
      });
      return res.status(response.status).json(data);
    }

    emitJsonLog('INFO', 'Loan application processed successfully', {
      application_id: appId,
      loan_type,
      decision: data.decision,
      apr: data.apr,
      risk_score: data.risk_score,
      db_query_ms: data.db_query_ms,
      processing_ms: data.processing_ms,
    });

    return res.status(200).json(data);
  } catch (err) {
    const durationSec = (performance.now() - startTime) / 1000.0;
    if (activeSpan) {
      activeSpan.recordException(err);
      activeSpan.setStatus({ code: SpanStatusCode.ERROR, message: err.message });
    }
    requestsCounter.add(1, {
      loan_type,
      status: 'error',
      decision: 'ERROR',
      scenario: effectiveMode,
    });
    requestDurationHist.record(durationSec, {
      loan_type,
      status: 'error',
      scenario: effectiveMode,
    });
    emitJsonLog('ERROR', `Unhandled exception calling loan-backend: ${err.message}`, {
      application_id: appId,
      loan_type,
      scenario: effectiveMode,
    });
    return res.status(502).json({
      application_id: appId,
      decision: 'ERROR',
      reason: err.message,
      trace_id: traceId,
    });
  }
});

app.get('/', (_req, res) => {
  res.setHeader('Content-Type', 'text/html');
  res.send(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Loan Portal (Node.js) — Grafana LGTM Live Demo</title>
  <style>
    :root { --bg: #0f172a; --card: #1e293b; --text: #f8fafc; --accent: #f97316; --border: #334155; }
    body { font-family: system-ui, -apple-system, sans-serif; background: var(--bg); color: var(--text); margin: 0; padding: 1.75rem; }
    .container { max-width: 1100px; margin: 0 auto; }
    .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem; border-bottom: 1px solid var(--border); padding-bottom: 1rem; }
    .badge { background: #334155; color: #38bdf8; padding: 0.25rem 0.6rem; border-radius: 999px; font-size: 0.78rem; font-weight: 600; }
    .status-pill { display: inline-flex; align-items: center; gap: 0.4rem; padding: 0.3rem 0.75rem; border-radius: 999px; font-size: 0.8rem; font-weight: 700; }
    .status-running { background: rgba(16, 185, 129, 0.2); color: #34d399; border: 1px solid #10b981; }
    .status-paused { background: rgba(234, 179, 8, 0.2); color: #facc15; border: 1px solid #eab308; }
    .card { background: var(--card); border: 1px solid var(--border); border-radius: 12px; padding: 1.35rem; margin-bottom: 1.2rem; }
    .load-bar { display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 1rem; background: #0f172a; border: 1px solid var(--border); border-radius: 10px; padding: 0.9rem 1.2rem; margin-bottom: 1.2rem; }
    .grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 0.9rem; }
    label { display: block; font-size: 0.8rem; color: #94a3b8; margin-bottom: 0.3rem; }
    input, select { width: 100%; box-sizing: border-box; padding: 0.55rem; border-radius: 6px; border: 1px solid var(--border); background: #0f172a; color: white; }
    .btns { display: flex; gap: 0.65rem; margin-top: 1.15rem; flex-wrap: wrap; }
    button { cursor: pointer; border: none; border-radius: 6px; padding: 0.6rem 1rem; font-weight: 600; color: white; background: var(--accent); font-size: 0.85rem; }
    button.green { background: #10b981; }
    button.warn { background: #eab308; color: #0f172a; }
    button.purple { background: #8b5cf6; }
    button.danger { background: #ef4444; }
    button.dark { background: #475569; }
    pre { background: #020617; padding: 0.9rem; border-radius: 8px; overflow-x: auto; border: 1px solid var(--border); font-size: 0.83rem; }
    a { color: #38bdf8; text-decoration: none; font-weight: 600; }
    table { width: 100%; border-collapse: collapse; font-size: 0.83rem; margin-top: 0.5rem; }
    th, td { text-align: left; padding: 0.5rem; border-bottom: 1px solid var(--border); }
    th { color: #94a3b8; }
    .rules { font-size: 0.82rem; color: #cbd5e1; line-height: 1.5; margin-top: 0.75rem; background: #0f172a; padding: 0.75rem 1rem; border-radius: 8px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <div>
        <h1 style="margin:0;">Loan Application Portal <span class="badge">Frontend: Node.js &rarr; Backend: Go &rarr; Underwriting: Python &rarr; PostgreSQL</span></h1>
        <p style="margin:0.35rem 0 0;color:#94a3b8;">Grafana LGTM Stack (Loki, Grafana, Tempo, Mimir) + Grafana Alloy Live Presentation Demo</p>
      </div>
      <div>
        <a href="http://localhost:3000/d/loan-lgtm-poc" target="_blank">Open Grafana Dashboard &rarr;</a>
      </div>
    </div>

    <!-- Traffic Generator Live Control Bar -->
    <div class="load-bar">
      <div style="display:flex;align-items:center;gap:1rem;">
        <span id="load_badge" class="status-pill status-running">RUNNING</span>
        <div>
          <div style="font-weight:600;font-size:0.92rem;">Background Load Generator (<code>loan-traffic-gen</code>)</div>
          <div id="load_stats" style="font-size:0.8rem;color:#94a3b8;">Sent: 0 requests | Last scenario: normal</div>
        </div>
      </div>
      <div style="display:flex;align-items:center;gap:0.65rem;">
        <select id="load_interval" style="width:auto;" onchange="changeLoadSpeed()">
          <option value="4000">Slow (1 req / 4s)</option>
          <option value="2000" selected>Normal (1 req / 2s)</option>
          <option value="500">Fast Stress (2 req / sec)</option>
        </select>
        <button id="load_toggle_btn" class="warn" onclick="toggleLoad()">Pause Load Generator</button>
      </div>
    </div>

    <div class="card">
      <h3 style="margin-top:0;">1. Submit Loan Application or Trigger Presentation Scenario</h3>
      <div class="grid">
        <div>
          <label>Applicant Name (use "Fraud" to trigger blacklist error)</label>
          <input id="applicant_name" value="Alex Rivera" />
        </div>
        <div>
          <label>Loan Type</label>
          <select id="loan_type">
            <option value="home">Home Mortgage</option>
            <option value="auto">Auto Loan</option>
            <option value="personal">Personal Loan</option>
            <option value="business">Small Business Loan</option>
          </select>
        </div>
        <div>
          <label>Loan Amount ($ USD — &gt;900k = Slow DB, 666 = DB Error)</label>
          <input id="amount" type="number" value="250000" />
        </div>
        <div>
          <label>Term (Months)</label>
          <input id="term_months" type="number" value="360" />
        </div>
        <div>
          <label>Annual Income ($ USD)</label>
          <input id="annual_income" type="number" value="145000" />
        </div>
        <div>
          <label>Credit Score (300-850 — &lt;450 = Bureau Timeout 504)</label>
          <input id="credit_score" type="number" value="760" />
        </div>
      </div>

      <div class="btns">
        <button class="green" onclick="submitLoan('normal')">Submit Custom Form</button>
        <button class="warn" onclick="runPreset('slow_db')">Demo: Slow PostgreSQL Query (~900ms)</button>
        <button class="danger" onclick="runPreset('fraud_error')">Demo: Fraud Blacklist Error (403)</button>
        <button class="danger" onclick="runPreset('bureau_timeout')">Demo: Credit Bureau Timeout (504)</button>
        <button class="purple" onclick="runPreset('db_error')">Demo: PostgreSQL Constraint Error (500)</button>
      </div>

      <div class="rules">
        <strong>Input Data Triggers (also work directly from form or API):</strong>
        Applicant Name containing <code>Fraud</code> or <code>Blacklist</code> &rarr; <strong>403 Fraud Blocked</strong> |
        Credit Score <code>&lt; 450</code> &rarr; <strong>504 Bureau Timeout</strong> |
        Amount <code>&gt; $900,000</code> &rarr; <strong>Slow PostgreSQL Analytical Query</strong> |
        Amount <code>= $666</code> &rarr; <strong>500 PostgreSQL SQL Exception</strong>
      </div>
    </div>

    <div class="card">
      <h3 style="margin-top:0;">2. Underwriting Result & Distributed Trace Link</h3>
      <div id="trace_link" style="margin-bottom:0.65rem;"></div>
      <pre id="result">Click any button above to generate a multi-service trace across Node.js &rarr; Go &rarr; Python &rarr; PostgreSQL...</pre>
    </div>

    <div class="card">
      <div style="display:flex;justify-content:space-between;align-items:center;">
        <h3 style="margin:0;">3. Recent Loan Applications in PostgreSQL</h3>
        <button class="dark" onclick="loadHistory()">Refresh DB Table (SELECT query trace)</button>
      </div>
      <table>
        <thead>
          <tr>
            <th>App ID</th>
            <th>Applicant</th>
            <th>Type</th>
            <th>Amount</th>
            <th>Score</th>
            <th>Decision</th>
            <th>DB Query (ms)</th>
            <th>Trace in Tempo</th>
          </tr>
        </thead>
        <tbody id="history_body">
          <tr><td colspan="8">Loading from PostgreSQL...</td></tr>
        </tbody>
      </table>
    </div>
  </div>

  <script>
    let currentLoadEnabled = true;

    async function refreshLoadStatus() {
      try {
        const res = await fetch('/api/load/status');
        const st = await res.json();
        currentLoadEnabled = st.enabled;
        const badge = document.getElementById('load_badge');
        const btn = document.getElementById('load_toggle_btn');
        if (st.enabled) {
          badge.textContent = 'RUNNING';
          badge.className = 'status-pill status-running';
          btn.textContent = 'Pause Load Generator';
          btn.className = 'warn';
        } else {
          badge.textContent = 'PAUSED';
          badge.className = 'status-pill status-paused';
          btn.textContent = 'Start Load Generator';
          btn.className = 'green';
        }
        document.getElementById('load_interval').value = String(st.interval_ms || 2000);
        document.getElementById('load_stats').textContent =
          'Auto-generated: ' + st.total_sent + ' requests | Last scenario: ' + st.last_scenario + ' (' + st.last_decision + ')';
      } catch (e) {}
    }

    async function toggleLoad() {
      await fetch('/api/load/control', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ enabled: !currentLoadEnabled })
      });
      await refreshLoadStatus();
    }

    async function changeLoadSpeed() {
      const interval_ms = parseInt(document.getElementById('load_interval').value, 10);
      await fetch('/api/load/control', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ interval_ms })
      });
      await refreshLoadStatus();
    }

    function tempoLink(traceId) {
      if (!traceId) return '';
      const url = 'http://localhost:3000/explore?left=' + encodeURIComponent(JSON.stringify({
        datasource: 'tempo',
        queries: [{ refId: 'A', queryType: 'traceql', query: traceId }]
      }));
      return '<a href="' + url + '" target="_blank">' + traceId.slice(0, 12) + '... &rarr;</a>';
    }

    async function runPreset(mode) {
      if (mode === 'slow_db') {
        document.getElementById('applicant_name').value = 'Sofia Chen (VIP High Value)';
        document.getElementById('loan_type').value = 'home';
        document.getElementById('amount').value = '980000';
        document.getElementById('credit_score').value = '775';
      } else if (mode === 'fraud_error') {
        document.getElementById('applicant_name').value = 'Victor Fraud Blacklist';
        document.getElementById('loan_type').value = 'personal';
        document.getElementById('amount').value = '55000';
        document.getElementById('credit_score').value = '680';
      } else if (mode === 'bureau_timeout') {
        document.getElementById('applicant_name').value = 'Liam O\\'Connor';
        document.getElementById('loan_type').value = 'auto';
        document.getElementById('amount').value = '32000';
        document.getElementById('credit_score').value = '410';
      } else if (mode === 'db_error') {
        document.getElementById('applicant_name').value = 'SQL Constraint Trigger';
        document.getElementById('loan_type').value = 'business';
        document.getElementById('amount').value = '666';
        document.getElementById('credit_score').value = '710';
      }
      await submitLoan(mode);
    }

    async function submitLoan(mode) {
      const payload = {
        applicant_name: document.getElementById('applicant_name').value,
        loan_type: document.getElementById('loan_type').value,
        amount: parseFloat(document.getElementById('amount').value),
        term_months: parseInt(document.getElementById('term_months').value, 10),
        annual_income: parseFloat(document.getElementById('annual_income').value),
        credit_score: parseInt(document.getElementById('credit_score').value, 10),
        simulate_mode: mode
      };
      document.getElementById('result').textContent = 'Executing distributed request (Node.js -> Go -> Python -> PostgreSQL)...';
      const resp = await fetch('/api/loans/apply', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      const data = await resp.json();
      document.getElementById('result').textContent = JSON.stringify(data, null, 2);
      if (data.trace_id) {
        document.getElementById('trace_link').innerHTML =
          'Trace ID: <code>' + data.trace_id + '</code> &mdash; ' + tempoLink(data.trace_id);
      }
      await loadHistory();
    }

    async function loadHistory() {
      try {
        const resp = await fetch('/api/loans/history');
        const data = await resp.json();
        const rows = (data.items || []).map(item => (
          '<tr>' +
            '<td><code>' + item.application_id + '</code></td>' +
            '<td>' + item.applicant_name + '</td>' +
            '<td>' + item.loan_type + '</td>' +
            '<td>$' + Number(item.amount).toLocaleString() + '</td>' +
            '<td>' + item.credit_score + '</td>' +
            '<td><strong>' + item.decision + '</strong></td>' +
            '<td>' + (item.db_query_ms || 0) + ' ms</td>' +
            '<td>' + tempoLink(item.trace_id) + '</td>' +
          '</tr>'
        )).join('');
        document.getElementById('history_body').innerHTML = rows || '<tr><td colspan="8">No records yet</td></tr>';
      } catch (e) {
        document.getElementById('history_body').innerHTML = '<tr><td colspan="8">Error loading history</td></tr>';
      }
    }

    refreshLoadStatus();
    setInterval(refreshLoadStatus, 2000);
    loadHistory();
  </script>
</body>
</html>`);
});

app.listen(PORT, () => {
  console.log(`loan-frontend (Node.js) listening on :${PORT}`);
});
