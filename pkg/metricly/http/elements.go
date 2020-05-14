package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michaelmosher/monitoring/pkg/metricly"
)

type sampleResponseData struct {
	Samples []struct {
		Data struct {
			Val float64
		}
	}
}

// FetchMetricValue returns the latest value of a given time-series data metric.
func (s Service) FetchMetricValue(metric metricly.Metric) (float64, error) {
	req, err := s.createSampleRequest(metric)

	if err != nil {
		return 0, fmt.Errorf("error creating API request: %v", err)
	}

	resp, err := s.HTTPClient.Do(req)

	if err != nil {
		return 0, fmt.Errorf("error executing API request: %v", err)
	}

	return handleSampleResponse(resp)
}

func (s Service) createSampleRequest(metric metricly.Metric) (*http.Request, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/elements/%s/metrics/%s/samples",
			apiBaseURL,
			metric.ElementID,
			metric.ID),
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating request object: %v", err)
	}

	req.Header.Add("Content-type", "application/json")
	req.SetBasicAuth(s.Username, s.Password)

	q := req.URL.Query()
	q.Add("duration", "PT1M")
	q.Add("rollup", "ZERO")

	req.URL.RawQuery = q.Encode()

	return req, nil
}

func handleSampleResponse(resp *http.Response) (float64, error) {
	defer resp.Body.Close()

	var d sampleResponseData
	err := json.NewDecoder(resp.Body).Decode(&d)

	if err != nil {
		return 0, fmt.Errorf("error decoding JSON: %v", err)
	}

	if len(d.Samples) == 0 {
		return 0, fmt.Errorf("0 results from getMetricResults API")
	}

	return d.Samples[0].Data.Val, nil
}
