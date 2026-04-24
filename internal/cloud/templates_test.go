package cloud

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- ListTemplates Tests ---

func TestListTemplates_SinglePage(t *testing.T) {
	templates := []CloudTemplate{
		{Name: "basic", Description: "Basic template", Free: true},
		{Name: "pro", Description: "Pro template", Free: false},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/templates" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		resp := ListTemplatesResponse{Templates: templates}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	result, err := client.ListTemplates()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 templates, got %d", len(result))
	}
	if result[0].Name != "basic" {
		t.Errorf("expected first template name 'basic', got %q", result[0].Name)
	}
	if result[1].Name != "pro" {
		t.Errorf("expected second template name 'pro', got %q", result[1].Name)
	}
}

func TestListTemplates_MultiplePages(t *testing.T) {
	page1Templates := []CloudTemplate{
		{Name: "template1", Description: "First", Free: true},
		{Name: "template2", Description: "Second", Free: true},
	}
	page2Templates := []CloudTemplate{
		{Name: "template3", Description: "Third", Free: false},
	}

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		cursor := r.URL.Query().Get("cursor")

		var resp ListTemplatesResponse
		if cursor == "" {
			// First page
			resp = ListTemplatesResponse{
				Templates:  page1Templates,
				NextCursor: "cursor-page-2",
			}
		} else if cursor == "cursor-page-2" {
			// Second page
			resp = ListTemplatesResponse{
				Templates:  page2Templates,
				NextCursor: "", // No more pages
			}
		} else {
			t.Errorf("unexpected cursor: %s", cursor)
			http.Error(w, "bad cursor", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	result, err := client.ListTemplates()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if requestCount != 2 {
		t.Errorf("expected 2 requests for pagination, got %d", requestCount)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 templates total, got %d", len(result))
	}
	if result[0].Name != "template1" {
		t.Errorf("expected first template 'template1', got %q", result[0].Name)
	}
	if result[2].Name != "template3" {
		t.Errorf("expected third template 'template3', got %q", result[2].Name)
	}
}

func TestListTemplates_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ListTemplatesResponse{Templates: []CloudTemplate{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	result, err := client.ListTemplates()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d templates", len(result))
	}
}

func TestListTemplates_NetworkError(t *testing.T) {
	// Create a server and immediately close it to simulate network error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	serverURL := server.URL
	server.Close()

	client := &Client{
		baseURL:    serverURL,
		httpClient: &http.Client{},
	}

	_, err := client.ListTemplates()
	if err == nil {
		t.Fatal("expected error for network failure")
	}

	var netErr *ErrNetworkFailure
	if !isNetworkError(err, &netErr) {
		t.Errorf("expected ErrNetworkFailure, got %T: %v", err, err)
	}
}

func TestListTemplates_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.ListTemplates()
	if err == nil {
		t.Fatal("expected error for unauthorized")
	}

	var authErr *ErrUnauthorized
	if !isUnauthorizedError(err, &authErr) {
		t.Errorf("expected ErrUnauthorized, got %T: %v", err, err)
	}
}

func TestListTemplates_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json{"))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.ListTemplates()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}

	var netErr *ErrNetworkFailure
	if !isNetworkError(err, &netErr) {
		t.Errorf("expected ErrNetworkFailure for decode error, got %T: %v", err, err)
	}
}

func TestListTemplates_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.ListTemplates()
	if err == nil {
		t.Fatal("expected error for server error")
	}

	var netErr *ErrNetworkFailure
	if !isNetworkError(err, &netErr) {
		t.Errorf("expected ErrNetworkFailure, got %T: %v", err, err)
	}
}

// --- FetchTemplate Tests ---

func TestFetchTemplate_Success(t *testing.T) {
	expectedContent := "title: My Template\nslides:\n  - content: Hello"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/templates/basic" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		resp := FetchTemplateResponse{
			Name:    "basic",
			Content: expectedContent,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	content, err := client.FetchTemplate("basic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}
}

func TestFetchTemplate_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.FetchTemplate("nonexistent")
	if err == nil {
		t.Fatal("expected error for not found")
	}

	var notFoundErr *ErrTemplateNotFound
	if !isTemplateNotFoundError(err, &notFoundErr) {
		t.Errorf("expected ErrTemplateNotFound, got %T: %v", err, err)
	}
	if notFoundErr.Name != "nonexistent" {
		t.Errorf("expected template name 'nonexistent', got %q", notFoundErr.Name)
	}
}

