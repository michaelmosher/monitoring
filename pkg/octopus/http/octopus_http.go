package http

import (
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
