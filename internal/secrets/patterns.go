package secrets

import (
	"regexp"
	"strings"
)

// SecretMatch records an occurrence of a detected credential or secret token.
type SecretMatch struct {
	PatternName string
	MatchedText string
	Description string
	IsMasked    bool
}

// SecretPattern defines a signature for identifying sensitive tokens.
type SecretPattern struct {
	Name        string
	Pattern     *regexp.Regexp
	Description string
}

var knownPatterns = []SecretPattern{
	{
		Name:        "OpenAI API Key",
		Pattern:     regexp.MustCompile(`\b(sk-(?:proj-)?[A-Za-z0-9_-]{32,})\b`),
		Description: "OpenAI API secret key",
	},
	{
		Name:        "Anthropic API Key",
		Pattern:     regexp.MustCompile(`\b(sk-ant-[A-Za-z0-9_-]{32,})\b`),
		Description: "Anthropic Claude API key",
	},
	{
		Name:        "AWS Access Key ID",
		Pattern:     regexp.MustCompile(`\b(AKIA[0-9A-Z]{16})\b`),
		Description: "Amazon Web Services Access Key ID",
	},
	{
		Name:        "AWS Secret Access Key",
		Pattern:     regexp.MustCompile(`(?i)(?:aws_secret_access_key|aws_secret|secret_key)\s*[:=]\s*["']?([A-Za-z0-9/+=]{40})["']?`),
		Description: "Amazon Web Services Secret Access Key",
	},
	{
		Name:        "GitHub Personal Access Token",
		Pattern:     regexp.MustCompile(`\b((?:ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9_]{36})\b`),
		Description: "GitHub classic personal access token",
	},
	{
		Name:        "GitHub Fine-grained PAT",
		Pattern:     regexp.MustCompile(`\b(github_pat_[A-Za-z0-9_]{82})\b`),
		Description: "GitHub fine-grained personal access token",
	},
	{
		Name:        "GitLab Personal Access Token",
		Pattern:     regexp.MustCompile(`\b(glpat-[a-zA-Z0-9\-_]{20,})\b`),
		Description: "GitLab personal access token",
	},
	{
		Name:        "Slack Bot Token",
		Pattern:     regexp.MustCompile(`\b(xoxb-[0-9]{11,13}-[0-9]{11,13}-[a-zA-Z0-9]{24})\b`),
		Description: "Slack bot API token",
	},
	{
		Name:        "Slack User Token",
		Pattern:     regexp.MustCompile(`\b(xoxp-[0-9]{11,13}-[0-9]{11,13}-[a-zA-Z0-9]{24,32})\b`),
		Description: "Slack user token",
	},
	{
		Name:        "Slack Incoming Webhook",
		Pattern:     regexp.MustCompile(`https://hooks\.slack\.com/services/T[0-9A-Za-z]+/B[0-9A-Za-z]+/[0-9A-Za-z]+`),
		Description: "Slack incoming webhook URL",
	},
	{
		Name:        "Stripe Secret Key",
		Pattern:     regexp.MustCompile(`\b((?:sk|rk)_(?:live|test)_[0-9a-zA-Z]{24,})\b`),
		Description: "Stripe live or test secret API key",
	},
	{
		Name:        "Discord Bot Token",
		Pattern:     regexp.MustCompile(`\b([MN][A-Za-z\d]{23,}\.[\w-]{6}\.[\w-]{27,})\b`),
		Description: "Discord bot authorization token",
	},
	{
		Name:        "Discord Webhook",
		Pattern:     regexp.MustCompile(`https://(?:ptb\.|canary\.)?discord(?:app)?\.com/api/webhooks/\d+/[A-Za-z0-9_-]+`),
		Description: "Discord incoming webhook URL",
	},
	{
		Name:        "SendGrid API Key",
		Pattern:     regexp.MustCompile(`\b(SG\.[a-zA-Z0-9_-]{22}\.[a-zA-Z0-9_-]{43})\b`),
		Description: "SendGrid API credentials",
	},
	{
		Name:        "Twilio Account SID & Auth",
		Pattern:     regexp.MustCompile(`\b(AC[a-f0-9]{32})\b`),
		Description: "Twilio account identifier",
	},
	{
		Name:        "Google Cloud API Key",
		Pattern:     regexp.MustCompile(`\b(AIza[0-9A-Za-z-_]{35})\b`),
		Description: "Google Cloud Platform API key",
	},
	{
		Name:        "JSON Web Token (JWT)",
		Pattern:     regexp.MustCompile(`\b(eyJ[A-Za-z0-9-_=]+\.eyJ[A-Za-z0-9-_=]+\.[A-Za-z0-9-_.+/=]*)\b`),
		Description: "JSON Web Token with payload",
	},
	{
		Name:        "Private RSA Key",
		Pattern:     regexp.MustCompile(`-----BEGIN (?:RSA )?PRIVATE KEY-----`),
		Description: "PEM encoded private cryptographic key",
	},
	{
		Name:        "Database Connection URL with Password",
		Pattern:     regexp.MustCompile(`(?i)(?:postgres|postgresql|mysql|mongodb|redis):\/\/[^:\s]+:([^@\s]+)@`),
		Description: "Database URI containing embedded plaintext password",
	},
	{
		Name:        "Bearer Authorization Header",
		Pattern:     regexp.MustCompile(`(?i)Bearer\s+([A-Za-z0-9\-_\.]{25,})`),
		Description: "Hardcoded Bearer authentication token",
	},
	{
		Name:        "Mailgun API Key",
		Pattern:     regexp.MustCompile(`\b(key-[0-9a-zA-Z]{32})\b`),
		Description: "Mailgun private API key",
	},
	{
		Name:        "Telegram Bot Token",
		Pattern:     regexp.MustCompile(`\b(\d{9,10}:[A-Za-z0-9_-]{35})\b`),
		Description: "Telegram bot HTTP API token",
	},
	{
		Name:        "Shopify Access Token",
		Pattern:     regexp.MustCompile(`\b(shpat_[a-fA-F0-9]{32})\b`),
		Description: "Shopify Admin API access token",
	},
	{
		Name:        "Postman API Key",
		Pattern:     regexp.MustCompile(`\b(PMAK-[a-zA-Z0-9]{24}-[a-zA-Z0-9]{34})\b`),
		Description: "Postman personal access API key",
	},
	{
		Name:        "HashiCorp Vault Token",
		Pattern:     regexp.MustCompile(`\b(s\.[a-zA-Z0-9]{24})\b`),
		Description: "HashiCorp Vault access token",
	},
}

