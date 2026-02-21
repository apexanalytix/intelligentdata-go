// Package intelligentdata provides a Go client for the Intelligent Data API by apexanalytix.
//
// Usage:
//
//	client := intelligentdata.NewClient("svm...")
//	resp, err := client.ValidateAddress(ctx, intelligentdata.AddressRequest{
//	    AddressLine1: "123 Main St", City: "New York", State: "NY",
//	    PostalCode: "10001", Country: "US",
//	})
package intelligentdata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

const (
	version        = "0.1.0"
	defaultBaseURL = "https://api.smartvmapi.com"
	maxRetries     = 3
)

// Client is the Intelligent Data API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	oauth      *oauth2TokenManager
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// WithHTTPClient provides a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithOAuth2 enables OAuth2 client credentials authentication.
func WithOAuth2(clientID, clientSecret, tokenURL string) Option {
	return func(c *Client) {
		if tokenURL == "" {
			tokenURL = c.baseURL + "/api/oauth/token"
		}
		c.oauth = newOAuth2TokenManager(clientID, clientSecret, tokenURL, c.httpClient)
	}
}

// NewClient creates a new Intelligent Data API client.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (map[string]interface{}, error) {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
	}

	url := c.baseURL + path
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "intelligentdata-go-sdk/"+version)

		if c.apiKey != "" {
			req.Header.Set("X-Api-Key", c.apiKey)
		} else if c.oauth != nil {
			token, err := c.oauth.getToken()
			if err != nil {
				return nil, err
			}
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < maxRetries-1 {
				sleep(ctx, backoff(attempt))
				continue
			}
			break
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == 429 {
			if attempt < maxRetries-1 {
				sleep(ctx, backoff(attempt))
				continue
			}
			return nil, &ApiError{StatusCode: 429, Message: "Rate limit exceeded"}
		}

		if resp.StatusCode >= 500 {
			lastErr = &ApiError{StatusCode: resp.StatusCode, Message: string(respBody)}
			if attempt < maxRetries-1 {
				sleep(ctx, backoff(attempt))
				continue
			}
			return nil, lastErr
		}

		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			msg := "Authentication failed"
			var raw map[string]interface{}
			if json.Unmarshal(respBody, &raw) == nil {
				if m, ok := raw["message"].(string); ok {
					msg = m
				}
			}
			return nil, &ApiError{StatusCode: resp.StatusCode, Message: msg, Raw: raw}
		}

		if resp.StatusCode >= 400 {
			msg := "Request failed"
			var raw map[string]interface{}
			if json.Unmarshal(respBody, &raw) == nil {
				if m, ok := raw["message"].(string); ok {
					msg = m
				}
			}
			return nil, &ApiError{StatusCode: resp.StatusCode, Message: msg, Raw: raw}
		}

		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		return result, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, &ApiError{StatusCode: 0, Message: "request failed after retries"}
}

func populateRaw(data map[string]interface{}, resp interface{}) {
	switch v := resp.(type) {
	case *AddressResponse:
		v.Raw = data
	case *TaxIdResponse:
		v.Raw = data
	case *BankAccountResponse:
		v.Raw = data
	case *BusinessLookupResponse:
		v.Raw = data
	case *SanctionsResponse:
		v.Raw = data
	case *DirectorsResponse:
		v.Raw = data
	}
}

// ── Public Methods ────────────────────────────────────────────────────────

// ValidateAddress validates and standardizes a postal address.
func (c *Client) ValidateAddress(ctx context.Context, req AddressRequest) (*AddressResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/validate/address", req)
	if err != nil {
		return nil, err
	}
	var resp AddressResponse
	remarshal(data, &resp)
	populateRaw(data, &resp)
	return &resp, nil
}

// ValidateTaxID validates a tax identification number.
func (c *Client) ValidateTaxID(ctx context.Context, req TaxIdRequest) (*TaxIdResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/validate/taxid", req)
	if err != nil {
		return nil, err
	}
	var resp TaxIdResponse
	remarshal(data, &resp)
	populateRaw(data, &resp)
	return &resp, nil
}

// ValidateBankAccount verifies bank account details.
func (c *Client) ValidateBankAccount(ctx context.Context, req BankAccountRequest) (*BankAccountResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/validate/bank", req)
	if err != nil {
		return nil, err
	}
	var resp BankAccountResponse
	remarshal(data, &resp)
	populateRaw(data, &resp)
	return &resp, nil
}

// LookupBusiness looks up official business registration data.
func (c *Client) LookupBusiness(ctx context.Context, req BusinessLookupRequest) (*BusinessLookupResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/enrich/business", req)
	if err != nil {
		return nil, err
	}
	var resp BusinessLookupResponse
	remarshal(data, &resp)
	populateRaw(data, &resp)
	return &resp, nil
}

// CheckSanctions screens an entity against global sanctions lists.
func (c *Client) CheckSanctions(ctx context.Context, req SanctionsRequest) (*SanctionsResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/risk/sanctions", req)
	if err != nil {
		return nil, err
	}
	var resp SanctionsResponse
	remarshal(data, &resp)
	populateRaw(data, &resp)
	return &resp, nil
}

// CheckDirectors checks for disqualified directors.
func (c *Client) CheckDirectors(ctx context.Context, req DirectorsRequest) (*DirectorsResponse, error) {
	data, err := c.doRequest(ctx, "POST", "/api/risk/directors", req)
	if err != nil {
		return nil, err
	}
	var resp DirectorsResponse
	remarshal(data, &resp)
	populateRaw(data, &resp)
	return &resp, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────

func remarshal(data map[string]interface{}, v interface{}) {
	b, _ := json.Marshal(data)
	json.Unmarshal(b, v)
}

func backoff(attempt int) time.Duration {
	return time.Duration(math.Pow(2, float64(attempt))) * time.Second
}

func sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}
