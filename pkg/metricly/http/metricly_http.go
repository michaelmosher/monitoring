package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const apiBaseURL = "https://us.cloudwisdom.virtana.com"

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type errorResponse struct {
	StatusCode int
	Error      string
	Message    string
	Path       string
}

type Service struct {
	HTTPClient httpClient
	Username   string
	Password   string
}

func handleErrorResponse(resp *http.Response, caller string) error {
	defer resp.Body.Close()

	var e errorResponse
	err := json.NewDecoder(resp.Body).Decode(&e)
	e.StatusCode = resp.StatusCode

	if err != nil {
		return fmt.Errorf("Error decoding JSON: %v", err)
	}

	return fmt.Errorf("Error retrieving %s data: %+v", caller, e)
}
