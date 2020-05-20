package cdc

import (
	"fmt"
	"log"
	"strings"

	"github.com/michaelmosher/monitoring/pkg/octopus"
)

type tenantMap map[string]octopus.Tenant
type octopusStatus struct {
	name   string
	online bool
}

func getOctopusStatuses(service octopus.Service) (map[string]octopusStatus, error) {
	tm, err := getOctopusTenants(service)

	if err != nil {
		return nil, fmt.Errorf("octopus.FetchTenants error: %s", err)
	}

	machines, err := service.FetchMachines()

	if err != nil {
		return nil, fmt.Errorf("octopus.FetchMachines error: %s", err)
	}

	statuses := make(map[string]octopusStatus)

	for _, m := range machines {
		if m.Status == "Disabled" || len(m.TenantIDs) == 0 {
			continue
		}

		statuses[getUAIDFromMachineName(m.Name)] = octopusStatus{
			online: (m.Status != "Offline"),
			name:   tm.findMachineName(m),
		}
	}

	return statuses, nil
}

func getOctopusTenants(service octopus.Service) (tenantMap, error) {
	tm := make(tenantMap)

	tenants, err := service.FetchTenants()

	if err != nil {
		return nil, err
	}

	for _, tenant := range tenants {
		tm[tenant.ID] = tenant
	}

	return tm, nil
}

func getUAIDFromMachineName(name string) string {
	return strings.Replace(name, "xx_polling-", "", 1)
}

func (tm tenantMap) findMachineName(m octopus.Machine) string {
	names := []string{}

	for t := range m.TenantIDs {
		names = append(names, tm[t].Name)
	}

	if len(names) == 0 || names[0] == "" {
		log.Printf("Couldn't find a name for %+v", m)
		return "unknown name"
	}

	return strings.Join(names, ", ")
}
