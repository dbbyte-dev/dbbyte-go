package dbbyte

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.dbbyte.com"
	defaultTimeout = 30 * time.Second
	userAgent      = "dbbyte-go/0.1.0"
)

// ClientOptions configures a Client.
type ClientOptions struct {
	APIKey string
	// BaseURL defaults to https://api.dbbyte.com. In production the
	// Cloudflare Worker gateway in front of api.dbbyte.com inserts
	// /api/v1 itself, so paths passed internally do NOT include /v1. If
	// you point BaseURL at a Next.js dev server directly (bypassing the
	// Worker), include /api/v1 in BaseURL instead.
	BaseURL string
	Timeout time.Duration
	// HTTPClient lets you supply your own *http.Client (custom transport,
	// proxy, etc). Defaults to a client with Timeout applied.
	HTTPClient *http.Client
}

type client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func newClient(opts ClientOptions) (*client, error) {
	if opts.APIKey == "" {
		return nil, fmt.Errorf("dbbyte: APIKey is required")
	}
	baseURL := opts.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: timeout}
	}

	return &client{apiKey: opts.APIKey, baseURL: baseURL, httpClient: httpClient}, nil
}

type queryParams map[string]string

func (c *client) buildURL(path string, query queryParams) string {
	u := c.baseURL + path
	if len(query) == 0 {
		return u
	}
	values := url.Values{}
	for k, v := range query {
		if v != "" {
			values.Set(k, v)
		}
	}
	if len(values) == 0 {
		return u
	}
	return u + "?" + values.Encode()
}

func (c *client) newRequest(ctx context.Context, method, fullURL string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}

func (c *client) do(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		if ctxErr := req.Context().Err(); ctxErr != nil {
			return nil, &Error{Message: "Request timed out", Status: 0}
		}
		return nil, err
	}
	return res, nil
}

func (c *client) errorFromResponse(res *http.Response) error {
	defer res.Body.Close()

	var retryAfter int
	if ra := res.Header.Get("Retry-After"); ra != "" {
		if n, err := strconv.Atoi(ra); err == nil {
			retryAfter = n
		}
	}

	raw, _ := io.ReadAll(res.Body)
	var body struct {
		Error string `json:"error"`
		Hint  string `json:"hint"`
	}
	message := res.Status
	hint := ""
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &body); err == nil {
			if body.Error != "" {
				message = body.Error
			}
			hint = body.Hint
		}
	}

	return &Error{Message: message, Status: res.StatusCode, Hint: hint, RetryAfter: retryAfter}
}

// requestJSON sends a request with an optional JSON body and decodes a
// JSON response into out. Pass a nil out to discard the response body.
func (c *client) requestJSON(ctx context.Context, method, path string, query queryParams, body interface{}, out interface{}) error {
	fullURL := c.buildURL(path, query)

	var reader io.Reader
	var contentType string
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(encoded)
		contentType = "application/json"
	}

	req, err := c.newRequest(ctx, method, fullURL, reader)
	if err != nil {
		return err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	res, err := c.do(req)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return c.errorFromResponse(res)
	}
	defer res.Body.Close()

	if out == nil {
		return nil
	}
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, out)
}

// requestForm uploads a single file field named "file" as multipart form
// data and decodes a JSON response into out.
func (c *client) requestForm(ctx context.Context, method, path string, query queryParams, filename string, fileBytes []byte, contentType string, out interface{}) error {
	fullURL := c.buildURL(path, query)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}
	if _, err := part.Write(fileBytes); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	req, err := c.newRequest(ctx, method, fullURL, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := c.do(req)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return c.errorFromResponse(res)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(raw) == 0 || out == nil {
		return nil
	}
	return json.Unmarshal(raw, out)
}

// requestBinary performs a request and returns the raw response body plus
// Content-Type / Content-Disposition-derived filename.
func (c *client) requestBinary(ctx context.Context, method, path string, query queryParams) (*DownloadResult, error) {
	fullURL := c.buildURL(path, query)

	req, err := c.newRequest(ctx, method, fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, c.errorFromResponse(res)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	var filename *string
	if disposition := res.Header.Get("Content-Disposition"); disposition != "" {
		if name, ok := parseDispositionFilename(disposition); ok {
			filename = &name
		}
	}

	return &DownloadResult{Data: data, ContentType: contentType, Filename: filename}, nil
}

func parseDispositionFilename(disposition string) (string, bool) {
	const marker = `filename="`
	idx := strings.Index(disposition, marker)
	if idx == -1 {
		return "", false
	}
	rest := disposition[idx+len(marker):]
	end := strings.Index(rest, `"`)
	if end == -1 {
		return "", false
	}
	return rest[:end], true
}
