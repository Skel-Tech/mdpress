// Package cloud provides a client for interacting with the mdpress cloud API.
package cloud

// CloudTemplate represents a template available from the mdpress cloud.
type CloudTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Free        bool   `json:"free"`
}

// ListTemplatesResponse is the API response for listing available templates.
type ListTemplatesResponse struct {
	Templates  []CloudTemplate `json:"templates"`
	NextCursor string          `json:"next_cursor,omitempty"`
}

// FetchTemplateResponse is the API response for fetching a single template.
type FetchTemplateResponse struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
