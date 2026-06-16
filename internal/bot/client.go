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
		return fmt.Errorf("ошибка маршалинга: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/plans",
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		return fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("backend вернул %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) GetPlansByDate(date string) ([]Plan, error) {
	var plans []Plan
	err := c.getJSON("/plans?date="+date, &plans)
	return plans, err
}

func (c *Client) GetPlansByDateRange(startDate, endDate string) ([]Plan, error) {
	var plans []Plan
	url := fmt.Sprintf("/plans/range?start=%s&end=%s", startDate, endDate)
	err := c.getJSON(url, &plans)
	return plans, err
}

func (c *Client) getJSON(path string, target interface{}) error {
	resp, err := c.httpClient.Get(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("backend вернул %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
