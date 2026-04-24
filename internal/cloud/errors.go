package cloud

import "fmt"

// ErrTemplateNotFound is returned when a requested template does not exist.
type ErrTemplateNotFound struct {
	Name string
}

func (e *ErrTemplateNotFound) Error() string {
	return fmt.Sprintf("template not found: %s", e.Name)
}

// ErrUnauthorized is returned when the API request is not authorized.
type ErrUnauthorized struct {
	Message string
}

func (e *ErrUnauthorized) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("unauthorized: %s", e.Message)
	}
	return "unauthorized: API key required or invalid"
}

// ErrNetworkFailure is returned when a network error occurs.
type ErrNetworkFailure struct {
	Err error
}

func (e *ErrNetworkFailure) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("network failure: %s", e.Err.Error())
	}
	return "network failure"
}

func (e *ErrNetworkFailure) Unwrap() error {
	return e.Err
}