// ScanString tests a string value against all known secret signatures.
func ScanString(text string) []*SecretMatch {
	var matches []*SecretMatch

	// Avoid scanning obvious n8n expression templates
	if strings.Contains(text, "{{") && strings.Contains(text, "}}") {
		return matches
	}

	for _, p := range knownPatterns {
		submatches := p.Pattern.FindAllString(text, -1)
		for _, m := range submatches {
			matches = append(matches, &SecretMatch{
				PatternName: p.Name,
				MatchedText: m,
				Description: p.Description,
			})
		}
	}

	return matches
}

// IsSensitiveKeyName checks whether a JSON parameter key represents a credential or secret.
func IsSensitiveKeyName(keyName string) bool {
	lower := strings.ToLower(keyName)
	sensitiveKeys := []string{
		"apikey",
		"api_key",
		"api-key",
		"secret",
		"secretkey",
		"secret_key",
		"password",
		"passwd",
		"token",
		"access_token",
		"accesstoken",
		"auth_token",
		"authtoken",
		"private_key",
		"privatekey",
		"client_secret",
		"clientsecret",
	}

	for _, k := range sensitiveKeys {
		if lower == k || strings.HasSuffix(lower, "_"+k) || strings.HasSuffix(lower, k) {
			return true
		}
	}
	return false
}

// MaskSecret returns a partially obscured representation of a sensitive value for safe reporting.
func MaskSecret(secret string) string {
	if len(secret) <= 8 {
		return "******"
	}
	prefix := secret[:4]
	suffix := secret[len(secret)-4:]
	return prefix + "..." + suffix
}
