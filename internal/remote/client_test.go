package remote

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemoteClient(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-N8N-API-KEY")
		if apiKey != "test-api-key" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.URL.Path == "/api/v1/workflows" {
			resp := map[string]interface{}{
				"data": []map[string]interface{}{
					{"id": "wf-1", "name": "Active Flow", "active": true},
					{"id": "wf-2", "name": "Draft Flow", "active": false},
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		if r.URL.Path == "/api/v1/workflows/wf-1" {
			detail := map[string]interface{}{
				"id":   "wf-1",
				"name": "Active Flow",
				"nodes": []map[string]interface{}{
					{"name": "Trigger", "type": "n8n-nodes-base.scheduleTrigger"},
				},
				"connections": map[string]interface{}{},
			}
			json.NewEncoder(w).Encode(detail)
			return
		}

		http.NotFound(w, r)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL, "test-api-key")
	workflows, err := client.FetchWorkflows(true) // activeOnly = true
	if err != nil {
		t.Fatalf("failed fetching workflows: %v", err)
	}

	if len(workflows) != 1 {
		t.Fatalf("expected 1 active workflow, got %d", len(workflows))
	}

	if workflows[0].Name != "Active Flow" {
		t.Errorf("expected workflow name 'Active Flow', got %s", workflows[0].Name)
	}
}
