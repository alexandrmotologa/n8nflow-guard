package remote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
	"github.com/alexandrmotologa/n8nflow-guard/internal/parser"
)

// Client handles communication with live n8n instance REST APIs.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient initializes an API client with standard timeouts.
func NewClient(baseURL, apiKey string) *Client {
	cleanURL := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if !strings.HasSuffix(cleanURL, "/api/v1") {
		cleanURL = cleanURL + "/api/v1"
	}

	return &Client{
		BaseURL: cleanURL,
		APIKey:  strings.TrimSpace(apiKey),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type listResponse struct {
	Data []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Active bool   `json:"active"`
	} `json:"data"`
	NextCursor string `json:"nextCursor"`
}

// FetchWorkflows retrieves all workflows (or only active ones) from the remote n8n server.
func (c *Client) FetchWorkflows(activeOnly bool) ([]*models.Workflow, error) {
	url := fmt.Sprintf("%s/workflows", c.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating list request: %w", err)
	}

	req.Header.Set("X-N8N-API-KEY", c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to n8n server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("n8n API returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var list listResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed decoding workflow list: %w", err)
	}

	var workflows []*models.Workflow
	for _, item := range list.Data {
		if activeOnly && !item.Active {
			continue
		}

		wf, err := c.FetchWorkflowDetail(item.ID)
		if err != nil {
			return nil, fmt.Errorf("failed fetching workflow %s (%s): %w", item.Name, item.ID, err)
		}
		workflows = append(workflows, wf)
	}

	return workflows, nil
}

// FetchWorkflowDetail retrieves the complete node and connection graph for a single workflow.
func (c *Client) FetchWorkflowDetail(workflowID string) (*models.Workflow, error) {
	url := fmt.Sprintf("%s/workflows/%s", c.BaseURL, workflowID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-N8N-API-KEY", c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed fetching detail (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	parsed, err := parser.ParseWorkflows(bodyBytes, fmt.Sprintf("remote://workflow/%s", workflowID))
	if err != nil {
		return nil, err
	}

	if len(parsed) == 0 {
		return nil, fmt.Errorf("empty workflow returned for ID %s", workflowID)
	}

	return parsed[0], nil
}
