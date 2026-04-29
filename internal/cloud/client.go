package cloud

import (
	"net/http"
	"os"
	"time"

	"github.com/skel-tech/mdpress/internal/auth"
)

const (
	// DefaultBaseURL is the default API endpoint for mdpress cloud.
	DefaultBaseURL = "https://api.mdpress.dev"

	// DefaultTimeout is the default HTTP timeout for API requests.
	DefaultTimeout = 10 * time.Second
)

// Client is an HTTP client for the mdpress cloud API.
type Client struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// NewClient creates a new cloud API client.
// It reads the MDPRESS_API_URL environment variable for the base URL,
// defaulting to https://api.mdpress.dev if not set.
func NewClient() *Client {
	baseURL := os.Getenv("MDPRESS_API_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		timeout: DefaultTimeout,
	}
}

// addAuthHeader adds the Authorization header to the request if an API key is available.
func (c *Client) addAuthHeader(req *http.Request) {
	creds, err := auth.Load()
	if err != nil || creds == nil {
		return
	}
	if creds.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+creds.APIKey)
	}
}
