package rate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type RestClient struct {
	httpClient *http.Client
	config     ClientConfig
}

func NewRestClient(httpClient *http.Client, config ClientConfig) *RestClient {
	config = config.Default()
	return &RestClient{
		httpClient: httpClient,
		config:     config,
	}
}

func (c *RestClient) Latest(ctx context.Context, req LatestRequest) (*LatestResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.config.Endpoint+"/api/latest.json", nil)
	if err != nil {
		return nil, err
	}

	if req.ETag != "" {
		hreq.Header.Set("If-None-Match", strconv.Quote(req.ETag))
	}

	if !req.Modified.IsZero() {
		hreq.Header.Set("If-Modified-Since", req.Modified.UTC().Format(http.TimeFormat))
	}

	query := hreq.URL.Query()
	query.Set("app_id", c.config.AppID)
	query.Set("base", req.Base)
	if len(req.Symbols) > 0 {
		query.Set("symbols", strings.Join(req.Symbols, ","))
	}
	if req.ShowAlternative {
		query.Set("show_alternative", "1")
	}
	hreq.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(hreq)
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

	res := &LatestResponse{
		ETag: resp.Header.Get("ETag"),
	}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, fmt.Errorf("rate: unmarshal response: %w", err)
	}

	return res, nil
}
