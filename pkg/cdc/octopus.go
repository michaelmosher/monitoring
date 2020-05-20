package cdc

import (
	"fmt"
	"net/http"

	"github.com/michaelmosher/monitoring/pkg/octopus"
	octopus_http "github.com/michaelmosher/monitoring/pkg/octopus/http"
)

func octopusHTTPService(httpClient *http.Client, url string, space string, apiKey string) octopus.Service {
	return octopus.New(
		octopus_http.New(httpClient, url, space, apiKey),
	)
}

func (s Service) getOfflineNUCs() ([]octopus.Machine, error) {
	offlineNUCs := []octopus.Machine{}

	allMachines, err := s.Octo.FetchMachines()

	if err != nil {
		return nil, fmt.Errorf("octopus.FetchMachines error: %s", err)
	}

	for _, machine := range allMachines {
		if machine.Status != "Offline" {
			continue
		}

		if _, ok := machine.Roles["side-server-appliances"]; ok == false {
			continue
		}

		offlineNUCs = append(offlineNUCs, machine)
	}

	return offlineNUCs, nil
}

func (s Service) getOctopusTenants() (tenantMap, error) {
	tm := make(tenantMap)

	tenants, err := s.Octo.FetchTenants()

	if err != nil {
		return nil, fmt.Errorf("octopus.FetchTenants error: %s", err)
	}

	for _, tenant := range tenants {
		tm[tenant.ID] = tenant
	}

	return tm, nil
}

func (s Service) getOctopusProjectIDs(projectNames ...string) ([]string, error) {
	projectIDs := make([]string, 0, len(projectNames))

	for _, name := range projectNames {
		project, err := s.Octo.FetchProject(name)

		if err != nil {
			return nil, fmt.Errorf("octopus.FetchProject(%s) error: %s", name, err)
		}

		projectIDs = append(projectIDs, project.ID)
	}

	return projectIDs, nil
}
