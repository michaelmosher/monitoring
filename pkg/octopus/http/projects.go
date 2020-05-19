package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michaelmosher/monitoring/pkg/octopus"
)

func (s Service) FetchProjects() ([]octopus.Project, error) {
	req, err := s.createDataRequest("projects/all")

	if err != nil {
		return nil, fmt.Errorf("error creating API request: %v", err)
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error executing API request: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, handleErrorResponse(resp, "projects")
	}

	return handleProjectsResponse(resp)

}

func (s Service) FetchProject(projectID string) (octopus.Project, error) {
	var p octopus.Project

	if projectID == "" || projectID == "all" {
		return p, fmt.Errorf("no u. Use FetchProjects instead")
	}

	req, err := s.createDataRequest(fmt.Sprintf("projects/%s", projectID))

	if err != nil {
		return p, fmt.Errorf("error creating API request: %v", err)
	}

	resp, err := s.httpClient.Do(req)

	if err != nil {
		return p, fmt.Errorf("error executing API request: %v", err)
	}

	if resp.StatusCode != 200 {
		return p, handleErrorResponse(resp, "project")
	}

	return p, handleProjectResponse(resp, &p)

}

func handleProjectsResponse(resp *http.Response) ([]octopus.Project, error) {
	defer resp.Body.Close()

	list := make([]octopus.Project, 0)
	dec := json.NewDecoder(resp.Body)

	// throw away opening '['
	dec.Token()

	for dec.More() {
		var p octopus.Project
		err := dec.Decode(&p)

		if err != nil {
			return list, fmt.Errorf("Error decoding JSON: %v", err)
		}

		list = append(list, p)
	}

	// throw away closing ']'
	dec.Token()

	return list, nil
}

func handleProjectResponse(resp *http.Response, p *octopus.Project) error {
	defer resp.Body.Close()

	err := json.NewDecoder(resp.Body).Decode(p)

	if err != nil {
		return fmt.Errorf("Error decoding JSON: %v", err)
	}

	return nil
}
