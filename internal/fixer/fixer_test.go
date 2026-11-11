package fixer

import (
	"strings"
	"testing"
)

func TestFixWorkflowContent(t *testing.T) {
	rawJSON := []byte(`{
		"name": "Payment Flow",
		"nodes": [
			{
				"name": "Stripe Charger",
				"type": "n8n-nodes-base.httpRequest",
				"parameters": {
					"apiKey": "sk_test_mockkey1234567890abcdef123456"
				}
			}
		],
		"pinData": {
			"Stripe Charger": [{"amount": 5000}]
		}
	}`)

	opts := FixOptions{
		InPlace:      false,
		DryRun:       true,
		StripPinData: true,
	}

	res, err := FixWorkflowContent(rawJSON, "payment.json", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.Modified {
		t.Errorf("expected workflow to be marked modified")
	}

	if res.SecretsReplaced != 1 {
		t.Errorf("expected 1 secret replaced, got %d", res.SecretsReplaced)
	}

	if !res.PinDataStripped {
		t.Errorf("expected pinData to be stripped")
	}

	if !strings.Contains(res.FixedJSON, "$env.STRIPE_CHARGER_APIKEY") {
		t.Errorf("expected fixed JSON to contain $env reference, got:\n%s", res.FixedJSON)
	}

	if strings.Contains(res.FixedJSON, "sk_test_") {
		t.Errorf("expected fixed JSON to NOT contain raw secret, got:\n%s", res.FixedJSON)
	}

	if !strings.Contains(res.FixedJSON, `"pinData": {}`) {
		t.Errorf("expected pinData to be cleared to empty object, got:\n%s", res.FixedJSON)
	}
}
