package taqlyn

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Config configures a Taqlyn server API client.
type Config struct {
	BaseURL       string
	ClientID      string
	PrivateKeyPEM string
	HTTPClient    *http.Client
}

// Client calls Taqlyn's signed server-side short-link API.
// It must not be used in mobile or other untrusted client applications.
type Client struct {
	baseURL    string
	clientID   string
	privateKey ed25519.PrivateKey
	httpClient *http.Client
	now        func() time.Time
}

// NewClient creates a signed Taqlyn server API client.
func NewClient(config Config) (*Client, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("taqlyn: BaseURL is required")
	}
	if config.ClientID == "" {
		return nil, fmt.Errorf("taqlyn: ClientID is required")
	}
	if config.PrivateKeyPEM == "" {
		return nil, fmt.Errorf("taqlyn: PrivateKeyPEM is required")
	}

	baseURL, err := url.Parse(config.BaseURL)
	if err != nil || baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, fmt.Errorf("taqlyn: BaseURL must be an absolute URL")
	}
	if baseURL.RawQuery != "" || baseURL.Fragment != "" {
		return nil, fmt.Errorf("taqlyn: BaseURL must not contain a query or fragment")
	}
	baseURL.Path = strings.TrimRight(baseURL.Path, "/")

	privateKey, err := parsePrivateKey(config.PrivateKeyPEM)
	if err != nil {
		return nil, err
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL:    strings.TrimRight(baseURL.String(), "/"),
		clientID:   config.ClientID,
		privateKey: privateKey,
		httpClient: httpClient,
		now:        time.Now,
	}, nil
}

// CreateShortLink creates a short link.
func (c *Client) CreateShortLink(ctx context.Context, input CreateShortLinkRequest) (*ShortLink, error) {
	var result ShortLink
	if err := c.doJSON(ctx, http.MethodPost, "/v1/short-links", input, http.StatusCreated, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetShortLink retrieves a short link by ID.
func (c *Client) GetShortLink(ctx context.Context, id string) (*ShortLink, error) {
	var result ShortLink
	path := "/v1/short-links/" + url.PathEscape(id)
	if err := c.doJSON(ctx, http.MethodGet, path, nil, http.StatusOK, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PatchShortLink updates a short link by ID.
func (c *Client) PatchShortLink(ctx context.Context, id string, input PatchShortLinkRequest) (*ShortLink, error) {
	var result ShortLink
	path := "/v1/short-links/" + url.PathEscape(id)
	if err := c.doJSON(ctx, http.MethodPatch, path, input, http.StatusOK, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteShortLink deletes a short link by ID.
func (c *Client) DeleteShortLink(ctx context.Context, id string) error {
	path := "/v1/short-links/" + url.PathEscape(id)
	return c.doJSON(ctx, http.MethodDelete, path, nil, http.StatusNoContent, nil)
}

func (c *Client) doJSON(
	ctx context.Context,
	method string,
	path string,
	input any,
	expectedStatus int,
	output any,
) error {
	var body []byte
	var err error
	if input != nil {
		body, err = json.Marshal(input)
		if err != nil {
			return fmt.Errorf("taqlyn: encode request: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("taqlyn: create request: %w", err)
	}
	for name, values := range signedHeaders(c.privateKey, method, path, c.clientID, c.now().Unix(), body) {
		for _, value := range values {
			request.Header.Add(name, value)
		}
	}
	request.Header.Set("Accept", "application/json")
	if input != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("taqlyn: request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != expectedStatus {
		return decodeAPIError(response)
	}
	if output == nil {
		return nil
	}
	if err := json.NewDecoder(response.Body).Decode(output); err != nil {
		return fmt.Errorf("taqlyn: decode response: %w", err)
	}
	return nil
}

func decodeAPIError(response *http.Response) error {
	var payload struct {
		Code    string `json:"error"`
		Message string `json:"message"`
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err == nil {
		_ = json.Unmarshal(data, &payload)
	}
	if payload.Message == "" {
		payload.Message = strings.TrimSpace(string(data))
	}
	if payload.Message == "" {
		payload.Message = http.StatusText(response.StatusCode)
	}
	return &APIError{
		StatusCode: response.StatusCode,
		Code:       payload.Code,
		Message:    payload.Message,
	}
}
