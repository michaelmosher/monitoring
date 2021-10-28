package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michaelmosher/monitoring/pkg/octopus"
)

func (s Service) FetchMachines() ([]octopus.Machine, error) {
	req, err := s.createDataRequest("machines/all")

	if err != nil {
		return nil, fmt.Errorf("error creating API request: %v", err)
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error executing API request: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, handleErrorResponse(resp, "machines")
	}

	return handleMachinesResponse(resp)
}

func (s Service) FetchMachine(machineID string) (octopus.Machine, error) {
	var m octopus.Machine

	if machineID == "" || machineID == "all" {
		return m, fmt.Errorf("no u. Use FetchMachines instead")
	}

	req, err := s.createDataRequest(fmt.Sprintf("machines/%s", machineID))

	if err != nil {
		return m, fmt.Errorf("error creating API request: %v", err)
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return m, fmt.Errorf("error executing API request: %v", err)
	}

	if resp.StatusCode != 200 {
		return m, handleErrorResponse(resp, "machine")
	}

	return m, handleMachineResponse(resp, &m)
}

func handleMachinesResponse(resp *http.Response) ([]octopus.Machine, error) {
	defer resp.Body.Close()

	type machineList struct {
		machines []octopus.Machine
	}

	list := machineList{}
	err := json.NewDecoder(resp.Body).Decode(&list.machines)
	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON: %v", err)
	}

	return list.machines, nil
}

func handleMachineResponse(resp *http.Response, m *octopus.Machine) error {
	defer resp.Body.Close()

	err := json.NewDecoder(resp.Body).Decode(m)

	if err != nil {
		return fmt.Errorf("Error decoding JSON: %v", err)
	}

	return nil
}
