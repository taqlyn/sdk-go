package taqlyn

import "fmt"

// APIError represents a non-success response from the Taqlyn API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("taqlyn: API request failed (%d, %s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("taqlyn: API request failed (%d): %s", e.StatusCode, e.Message)
}
