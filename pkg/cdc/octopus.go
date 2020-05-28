package cdc

import (
	"fmt"

	"github.com/michaelmosher/monitoring/pkg/octopus"
)

const (
	nucOctopusRole = "side-server-appliances"
	vmOctopusRole  = "linux-server"
)

func getOfflineNUCs(octo octopusClient) ([]octopus.Machine, error) {
	offlineNUCs := []octopus.Machine{}

	allMachines, err := octo.FetchMachines()

	if err != nil {
		return nil, fmt.Errorf("octopus.FetchMachines error: %s", err)
	}

	for _, machine := range allMachines {
		if machine.Status != "Offline" {
			continue
		}

		if _, ok := machine.Roles[nucOctopusRole]; ok == false {
			continue
		}

		offlineNUCs = append(offlineNUCs, machine)
	}

	return offlineNUCs, nil
}

func getOnlineMachines(octo octopusClient) ([]octopus.Machine, error) {
	onlineNUCs := []octopus.Machine{}

	allMachines, err := octo.FetchMachines()

	if err != nil {
		return nil, fmt.Errorf("octopus.FetchMachines error: %s", err)
	}

	for _, machine := range allMachines {
		if machine.Status == "Offline" {
			continue
		}

		_, sideServer := machine.Roles[nucOctopusRole]
		_, linuxServer := machine.Roles[vmOctopusRole]

		if !sideServer && !linuxServer {
			continue
		}

		onlineNUCs = append(onlineNUCs, machine)
	}

	return onlineNUCs, nil
}

func getOctopusTenants(octo octopusClient) (map[string]octopus.Tenant, error) {
	tm := make(map[string]octopus.Tenant)

	tenants, err := octo.FetchTenants()

	if err != nil {
		return nil, fmt.Errorf("octopus.FetchTenants error: %s", err)
	}

	for _, tenant := range tenants {
		tm[tenant.ID] = tenant
	}

	return tm, nil
}

func getOctopusProjectIDs(octo octopusClient, projectNames ...string) ([]string, error) {
	projectIDs := make([]string, 0, len(projectNames))

	for _, name := range projectNames {
		project, err := octo.FetchProject(name)

		if err != nil {
			return nil, fmt.Errorf("octopus.FetchProject(%s) error: %s", name, err)
		}

		projectIDs = append(projectIDs, project.ID)
	}

	return projectIDs, nil
}
