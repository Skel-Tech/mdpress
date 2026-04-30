package cloud

import (
	"os"
	"testing"
)

func TestNewClient_DefaultBaseURL(t *testing.T) {
	// Clear any existing env var
	oldURL := os.Getenv("MDPRESS_API_URL")
	os.Unsetenv("MDPRESS_API_URL")
	defer func() {
		if oldURL != "" {
			os.Setenv("MDPRESS_API_URL", oldURL)
		}
	}()

	client := NewClient()
	if client.baseURL != DefaultBaseURL {
		t.Errorf("expected baseURL %q, got %q", DefaultBaseURL, client.baseURL)
	}
}

func TestNewClient_CustomBaseURL(t *testing.T) {
	customURL := "https://custom.api.example.com"

	oldURL := os.Getenv("MDPRESS_API_URL")
	os.Setenv("MDPRESS_API_URL", customURL)
	defer func() {
		if oldURL != "" {
			os.Setenv("MDPRESS_API_URL", oldURL)
		} else {
			os.Unsetenv("MDPRESS_API_URL")
		}
	}()

	client := NewClient()
	if client.baseURL != customURL {
		t.Errorf("expected baseURL %q, got %q", customURL, client.baseURL)
	}
}

func TestNewClient_DefaultTimeout(t *testing.T) {
	client := NewClient()
	if client.timeout != DefaultTimeout {
		t.Errorf("expected timeout %v, got %v", DefaultTimeout, client.timeout)
	}
	if client.httpClient.Timeout != DefaultTimeout {
		t.Errorf("expected httpClient timeout %v, got %v", DefaultTimeout, client.httpClient.Timeout)
	}
}

func TestNewClient_HttpClientNotNil(t *testing.T) {
	client := NewClient()
	if client.httpClient == nil {
		t.Error("expected httpClient to be non-nil")
	}
}
