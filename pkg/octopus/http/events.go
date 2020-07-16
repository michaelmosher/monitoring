package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michaelmosher/monitoring/pkg/octopus"
)

type EventsResponse struct {
	Events []octopus.Event `json:"Items"`
}

func (s Service) FetchEvents(filter map[string]string) ([]octopus.Event, error) {
	queryString := ""

	for key, value := range filter {
		queryString = fmt.Sprintf("%s&%s=%s", queryString, key, value)
	}

	req, err := s.createDataRequest(fmt.Sprintf("events?%s", queryString))

	if err != nil {
		return nil, fmt.Errorf("error creating API request: %v", err)
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error executing API request: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, handleErrorResponse(resp, "events")
	}

	return handleEventsResponse(resp)
}

func handleEventsResponse(resp *http.Response) ([]octopus.Event, error) {
	defer resp.Body.Close()

	var e EventsResponse
	err := json.NewDecoder(resp.Body).Decode(&e)

	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON: %v", err)
	}

	return e.Events, nil
}