func TestFetchTemplate_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.FetchTemplate("pro-template")
	if err == nil {
		t.Fatal("expected error for unauthorized")
	}

	var authErr *ErrUnauthorized
	if !isUnauthorizedError(err, &authErr) {
		t.Errorf("expected ErrUnauthorized, got %T: %v", err, err)
	}
}

func TestFetchTemplate_NetworkError(t *testing.T) {
	// Create a server and immediately close it to simulate network error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	serverURL := server.URL
	server.Close()

	client := &Client{
		baseURL:    serverURL,
		httpClient: &http.Client{},
	}

	_, err := client.FetchTemplate("basic")
	if err == nil {
		t.Fatal("expected error for network failure")
	}

	var netErr *ErrNetworkFailure
	if !isNetworkError(err, &netErr) {
		t.Errorf("expected ErrNetworkFailure, got %T: %v", err, err)
	}
}

func TestFetchTemplate_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.FetchTemplate("basic")
	if err == nil {
		t.Fatal("expected error for server error")
	}

	var netErr *ErrNetworkFailure
	if !isNetworkError(err, &netErr) {
		t.Errorf("expected ErrNetworkFailure, got %T: %v", err, err)
	}
}

func TestFetchTemplate_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not valid json"))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.FetchTemplate("basic")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}

	var netErr *ErrNetworkFailure
	if !isNetworkError(err, &netErr) {
		t.Errorf("expected ErrNetworkFailure for decode error, got %T: %v", err, err)
	}
}

