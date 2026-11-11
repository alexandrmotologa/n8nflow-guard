const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const CHROME_PATH = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe";
const FFMPEG_PATH = "C:\\Users\\alexander\\AppData\\Local\\Microsoft\\WinGet\\Links\\ffmpeg.exe";
const TEMP_DIR = "C:\\Users\\alexander\\AppData\\Local\\Temp";
const OUTPUT_DIR = "B:\\workgit\\n8nflow-guard\\docs\\images";

function wrapInFrame(title, badgeText, contentHtml, width = 1200, height = 750) {
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
    width: ${width}px;
    height: ${height}px;
    padding: 24px;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  }
  .window {
    width: 100%;
    height: 100%;
    background: #0f172a;
    border-radius: 12px;
    border: 1px solid #1e293b;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.7);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  .header {
    background: #1e293b;
    padding: 12px 18px;
    display: flex;
    align-items: center;
    border-bottom: 1px solid #334155;
    flex-shrink: 0;
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
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    padding-left: 30px;
  }
  .step-badge {
    background: #0284c7;
    color: #ffffff;
    font-size: 11px;
    font-weight: 700;
    padding: 2px 8px;
    border-radius: 4px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  .body {
    padding: 22px 26px;
    color: #f1f5f9;
    font-family: ui-monospace, SFMono-Regular, "Cascadia Code", "Source Code Pro", Menlo, Monaco, Consolas, monospace;
    font-size: 13.5px;
    line-height: 1.55;
    flex-grow: 1;
    overflow: hidden;
  }
  .prompt-line {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid #1e293b;
  }
  .prompt-symbol { color: #00f5ff; font-weight: bold; }
  .prompt-cmd { color: #f8fafc; font-weight: 600; font-size: 14px; }
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
  .loc { color: #64748b; }
  .rule { color: #38bdf8; font-weight: 600; }
  .msg { color: #f1f5f9; }
  .hint { color: #94a3b8; padding-left: 20px; font-size: 12.5px; }
  .suggested { color: #34d399; padding-left: 20px; font-size: 12.5px; }
  .diff-add { color: #4ade80; background: rgba(34, 197, 94, 0.1); }
  .diff-del { color: #f87171; background: rgba(239, 68, 68, 0.1); }
  .stat-card {
    background: #1e293b;
    border: 1px solid #334155;
    border-radius: 8px;
    padding: 12px 18px;
    margin-top: 14px;
    display: flex;
    justify-content: space-between;
  }
  .stat-item { text-align: center; }
  .stat-num { font-size: 19px; font-weight: 700; color: #00f5ff; }
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
      <div class="step-badge">${badgeText}</div>
    </div>
    <div class="body">
      ${contentHtml}
    </div>
  </div>
</body>
</html>`;
}

// Frame 1: Scanning & Detection
const f1Content = `
  <div class="prompt-line">
    <span class="prompt-symbol">❯</span>
    <span class="prompt-cmd">n8nflow-guard lint ./workflows/*.json</span>
  </div>
  <div style="margin-bottom: 10px; color: #e2e8f0; font-weight: 600;">
    FILE: <span style="color: #38bdf8;">./workflows/invoice_dispatch.json</span>
  </div>
  <div style="margin-bottom: 6px;">
    <span class="badge badge-critical">CRITICAL</span> <span class="loc">line 27</span> <span class="rule">SEC-001</span> <span class="msg">Node 'HTTP Request' contains hardcoded OpenAI API Key (sk-p...1a0b)</span><br>
    <span class="suggested">Suggested: {{ $env.HTTP_REQUEST_HEADER_0_VALUE }} (auto-fixable with --fix)</span>
  </div>
  <div style="margin-bottom: 6px;">
    <span class="badge badge-critical">CRITICAL</span> <span class="loc">line 43</span> <span class="rule">SEC-001</span> <span class="msg">Node 'HTTP Request' contains hardcoded Payment Provider Secret Key (sec_...9f2a)</span><br>
    <span class="suggested">Suggested: {{ $env.PAYMENT_PROVIDER_SECRET }} (auto-fixable with --fix)</span>
  </div>
  <div style="margin-bottom: 6px;">
    <span class="badge badge-high">HIGH</span> <span class="loc">line 1</span> <span class="rule">REL-001</span> <span class="msg">Workflow does not configure an errorWorkflow handler in settings</span><br>
    <span class="hint">Hint: Set errorWorkflow in settings to route uncaught failures to a DLQ triage workflow.</span>
  </div>
  <div style="margin-bottom: 6px;">
    <span class="badge badge-high">HIGH</span> <span class="loc">line 6</span> <span class="rule">SEC-002</span> <span class="msg">Webhook node 'Webhook' has authentication disabled ('none')</span>
  </div>
  <div style="margin-bottom: 6px;">
    <span class="badge badge-high">HIGH</span> <span class="loc">line 55</span> <span class="rule">SEC-009</span> <span class="msg">Database node 'Postgres' uses raw string interpolation in SQL query</span>
  </div>
  <div style="margin-bottom: 6px;">
    <span class="badge badge-medium">MEDIUM</span> <span class="loc">line 86</span> <span class="rule">BEST-006</span> <span class="msg">Workflow includes committed test pinData for 1 nodes</span>
  </div>
  <div class="stat-card">
    <div class="stat-item"><div class="stat-num" style="color: #ef4444;">F (24/100)</div><div class="stat-lbl">Security Score</div></div>
    <div class="stat-item"><div class="stat-num" style="color: #f87171;">2</div><div class="stat-lbl">Critical</div></div>
    <div class="stat-item"><div class="stat-num" style="color: #fb923c;">4</div><div class="stat-lbl">High</div></div>
    <div class="stat-item"><div class="stat-num" style="color: #fbbf24;">7</div><div class="stat-lbl">Medium</div></div>
    <div class="stat-item"><div class="stat-num" style="color: #4ade80;">3 Fixable</div><div class="stat-lbl">Auto Remediation</div></div>
  </div>
`;

// Frame 2: Auto-Remediation & Diff
const f2Content = `
  <div class="prompt-line">
    <span class="prompt-symbol">❯</span>
    <span class="prompt-cmd">n8nflow-guard fix --dry-run ./workflows/*.json</span>
  </div>
  <div style="color: #94a3b8; margin-bottom: 10px;">
    Scanning for auto-remediable security and hygiene issues across 3 workflow files...
  </div>
  <div style="font-weight: 600; color: #e2e8f0; margin-bottom: 8px;">
    --- Diff for ./workflows/invoice_dispatch.json ---
  </div>
  <div style="background: #0b0f19; padding: 14px; border-radius: 8px; border: 1px solid #1e293b; font-size: 13px;">
    <div class="diff-del">- &quot;value&quot;: &quot;Bearer sk-proj-9a8b7c6d5e4f3a2b1c0d9e8f7a6b5c4d3e2f1a0b&quot;</div>
    <div class="diff-add">+ &quot;value&quot;: &quot;={{ $env.HTTP_REQUEST_HEADER_0_VALUE }}&quot;</div>
    <div style="color: #475569; margin: 4px 0;">@@ -42,7 +42,7 @@ Node: Stripe Charge @@</div>
    <div class="diff-del">- &quot;apiKey&quot;: &quot;mock_stripe_secret_key_9f2a0b1c&quot;</div>
    <div class="diff-add">+ &quot;apiKey&quot;: &quot;={{ $env.PAYMENT_PROVIDER_SECRET }}&quot;</div>
    <div style="color: #475569; margin: 4px 0;">@@ -86,12 +86,2 @@ pinData Sanitization @@</div>
    <div class="diff-del">- &quot;pinData&quot;: { &quot;HTTP Request&quot;: [{ &quot;customer_id&quot;: &quot;c_123&quot;, &quot;card&quot;: &quot;4111...4444&quot; }] }</div>
    <div class="diff-add">+ &quot;pinData&quot;: {}</div>
  </div>
  <div style="margin-top: 14px; padding: 12px 16px; background: rgba(34, 197, 94, 0.1); border: 1px solid #16a34a; border-radius: 6px; color: #4ade80;">
    ✓ Extracted 2 credentials into <strong>.env.example</strong> with auto-generated environment variables.<br>
    ✓ Stripped confidential test payloads from <strong>pinData</strong>.<br>
    ✓ Summary: 1 file modified, 2 secrets sanitized, 1 pinData cleared. Zero syntax errors.
  </div>
`;

// Frame 3: Interactive Security Report Dashboard
const f3Content = `
  <div class="prompt-line">
    <span class="prompt-symbol">❯</span>
    <span class="prompt-cmd">n8nflow-guard lint --format html -o security_report.html</span>
  </div>
  <div style="background: #0b0f19; border: 1px solid #1e293b; border-radius: 8px; padding: 16px;">
    <div style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #1e293b; padding-bottom: 12px; margin-bottom: 12px;">
      <div>
        <div style="font-size: 16px; font-weight: 700; color: #f8fafc;">n8nflow-guard Security Report</div>
        <div style="font-size: 12px; color: #64748b;">Audited 3 workflow(s) &bull; OWASP API Security (2023) Aligned</div>
      </div>
      <div style="text-align: right;">
        <span style="font-size: 22px; font-weight: 800; color: #ef4444;">F</span>
        <span style="font-size: 14px; color: #94a3b8;"> (24/100)</span>
      </div>
    </div>
    <div style="display: grid; grid-template-columns: repeat(5, 1fr); gap: 8px; margin-bottom: 14px;">
      <div style="background: #1e293b; padding: 8px; border-radius: 6px; text-align: center;"><div style="font-size: 16px; font-weight: 700; color: #f1f5f9;">32</div><div style="font-size: 10px; color: #94a3b8;">Total Issues</div></div>
      <div style="background: #1e293b; padding: 8px; border-radius: 6px; text-align: center;"><div style="font-size: 16px; font-weight: 700; color: #f87171;">2</div><div style="font-size: 10px; color: #94a3b8;">Critical</div></div>
      <div style="background: #1e293b; padding: 8px; border-radius: 6px; text-align: center;"><div style="font-size: 16px; font-weight: 700; color: #fb923c;">7</div><div style="font-size: 10px; color: #94a3b8;">High</div></div>
      <div style="background: #1e293b; padding: 8px; border-radius: 6px; text-align: center;"><div style="font-size: 16px; font-weight: 700; color: #fbbf24;">14</div><div style="font-size: 10px; color: #94a3b8;">Medium</div></div>
      <div style="background: #1e293b; padding: 8px; border-radius: 6px; text-align: center;"><div style="font-size: 16px; font-weight: 700; color: #4ade80;">3</div><div style="font-size: 10px; color: #94a3b8;">Auto-Fixable</div></div>
    </div>
    <div style="display: flex; gap: 8px; font-size: 12px; margin-bottom: 10px;">
      <span style="background: #0284c7; color: white; padding: 3px 10px; border-radius: 4px; font-weight: 600;">All (32)</span>
      <span style="background: #1e293b; color: #94a3b8; padding: 3px 10px; border-radius: 4px;">Critical (2)</span>
      <span style="background: #1e293b; color: #94a3b8; padding: 3px 10px; border-radius: 4px;">High (7)</span>
      <span style="background: #1e293b; color: #94a3b8; padding: 3px 10px; border-radius: 4px;">Medium (14)</span>
    </div>
    <div style="background: #1e293b; padding: 10px 14px; border-radius: 6px; border-left: 3px solid #dc2626; font-size: 12px;">
      <span style="color: #f87171; font-weight: 700;">SEC-001: hardcoded-credentials</span> &bull; <span style="color: #94a3b8;">./workflows/invoice_dispatch.json:27</span><br>
      <span style="color: #e2e8f0;">Node 'HTTP Request' contains hardcoded OpenAI API Key (sk-p...1a0b)</span>
    </div>
  </div>
`;

// Frame 4: DAG Topology & Cycle Verification
const f4Content = `
  <div class="prompt-line">
    <span class="prompt-symbol">❯</span>
    <span class="prompt-cmd">n8nflow-guard graph --format ascii ./workflows/OrderFulfillment.json</span>
  </div>
  <div style="font-family: inherit; font-size: 13px; line-height: 1.5; background: #0b0f19; padding: 14px; border-radius: 8px; border: 1px solid #1e293b;">
    <div style="color: #38bdf8; font-weight: 600;">Workflow: Order Fulfillment Pipeline (5 nodes, 0 cycles, 0 orphans)</div>
    <div style="color: #475569;">--------------------------------------------------------------------------------</div>
    <div><span style="color: #34d399; font-weight: bold;">[Trigger]</span> <span style="color: #f1f5f9;">Shopify New Order Webhook</span> <span style="color: #64748b;">(n8n-nodes-base.webhook)</span></div>
    <div style="padding-left: 20px; color: #94a3b8;">-&gt; to: <span style="color: #38bdf8;">Validate Inventory Stock</span></div>
    <div style="margin-top: 6px;"><span style="color: #38bdf8; font-weight: bold;">[Action]</span>  <span style="color: #f1f5f9;">Validate Inventory Stock</span> <span style="color: #64748b;">(n8n-nodes-base.postgres)</span></div>
    <div style="padding-left: 20px; color: #94a3b8;">&lt;- from: <span style="color: #94a3b8;">Shopify New Order Webhook</span> &nbsp; -&gt; to: <span style="color: #38bdf8;">Dispatch Stripe Payment</span></div>
    <div style="margin-top: 6px;"><span style="color: #38bdf8; font-weight: bold;">[Action]</span>  <span style="color: #f1f5f9;">Dispatch Stripe Payment</span> <span style="color: #64748b;">(n8n-nodes-base.httpRequest)</span></div>
    <div style="padding-left: 20px; color: #94a3b8;">&lt;- from: <span style="color: #94a3b8;">Validate Inventory Stock</span> &nbsp; -&gt; to: <span style="color: #38bdf8;">Send Customer Confirmation</span></div>
    <div style="margin-top: 6px;"><span style="color: #38bdf8; font-weight: bold;">[Action]</span>  <span style="color: #f1f5f9;">Send Customer Confirmation</span> <span style="color: #64748b;">(n8n-nodes-base.emailSend)</span></div>
    <div style="padding-left: 20px; color: #4ade80;">-&gt; [Terminal Output] (Execution Completed Cleanly)</div>
  </div>
  <div style="margin-top: 14px; padding: 10px 14px; background: rgba(34, 197, 94, 0.1); border: 1px solid #16a34a; border-radius: 6px; color: #4ade80; font-size: 13px;">
    ✓ Graph Topology: Valid Directed Acyclic Graph (DAG)<br>
    ✓ 0 Circular Dependencies &bull; 0 Orphan Disconnected Nodes &bull; Topologically Sorted
  </div>
`;

console.log('Rendering GIF frames...');
const frames = [
  { name: 'frame_01', html: wrapInFrame('n8nflow-guard lint', 'Step 1: Security Audit', f1Content) },
  { name: 'frame_02', html: wrapInFrame('n8nflow-guard fix', 'Step 2: Auto Remediation', f2Content) },
  { name: 'frame_03', html: wrapInFrame('n8nflow-guard dashboard', 'Step 3: HTML Dashboard', f3Content) },
  { name: 'frame_04', html: wrapInFrame('n8nflow-guard graph', 'Step 4: DAG Topology', f4Content) }
];

frames.forEach((f, idx) => {
  const tempHtml = path.join(TEMP_DIR, `${f.name}.html`);
  const tempPng = path.join(TEMP_DIR, `${f.name}.png`);
  fs.writeFileSync(tempHtml, f.html, 'utf8');
  const cmd = `"${CHROME_PATH}" --headless=new --no-sandbox --hide-scrollbars --window-size=1200,750 --screenshot="${tempPng}" "file:///${tempHtml.replace(/\\/g, '/')}"`;
  execSync(cmd, { stdio: 'inherit' });
  console.log(`✓ Rendered ${f.name}.png`);
});

const outGifPath = path.join(OUTPUT_DIR, 'n8nflow-guard_demo.gif');
console.log('Compiling animated GIF via ffmpeg...');

// Each frame shown for 2.5 seconds (fps=0.4 input or framerate 1/2.5)
const ffmpegCmd = `"${FFMPEG_PATH}" -y -framerate 0.4 -i "${TEMP_DIR}\\frame_%02d.png" -vf "fps=10,scale=1000:-1:flags=lanczos,split[s0][s1];[s0]palettegen=max_colors=128:reserve_transparent=0[p];[s1][p]paletteuse=dither=bayer:bayer_scale=3" "${outGifPath}"`;
execSync(ffmpegCmd, { stdio: 'inherit' });

console.log(`✓ Successfully compiled animated demo GIF: ${outGifPath}`);
