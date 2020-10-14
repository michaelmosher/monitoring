package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michaelmosher/monitoring/pkg/metricly"
)

type metricsResponseData struct {
	Page struct {
		Content []metricly.Metric
	}
	NumberOfElements int
	Last             bool
}

func (s Service) FetchMetrics(query metricly.MetricQuery) ([]metricly.Metric, error) {
	req, err := s.createMetricsRequest(query)

	if err != nil {
		return nil, fmt.Errorf("error creating API request: %v", err)
	}

	resp, err := s.HTTPClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error executing API request: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, handleErrorResponse(resp, "metrics")
	}

	return handleMetricsResponse(resp)
}

func (s Service) createMetricsRequest(query metricly.MetricQuery) (*http.Request, error) {
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/metrics/elasticsearch/metricQuery", apiBaseURL),
		bytes.NewReader(queryBytes),
	)

	if err != nil {
		return nil, fmt.Errorf("error creating request object: %v", err)
	}

	req.Header.Add("Content-type", "application/json")
	req.SetBasicAuth(s.Username, s.Password)

	return req, nil
}

func handleMetricsResponse(resp *http.Response) ([]metricly.Metric, error) {
	defer resp.Body.Close()

	var d metricsResponseData
	err := json.NewDecoder(resp.Body).Decode(&d)

	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return d.Page.Content, nil
}
