package reporter

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// HTMLReporter generates a self-contained interactive HTML audit report.
type HTMLReporter struct{}

func (r *HTMLReporter) Report(w io.Writer, issues []models.Issue, summary Summary) error {
	deductions := (summary.CriticalIssues * 15) + (summary.HighIssues * 8) + (summary.MediumIssues * 3) + summary.LowIssues
	score := 100 - deductions
	if score < 0 {
		score = 0
	}

	grade := "A+"
	gradeColor := "#22c55e"
	switch {
	case score >= 95:
		grade = "A+"
		gradeColor = "#22c55e"
	case score >= 90:
		grade = "A"
		gradeColor = "#16a34a"
	case score >= 80:
		grade = "B"
		gradeColor = "#3b82f6"
	case score >= 70:
		grade = "C"
		gradeColor = "#eab308"
	case score >= 60:
		grade = "D"
		gradeColor = "#f97316"
	default:
		grade = "F"
		gradeColor = "#ef4444"
	}

	issuesJSON, _ := json.Marshal(issues)

	htmlTmpl := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>n8nflow-guard Security Audit Report</title>
  <style>
    :root {
      --bg: #0f172a;
      --card-bg: #1e293b;
      --border: #334155;
      --text: #f8fafc;
      --text-muted: #94a3b8;
      --critical: #ef4444;
      --high: #f97316;
      --medium: #eab308;
      --low: #06b6d4;
      --success: #22c55e;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: var(--bg); color: var(--text); padding: 2rem; line-height: 1.5; }
    .container { max-width: 1200px; margin: 0 auto; }
    header { display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border); padding-bottom: 1.5rem; margin-bottom: 2rem; }
    h1 { font-size: 1.75rem; font-weight: 700; }
    .subtitle { color: var(--text-muted); font-size: 0.9rem; margin-top: 0.25rem; }
    .score-badge { display: flex; align-items: center; gap: 1rem; background: var(--card-bg); border: 1px solid var(--border); padding: 0.75rem 1.5rem; border-radius: 0.75rem; }
    .grade { font-size: 2.25rem; font-weight: 800; color: %s; }
    .score-meta { text-align: right; }
    .score-val { font-size: 1.25rem; font-weight: 700; }
    .metrics-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
    .metric-card { background: var(--card-bg); border: 1px solid var(--border); padding: 1.25rem; border-radius: 0.5rem; text-align: center; }
    .metric-num { font-size: 1.75rem; font-weight: 700; margin-top: 0.25rem; }
    .critical { color: var(--critical); }
    .high { color: var(--high); }
    .medium { color: var(--medium); }
    .low { color: var(--low); }
    .controls { display: flex; gap: 1rem; margin-bottom: 1.5rem; flex-wrap: wrap; }
    .search-box { flex: 1; min-width: 250px; padding: 0.75rem 1rem; background: var(--card-bg); border: 1px solid var(--border); color: var(--text); border-radius: 0.5rem; font-size: 0.95rem; }
    .filter-btn { padding: 0.5rem 1rem; background: var(--card-bg); border: 1px solid var(--border); color: var(--text); border-radius: 0.5rem; cursor: pointer; font-size: 0.85rem; font-weight: 600; }
    .filter-btn.active { background: #3b82f6; border-color: #3b82f6; }
    .issues-list { display: flex; flex-direction: column; gap: 1rem; }
    .issue-card { background: var(--card-bg); border: 1px solid var(--border); border-radius: 0.5rem; padding: 1.25rem; }
    .issue-header { display: flex; justify-content: space-between; align-items: flex-start; gap: 1rem; margin-bottom: 0.5rem; }
    .badge { font-size: 0.75rem; font-weight: 700; padding: 0.2rem 0.5rem; border-radius: 0.25rem; text-transform: uppercase; }
    .badge-critical { background: rgba(239, 68, 68, 0.2); color: var(--critical); border: 1px solid var(--critical); }
    .badge-high { background: rgba(249, 115, 22, 0.2); color: var(--high); border: 1px solid var(--high); }
    .badge-medium { background: rgba(234, 179, 8, 0.2); color: var(--medium); border: 1px solid var(--medium); }
    .badge-low { background: rgba(6, 182, 212, 0.2); color: var(--low); border: 1px solid var(--low); }
    .issue-loc { color: var(--text-muted); font-size: 0.85rem; font-family: monospace; }
    .issue-msg { font-size: 1rem; font-weight: 600; margin-bottom: 0.5rem; }
    .issue-hint { font-size: 0.85rem; color: var(--text-muted); background: rgba(0,0,0,0.2); padding: 0.5rem 0.75rem; border-radius: 0.25rem; }
    .empty-state { text-align: center; padding: 3rem; color: var(--text-muted); font-size: 1.1rem; }
  </style>
</head>
<body>
  <div class="container">
    <header>
      <div>
        <h1>n8nflow-guard Security Report</h1>
        <div class="subtitle">Audited %d workflow(s) in %d file(s)</div>
      </div>
      <div class="score-badge">
        <div class="grade">%s</div>
        <div class="score-meta">
          <div class="score-val">%d / 100</div>
          <div class="subtitle">Security Score</div>
        </div>
      </div>
    </header>

    <div class="metrics-grid">
      <div class="metric-card">
        <div class="subtitle">Total Issues</div>
        <div class="metric-num">%d</div>
      </div>
      <div class="metric-card">
        <div class="subtitle">Critical</div>
        <div class="metric-num critical">%d</div>
      </div>
      <div class="metric-card">
        <div class="subtitle">High</div>
        <div class="metric-num high">%d</div>
      </div>
      <div class="metric-card">
        <div class="subtitle">Medium</div>
        <div class="metric-num medium">%d</div>
      </div>
      <div class="metric-card">
        <div class="subtitle">Low</div>
        <div class="metric-num low">%d</div>
      </div>
      <div class="metric-card">
        <div class="subtitle">Auto-Fixable</div>
        <div class="metric-num" style="color: var(--success)">%d</div>
      </div>
    </div>

    <div class="controls">
      <input type="text" id="search" class="search-box" placeholder="Filter issues by node, rule, file, or keyword..." oninput="renderIssues()">
      <button class="filter-btn active" onclick="setSeverity('ALL')">All</button>
      <button class="filter-btn" onclick="setSeverity('CRITICAL')">Critical</button>
      <button class="filter-btn" onclick="setSeverity('HIGH')">High</button>
      <button class="filter-btn" onclick="setSeverity('MEDIUM')">Medium</button>
      <button class="filter-btn" onclick="setSeverity('LOW')">Low</button>
    </div>

    <div id="issues-container" class="issues-list"></div>
  </div>

  <script>
    const issues = %s;
    let selectedSeverity = 'ALL';

    function setSeverity(sev) {
      selectedSeverity = sev;
      document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.classList.toggle('active', btn.textContent.toUpperCase() === sev);
      });
      renderIssues();
    }

    function renderIssues() {
      const q = document.getElementById('search').value.toLowerCase();
      const container = document.getElementById('issues-container');
      container.innerHTML = '';

      const filtered = issues.filter(iss => {
        if (selectedSeverity !== 'ALL' && iss.severity !== selectedSeverity) return false;
        if (!q) return true;
        return (iss.message || '').toLowerCase().includes(q) ||
               (iss.ruleId || '').toLowerCase().includes(q) ||
               (iss.nodeName || '').toLowerCase().includes(q) ||
               (iss.filePath || '').toLowerCase().includes(q);
      });

      if (filtered.length === 0) {
        container.innerHTML = '<div class="empty-state">No issues matching criteria.</div>';
        return;
      }

      filtered.forEach(iss => {
        const card = document.createElement('div');
        card.className = 'issue-card';
        const sevClass = (iss.severity || 'low').toLowerCase();

        card.innerHTML = ` + "`" + `
          <div class="issue-header">
            <div>
              <span class="badge badge-${sevClass}">${iss.severity}</span>
              <strong style="margin-left: 0.5rem;">${iss.ruleId}</strong>: ${iss.ruleName}
            </div>
            <div class="issue-loc">${iss.filePath}:${iss.line}</div>
          </div>
          <div class="issue-msg">${iss.message}</div>
          ${iss.remediation ? '<div class="issue-hint"><strong>Hint:</strong> ' + iss.remediation + '</div>' : ''}
        ` + "`" + `;
        container.appendChild(card);
      });
    }

    renderIssues();
  </script>
</body>
</html>`

	_, err := fmt.Fprintf(w, htmlTmpl,
		gradeColor,
		summary.TotalWorkflows,
		summary.TotalFilesScanned,
		grade,
		score,
		summary.TotalIssues,
		summary.CriticalIssues,
		summary.HighIssues,
		summary.MediumIssues,
		summary.LowIssues,
		summary.FixableIssues,
		string(issuesJSON),
	)
	return err
}

func escapeHTML(s string) string {
	return html.EscapeString(strings.TrimSpace(s))
}
