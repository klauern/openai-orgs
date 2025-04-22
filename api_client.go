package openaiorgs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// DefaultBaseURL is the default endpoint for the OpenAI Organizations API.
const DefaultBaseURL = "https://api.openai.com/v1"

// Client represents an OpenAI Organizations API client.
// It handles authentication, request retries, and provides methods for interacting with the API.
type Client struct {
	// client is the underlying HTTP client used for making requests.
	client *resty.Client
	// BaseURL is the base endpoint for API requests.
	BaseURL string
}

// withLowRateLimits configures a resty client with conservative rate limiting settings.
// It sets up retry behavior for 5xx errors and rate limit (429) responses.
func withLowRateLimits(client *resty.Client) *resty.Client {
	return client.
		SetRetryCount(20).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(5 * time.Minute).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			return response.StatusCode() >= 500 || response.StatusCode() == 429
		})
}

// NewClient creates a new OpenAI Organizations API client.
// It configures the client with the provided base URL and authentication token.
// If baseURL is empty, DefaultBaseURL is used.
func NewClient(baseURL, token string) *Client {
	client := resty.New()
	client = withLowRateLimits(client)

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	client.SetBaseURL(baseURL)
	client.SetAuthToken(token)
	client.SetHeader("Content-Type", "application/json")

	return &Client{
		client:  client,
		BaseURL: baseURL,
	}
}

// GetSingle makes a GET request to retrieve a single resource of type T.
// It handles JSON unmarshaling and error wrapping.
// Returns a pointer to the resource and any error encountered.
func GetSingle[T any](client *resty.Client, endpoint string) (*T, error) {
	resp, err := client.R().
		ExpectContentType("application/json").
		Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	var result T
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

// Get makes a GET request to retrieve a list of resources of type T.
// It supports query parameters and handles pagination through the ListResponse type.
// Returns a pointer to the ListResponse containing the resources and any error encountered.
func Get[T any](client *resty.Client, endpoint string, queryParams map[string]string) (*ListResponse[T], error) {
	resp, err := client.R().
		SetQueryParams(queryParams).
		ExpectContentType("application/json").
		Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	var listResp ListResponse[T]
	err = json.Unmarshal(resp.Body(), &listResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &listResp, nil
}

// Post makes a POST request to create a new resource of type T.
// It handles JSON marshaling of the request body and unmarshaling of the response.
// Returns a pointer to the created resource and any error encountered.
func Post[T any](client *resty.Client, endpoint string, body any) (*T, error) {
	resp, err := client.R().
		SetBody(body).
		ExpectContentType("application/json").
		Post(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	var result T
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

// Delete makes a DELETE request to remove a resource.
// It returns an error if the request fails or returns a non-2xx status code.
func Delete[T any](client *resty.Client, endpoint string) error {
	resp, err := client.R().Delete(endpoint)
	if err != nil {
		return fmt.Errorf("error making DELETE request: %v", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}
