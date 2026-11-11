const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const CHROME_PATH = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe";
const TEMP_DIR = "C:\\Users\\alexander\\AppData\\Local\\Temp";
const OUTPUT_DIR = "B:\\workgit\\n8nflow-guard\\docs\\images";

function wrapInTerminal(title, command, contentHtml, width = 1200, height = 750) {
  return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    background: #090d16;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    padding: 30px;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  }
  .window {
    width: 100%;
    max-width: ${width}px;
    background: #0f172a;
    border-radius: 12px;
    border: 1px solid #1e293b;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.7), 0 0 0 1px rgba(255, 255, 255, 0.05);
    overflow: hidden;
  }
  .header {
    background: #1e293b;
    padding: 12px 18px;
    display: flex;
    align-items: center;
    border-bottom: 1px solid #334155;
  }
  .dots {
    display: flex;
    gap: 8px;
  }
  .dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
  }
  .dot-red { background: #ef4444; }
  .dot-yellow { background: #eab308; }
  .dot-green { background: #22c55e; }
  .title {
    margin: 0 auto;
    color: #94a3b8;
    font-size: 13px;
    font-weight: 500;
    letter-spacing: 0.3px;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    padding-right: 50px;
  }
  .body {
    padding: 24px 28px;
    color: #f1f5f9;
    font-family: ui-monospace, SFMono-Regular, "Cascadia Code", "Source Code Pro", Menlo, Monaco, Consolas, monospace;
    font-size: 13.5px;
    line-height: 1.6;
    overflow-x: auto;
  }
  .prompt-line {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 18px;
    padding-bottom: 14px;
    border-bottom: 1px solid #1e293b;
  }
  .prompt-symbol { color: #00f5ff; font-weight: bold; }
  .prompt-cmd { color: #f8fafc; font-weight: 600; }
  .badge {
    display: inline-block;
    padding: 1px 7px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 0.5px;
  }
  .badge-critical { background: #dc2626; color: #ffffff; }
  .badge-high { background: #ea580c; color: #ffffff; }
  .badge-medium { background: #d97706; color: #ffffff; }
  .badge-low { background: #0284c7; color: #ffffff; }
  .badge-fixable { background: #16a34a; color: #ffffff; }
  .loc { color: #64748b; }
  .rule { color: #38bdf8; font-weight: 600; }
  .msg { color: #f1f5f9; }
  .hint { color: #94a3b8; padding-left: 20px; font-size: 12.5px; }
  .hint-accent { color: #38bdf8; }
  .suggested { color: #34d399; padding-left: 20px; font-size: 12.5px; }
  .diff-add { color: #4ade80; background: rgba(34, 197, 94, 0.1); }
  .diff-del { color: #f87171; background: rgba(239, 68, 68, 0.1); }
  .stat-card {
    background: #1e293b;
    border: 1px solid #334155;
    border-radius: 8px;
    padding: 14px 18px;
    margin-top: 20px;
    display: flex;
    justify-content: space-between;
  }
  .stat-item { text-align: center; }
  .stat-num { font-size: 20px; font-weight: 700; color: #00f5ff; }
  .stat-lbl { font-size: 11px; color: #94a3b8; text-transform: uppercase; letter-spacing: 0.5px; }
</style>
</head>
<body>
  <div class="window">
    <div class="header">
      <div class="dots">
        <div class="dot dot-red"></div>
        <div class="dot dot-yellow"></div>
        <div class="dot dot-green"></div>
      </div>
      <div class="title">${title}</div>
    </div>
    <div class="body">
      <div class="prompt-line">
        <span class="prompt-symbol">❯</span>
        <span class="prompt-cmd">${command}</span>
      </div>
      ${contentHtml}
    </div>
  </div>
</body>
</html>`;
}

function takeScreenshot(htmlContent, tempName, outName, width = 1200, height = 780) {
  const tempHtmlPath = path.join(TEMP_DIR, `${tempName}.html`);
  const tempPngPath = path.join(TEMP_DIR, `${tempName}.png`);
  const outPngPath = path.join(OUTPUT_DIR, outName);

  fs.writeFileSync(tempHtmlPath, htmlContent, 'utf8');

  const cmd = `"${CHROME_PATH}" --headless=new --no-sandbox --hide-scrollbars --window-size=${width},${height} --screenshot="${tempPngPath}" "file:///${tempHtmlPath.replace(/\\/g, '/')}"`;
  execSync(cmd, { stdio: 'inherit' });

  fs.copyFileSync(tempPngPath, outPngPath);
  console.log(`✓ Saved ${outName}`);
}

// 1. Terminal Lint Screenshot
const lintHtml = `
<div style="margin-bottom: 12px; color: #e2e8f0; font-weight: 600;">
  FILE: <span style="color: #38bdf8;">./workflows/invoice_dispatch.json</span>
</div>

<div style="margin-bottom: 8px;">
  <span class="badge badge-critical">CRITICAL</span> <span class="loc">line 27</span> <span class="rule">SEC-001</span> <span class="msg">Node 'HTTP Request' contains hardcoded OpenAI API Key (sk-p...1a0b)</span><br>
  <span class="hint">Hint: Move sensitive values into an n8n credential vault or environment variable ({{ $env.VAR }}).</span><br>
  <span class="suggested">Suggested: {{ $env.HTTP_REQUEST_HEADER_0_VALUE }} (auto-fixable with --fix)</span>
</div>

<div style="margin-bottom: 8px;">
  <span class="badge badge-critical">CRITICAL</span> <span class="loc">line 43</span> <span class="rule">SEC-001</span> <span class="msg">Node 'HTTP Request' contains hardcoded Payment Provider Secret Key (sec_...9f2a)</span><br>
  <span class="hint">Hint: Never commit production secret keys into workflow source control.</span><br>
  <span class="suggested">Suggested: {{ $env.PAYMENT_PROVIDER_SECRET }} (auto-fixable with --fix)</span>
</div>

<div style="margin-bottom: 8px;">
  <span class="badge badge-high">HIGH</span> <span class="loc">line 1</span> <span class="rule">REL-001</span> <span class="msg">Workflow does not configure an errorWorkflow handler in settings</span><br>
  <span class="hint">Hint: Set errorWorkflow in workflow settings to designate an automated dead-letter triage flow.</span>
</div>

<div style="margin-bottom: 8px;">
  <span class="badge badge-high">HIGH</span> <span class="loc">line 6</span> <span class="rule">SEC-002</span> <span class="msg">Webhook node 'Webhook' has authentication disabled ('none')</span><br>
  <span class="hint">Hint: Configure authentication (headerAuth, basicAuth, or jwtAuth) on inbound webhook endpoints.</span>
</div>

<div style="margin-bottom: 8px;">
  <span class="badge badge-high">HIGH</span> <span class="loc">line 55</span> <span class="rule">SEC-009</span> <span class="msg">Database node 'Postgres' uses raw expression interpolation in SQL statement</span><br>
  <span class="hint">Hint: Use parameterized queries ($1, $2) to eliminate SQL injection vectors.</span>
</div>

<div style="margin-bottom: 8px;">
  <span class="badge badge-medium">MEDIUM</span> <span class="loc">line 27</span> <span class="rule">REL-002</span> <span class="msg">HTTP Request node 'HTTP Request' does not enable retry-on-failure</span><br>
  <span class="hint">Hint: Enable retryOnFail with maxTries=3 and exponential backoff to handle transient 5xx faults.</span>
</div>

<div style="margin-bottom: 8px;">
  <span class="badge badge-medium">MEDIUM</span> <span class="loc">line 86</span> <span class="rule">BEST-006</span> <span class="msg">Workflow includes committed test pinData for 1 nodes</span><br>
  <span class="hint">Hint: Clear pinned test data or run 'n8nflow-guard fix --strip-pindata' before pushing.</span>
</div>

<div class="stat-card">
  <div class="stat-item"><div class="stat-num" style="color: #ef4444;">F (24/100)</div><div class="stat-lbl">Security Score</div></div>
  <div class="stat-item"><div class="stat-num" style="color: #f87171;">2</div><div class="stat-lbl">Critical</div></div>
  <div class="stat-item"><div class="stat-num" style="color: #fb923c;">4</div><div class="stat-lbl">High</div></div>
  <div class="stat-item"><div class="stat-num" style="color: #fbbf24;">7</div><div class="stat-lbl">Medium</div></div>
  <div class="stat-item"><div class="stat-num" style="color: #38bdf8;">4</div><div class="stat-lbl">Low</div></div>
  <div class="stat-item"><div class="stat-num" style="color: #4ade80;">3 Fixable</div><div class="stat-lbl">Auto Remediation</div></div>
</div>
`;

// 2. Terminal Fix Diff Screenshot
const fixHtml = `
<div style="margin-bottom: 12px; color: #94a3b8;">
  Scanning for auto-remediable security and hygiene issues across 3 workflow files...
</div>

<div style="font-weight: 600; color: #e2e8f0; margin-bottom: 10px;">
  --- Diff for ./workflows/invoice_dispatch.json ---
</div>

<div style="font-family: inherit; font-size: 13px; line-height: 1.5; background: #0b0f19; padding: 16px; border-radius: 8px; border: 1px solid #1e293b;">
  <div class="diff-del">- &quot;value&quot;: &quot;Bearer sk-proj-9a8b7c6d5e4f3a2b1c0d9e8f7a6b5c4d3e2f1a0b&quot;</div>
  <div class="diff-add">+ &quot;value&quot;: &quot;={{ $env.HTTP_REQUEST_HEADER_0_VALUE }}&quot;</div>
  <div style="color: #475569; margin: 4px 0;">@@ -42,7 +42,7 @@ Node: Stripe Charge @@</div>
  <div class="diff-del">- &quot;apiKey&quot;: &quot;mock_stripe_secret_key_9f2a0b1c&quot;</div>
  <div class="diff-add">+ &quot;apiKey&quot;: &quot;={{ $env.PAYMENT_PROVIDER_SECRET }}&quot;</div>
  <div style="color: #475569; margin: 4px 0;">@@ -86,12 +86,2 @@ pinData Sanitization @@</div>
  <div class="diff-del">- &quot;pinData&quot;: { &quot;HTTP Request&quot;: [{ &quot;customer_id&quot;: &quot;c_123&quot;, &quot;card&quot;: &quot;4111...4444&quot; }] }</div>
  <div class="diff-add">+ &quot;pinData&quot;: {}</div>
</div>

<div style="margin-top: 18px; padding: 12px 16px; background: rgba(34, 197, 94, 0.1); border: 1px solid #16a34a; border-radius: 6px; color: #4ade80;">
  ✓ Extracted 2 credentials into <strong>.env.example</strong> with auto-generated environment variables.<br>
  ✓ Stripped confidential test payloads from <strong>pinData</strong>.<br>
  ✓ Summary: 1 file modified, 2 secrets sanitized, 1 pinData cleared. Zero syntax errors.
</div>
`;

// 3. Terminal DAG Graph Visualization Screenshot
const graphHtml = `
<div style="color: #94a3b8; margin-bottom: 12px;">
  Topological Execution Path & Cycle Analysis for <strong>OrderFulfillment.json</strong>:
</div>

<div style="font-family: inherit; font-size: 13px; line-height: 1.5; background: #0b0f19; padding: 16px; border-radius: 8px; border: 1px solid #1e293b;">
  <div style="color: #38bdf8; font-weight: 600;">Workflow: Order Fulfillment Pipeline (5 nodes, 0 cycles, 0 orphans)</div>
  <div style="color: #475569;">--------------------------------------------------------------------------------</div>
  <div><span style="color: #34d399; font-weight: bold;">[Trigger]</span> <span style="color: #f1f5f9;">Shopify New Order Webhook</span> <span style="color: #64748b;">(n8n-nodes-base.webhook)</span></div>
  <div style="padding-left: 20px; color: #94a3b8;">-&gt; to: <span style="color: #38bdf8;">Validate Inventory Stock</span></div>
  <div style="margin-top: 8px;"><span style="color: #38bdf8; font-weight: bold;">[Action]</span>  <span style="color: #f1f5f9;">Validate Inventory Stock</span> <span style="color: #64748b;">(n8n-nodes-base.postgres)</span></div>
  <div style="padding-left: 20px; color: #94a3b8;">&lt;- from: <span style="color: #94a3b8;">Shopify New Order Webhook</span></div>
  <div style="padding-left: 20px; color: #94a3b8;">-&gt; to:   <span style="color: #38bdf8;">Dispatch Stripe Payment</span></div>
  <div style="margin-top: 8px;"><span style="color: #38bdf8; font-weight: bold;">[Action]</span>  <span style="color: #f1f5f9;">Dispatch Stripe Payment</span> <span style="color: #64748b;">(n8n-nodes-base.httpRequest)</span></div>
  <div style="padding-left: 20px; color: #94a3b8;">&lt;- from: <span style="color: #94a3b8;">Validate Inventory Stock</span></div>
  <div style="padding-left: 20px; color: #94a3b8;">-&gt; to:   <span style="color: #38bdf8;">Send Customer Confirmation Email</span></div>
  <div style="margin-top: 8px;"><span style="color: #38bdf8; font-weight: bold;">[Action]</span>  <span style="color: #f1f5f9;">Send Customer Confirmation Email</span> <span style="color: #64748b;">(n8n-nodes-base.emailSend)</span></div>
  <div style="padding-left: 20px; color: #94a3b8;">&lt;- from: <span style="color: #94a3b8;">Dispatch Stripe Payment</span></div>
  <div style="padding-left: 20px; color: #4ade80;">-&gt; [Terminal Output]</div>
</div>

<div style="margin-top: 16px; color: #4ade80; display: flex; align-items: center; gap: 8px;">
  <span>✓ Graph Topology Status: DAG Verified (Single Connected Component, 0 Deadlocks)</span>
</div>
`;

console.log('Generating screenshots...');
takeScreenshot(wrapInTerminal('n8nflow-guard lint --format terminal', 'n8nflow-guard lint ./workflows/*.json', lintHtml), 't_lint', 'n8nflow-guard_terminal.png');
takeScreenshot(wrapInTerminal('n8nflow-guard fix --dry-run', 'n8nflow-guard fix --dry-run ./workflows/*.json', fixHtml), 't_fix', 'n8nflow-guard_fix_diff.png');
takeScreenshot(wrapInTerminal('n8nflow-guard graph --format ascii', 'n8nflow-guard graph ./workflows/OrderFulfillment.json', graphHtml), 't_graph', 'n8nflow-guard_graph.png');

console.log('All individual screenshots captured successfully.');