func TestFetchTemplate_URLEncoding(t *testing.T) {
	// Test that template names with special characters are properly URL-encoded
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.URL.RawPath contains the encoded path (r.URL.Path is decoded)
		// Check that the raw path was properly encoded
		expectedRaw := "/v1/templates/my%20template%2Fwith%20special"
		if r.URL.RawPath != expectedRaw {
			t.Errorf("unexpected raw path: got %q, want %q", r.URL.RawPath, expectedRaw)
		}

		resp := FetchTemplateResponse{
			Name:    "my template/with special",
			Content: "content",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.FetchTemplate("my template/with special")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Auth Header Tests ---

func TestListTemplates_AuthHeaderIncluded(t *testing.T) {
	// Set up temp home directory with auth credentials
	tempHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	// Create auth file
	authDir := filepath.Join(tempHome, ".config", "mdpress")
	if err := os.MkdirAll(authDir, 0755); err != nil {
		t.Fatal(err)
	}
	authFile := filepath.Join(authDir, "auth.yml")
	authContent := "api_key: mdp_test_secret_key\n"
	if err := os.WriteFile(authFile, []byte(authContent), 0600); err != nil {
		t.Fatal(err)
	}

	var receivedAuthHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get("Authorization")
		resp := ListTemplatesResponse{Templates: []CloudTemplate{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.ListTemplates()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedHeader := "Bearer mdp_test_secret_key"
	if receivedAuthHeader != expectedHeader {
		t.Errorf("expected Authorization header %q, got %q", expectedHeader, receivedAuthHeader)
	}
}

func TestListTemplates_AuthHeaderOmittedWhenNoCredentials(t *testing.T) {
	// Set up temp home directory without auth credentials
	tempHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	var receivedAuthHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get("Authorization")
		resp := ListTemplatesResponse{Templates: []CloudTemplate{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.ListTemplates()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedAuthHeader != "" {
		t.Errorf("expected no Authorization header, got %q", receivedAuthHeader)
	}
}

func TestFetchTemplate_AuthHeaderIncluded(t *testing.T) {
	// Set up temp home directory with auth credentials
	tempHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	// Create auth file
	authDir := filepath.Join(tempHome, ".config", "mdpress")
	if err := os.MkdirAll(authDir, 0755); err != nil {
		t.Fatal(err)
	}
	authFile := filepath.Join(authDir, "auth.yml")
	authContent := "api_key: mdp_another_key\n"
	if err := os.WriteFile(authFile, []byte(authContent), 0600); err != nil {
		t.Fatal(err)
	}

	var receivedAuthHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get("Authorization")
		resp := FetchTemplateResponse{Name: "test", Content: "content"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.FetchTemplate("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedHeader := "Bearer mdp_another_key"
	if receivedAuthHeader != expectedHeader {
		t.Errorf("expected Authorization header %q, got %q", expectedHeader, receivedAuthHeader)
	}
}

func TestFetchTemplate_AuthHeaderOmittedWhenNoCredentials(t *testing.T) {
	// Set up temp home directory without auth credentials
	tempHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	var receivedAuthHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get("Authorization")
		resp := FetchTemplateResponse{Name: "test", Content: "content"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.FetchTemplate("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedAuthHeader != "" {
		t.Errorf("expected no Authorization header, got %q", receivedAuthHeader)
	}
}

func TestAuthHeader_EmptyAPIKey(t *testing.T) {
	// Set up temp home directory with empty API key
	tempHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	// Create auth file with empty api_key
	authDir := filepath.Join(tempHome, ".config", "mdpress")
	if err := os.MkdirAll(authDir, 0755); err != nil {
		t.Fatal(err)
	}
	authFile := filepath.Join(authDir, "auth.yml")
	authContent := "api_key: \"\"\nemail: user@example.com\n"
	if err := os.WriteFile(authFile, []byte(authContent), 0600); err != nil {
		t.Fatal(err)
	}

	var receivedAuthHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get("Authorization")
		resp := ListTemplatesResponse{Templates: []CloudTemplate{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.ListTemplates()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty API key should not set auth header
	if receivedAuthHeader != "" {
		t.Errorf("expected no Authorization header for empty API key, got %q", receivedAuthHeader)
	}
}

// --- Helper functions for error type checking ---

func isNetworkError(err error, target **ErrNetworkFailure) bool {
	if e, ok := err.(*ErrNetworkFailure); ok {
		*target = e
		return true
	}
	return false
}

func isUnauthorizedError(err error, target **ErrUnauthorized) bool {
	if e, ok := err.(*ErrUnauthorized); ok {
		*target = e
		return true
	}
	return false
}

func isTemplateNotFoundError(err error, target **ErrTemplateNotFound) bool {
	if e, ok := err.(*ErrTemplateNotFound); ok {
		*target = e
		return true
	}
	return false
}

// --- Integration-style tests for custom base URL ---

func TestCustomBaseURL_ListTemplates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ListTemplatesResponse{
			Templates: []CloudTemplate{{Name: "custom", Description: "From custom URL", Free: true}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	oldURL := os.Getenv("MDPRESS_API_URL")
	os.Setenv("MDPRESS_API_URL", server.URL)
	defer func() {
		if oldURL != "" {
			os.Setenv("MDPRESS_API_URL", oldURL)
		} else {
			os.Unsetenv("MDPRESS_API_URL")
		}
	}()

	client := NewClient()
	// Override httpClient to use test server's client
	client.httpClient = server.Client()

	result, err := client.ListTemplates()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 template, got %d", len(result))
	}
	if result[0].Name != "custom" {
		t.Errorf("expected template name 'custom', got %q", result[0].Name)
	}
}

func TestCustomBaseURL_FetchTemplate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/templates/") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := FetchTemplateResponse{Name: "custom", Content: "custom content"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	oldURL := os.Getenv("MDPRESS_API_URL")
	os.Setenv("MDPRESS_API_URL", server.URL)
	defer func() {
		if oldURL != "" {
			os.Setenv("MDPRESS_API_URL", oldURL)
		} else {
			os.Unsetenv("MDPRESS_API_URL")
		}
	}()

	client := NewClient()
	// Override httpClient to use test server's client
	client.httpClient = server.Client()

	content, err := client.FetchTemplate("custom")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != "custom content" {
		t.Errorf("expected content 'custom content', got %q", content)
	}
}
