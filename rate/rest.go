package rate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RestClient struct {
	httpClient *http.Client
	config     Config
	symbols    string
}

func NewRestClient(httpClient *http.Client, config Config) *RestClient {
	config = config.Normalize()
	symbols := strings.Join(config.Symbols, ",")
	return &RestClient{
		httpClient: httpClient,
		config:     config,
		symbols:    symbols,
	}
}

func (c *RestClient) Base() string {
	return c.config.Base
}

func (c *RestClient) Latest(ctx context.Context) (*Latest, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.config.Endpoint+"/api/latest.json", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("app_id", c.config.AppID)
	if c.symbols != "" {
		q.Set("symbols", c.symbols)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rate: request latest: %w", err)
	}
	defer resp.Body.Close() // nolint: errcheck

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("rate: read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ClientError{
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}

	res := &Latest{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, fmt.Errorf("rate: unmarshal response: %w", err)
	}

	return res, nil
}
