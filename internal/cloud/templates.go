package cloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ListTemplates fetches all available templates from the cloud API.
// It automatically handles cursor-based pagination, fetching all pages.
func (c *Client) ListTemplates() ([]CloudTemplate, error) {
	var allTemplates []CloudTemplate
	cursor := ""

	for {
		templates, nextCursor, err := c.listTemplatesPage(cursor)
		if err != nil {
			return nil, err
		}

		allTemplates = append(allTemplates, templates...)

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	return allTemplates, nil
}

// listTemplatesPage fetches a single page of templates.
func (c *Client) listTemplatesPage(cursor string) ([]CloudTemplate, string, error) {
	endpoint := c.baseURL + "/v1/templates"

	if cursor != "" {
		endpoint += "?cursor=" + url.QueryEscape(cursor)
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", &ErrNetworkFailure{Err: err}
	}

	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", &ErrNetworkFailure{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, "", &ErrUnauthorized{}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", &ErrNetworkFailure{
			Err: fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body)),
		}
	}

	var listResp ListTemplatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, "", &ErrNetworkFailure{Err: fmt.Errorf("failed to decode response: %w", err)}
	}

	return listResp.Templates, listResp.NextCursor, nil
}

// FetchTemplate retrieves a single template's content by name.
// Returns the template content as YAML string.
func (c *Client) FetchTemplate(name string) (string, error) {
	endpoint := c.baseURL + "/v1/templates/" + url.PathEscape(name)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", &ErrNetworkFailure{Err: err}
	}

	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", &ErrNetworkFailure{Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Success, continue to parse response
	case http.StatusNotFound:
		return "", &ErrTemplateNotFound{Name: name}
	case http.StatusUnauthorized:
		return "", &ErrUnauthorized{}
	default:
		body, _ := io.ReadAll(resp.Body)
		return "", &ErrNetworkFailure{
			Err: fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body)),
		}
	}

	var fetchResp FetchTemplateResponse
	if err := json.NewDecoder(resp.Body).Decode(&fetchResp); err != nil {
		return "", &ErrNetworkFailure{Err: fmt.Errorf("failed to decode response: %w", err)}
	}

	return fetchResp.Content, nil
}
