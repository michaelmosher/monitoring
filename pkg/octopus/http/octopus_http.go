package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type httpDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type errorResponse struct {
	StatusCode   int
	ErrorMessage string
}

type Service struct {
	httpClient httpDoer
	apiBaseURL string
	apiKey     string
}

// New creates an instance of an Octopus client, ready to call some APIs.
// Because the `octopus` struct is private, this is the only public way
// to obtain an instance.
func New(doer httpDoer, instanceURL string, space string, apiKey string) Service {
	return Service{
		httpClient: doer,
		apiBaseURL: fmt.Sprintf("%s/api/%s", instanceURL, space),
		apiKey:     apiKey,
	}
}

func (s Service) createDataRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%s", s.apiBaseURL, url),
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating request object: %v", err)
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("X-Octopus-ApiKey", s.apiKey)

	return req, nil
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
