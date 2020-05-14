package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michaelmosher/monitoring/pkg/octopus"
)

func (s Service) FetchMachines() ([]octopus.Machine, error) {
	req, err := s.createMachineDataRequest("all")

	if err != nil {
		return nil, fmt.Errorf("error creating API request: %v", err)
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error executing API request: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, handleMachineErrorRepsonse(resp)
	}

	return handleMachinesResponse(resp)
}

func (s Service) FetchMachine(machineID string) (octopus.Machine, error) {
	var m octopus.Machine

	if machineID == "" || machineID == "all" {
		return m, fmt.Errorf("no u. Use FetchMachines instead")
	}

	req, err := s.createMachineDataRequest(machineID)

	if err != nil {
		return m, fmt.Errorf("error creating API request: %v", err)
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return m, fmt.Errorf("error executing API request: %v", err)
	}

	if resp.StatusCode != 200 {
		return m, handleMachineErrorRepsonse(resp)
	}

	return m, handleMachineResponse(resp, &m)
}

func (s Service) createMachineDataRequest(machineID string) (*http.Request, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/machines/%s", s.apiBaseURL, machineID),
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating request object: %v", err)
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("X-Octopus-ApiKey", s.apiKey)

	return req, nil
}

func handleMachineErrorRepsonse(resp *http.Response) error {
	defer resp.Body.Close()

	var e errorResponse
	err := json.NewDecoder(resp.Body).Decode(&e)
	e.StatusCode = resp.StatusCode

	if err != nil {
		return fmt.Errorf("Error decoding JSON: %v", err)
	}

	return fmt.Errorf("Error retrieving machine data: %+v", e)
}

func handleMachinesResponse(resp *http.Response) ([]octopus.Machine, error) {
	defer resp.Body.Close()

	list := make([]octopus.Machine, 0)
	dec := json.NewDecoder(resp.Body)

	// throw away opening '['
	dec.Token()

	for dec.More() {
		var m octopus.Machine
		err := dec.Decode(&m)

		if err != nil {
			return list, fmt.Errorf("Error decoding JSON: %v", err)
		}

		list = append(list, m)
	}

	// throw away closing ']'
	dec.Token()

	return list, nil
}

func handleMachineResponse(resp *http.Response, m *octopus.Machine) error {
	defer resp.Body.Close()

	err := json.NewDecoder(resp.Body).Decode(m)

	if err != nil {
		return fmt.Errorf("Error decoding JSON: %v", err)
	}

	return nil
}
