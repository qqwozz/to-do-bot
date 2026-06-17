package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) CreatePlan(plan PlanRequest) error {
	data, err := json.Marshal(plan)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/plans", "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("backend %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *Client) GetPlansByDate(date string) ([]Plan, error) {
	var plans []Plan
	err := c.getJSON("/plans?date="+date, &plans)
	return plans, err
}

func (c *Client) GetPlansByDateRange(start, end string) ([]Plan, error) {
	var plans []Plan
	err := c.getJSON(fmt.Sprintf("/plans/range?start=%s&end=%s", start, end), &plans)
	return plans, err
}

func (c *Client) getJSON(path string, target any) error {
	resp, err := c.httpClient.Get(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("backend %d: %s", resp.StatusCode, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(target)
}
