package http

import (
	"net/http"
)

const apiBaseURL = "https://us.cloudwisdom.virtana.com"

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Service struct {
	HTTPClient httpClient
	Username   string
	Password   string
}
