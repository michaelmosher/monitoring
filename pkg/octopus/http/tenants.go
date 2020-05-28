package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michaelmosher/monitoring/pkg/octopus"
)

func (s Service) FetchTenants() ([]octopus.Tenant, error) {
	req, err := s.createDataRequest("tenantvariables/all")

	if err != nil {
		return nil, fmt.Errorf("error creating API request: %v", err)
	}

	resp, err := s.httpClient.Do(req)

	if resp.StatusCode != 200 {
		return nil, handleErrorResponse(resp, "tenants")
	}

	if err != nil {
		return nil, fmt.Errorf("error executing API request: %v", err)
	}

	return handleTenantsResponse(resp)
}

func (s Service) FetchTenant(tenantID string) (octopus.Tenant, error) {
	var t octopus.Tenant

	if tenantID == "" || tenantID == "all" {
		return t, fmt.Errorf("no u. Use FetchTenants instead")
	}

	req, err := s.createDataRequest(fmt.Sprintf("tenants/%s", tenantID))

	if err != nil {
		return t, fmt.Errorf("error creating API request: %v", err)
	}

	resp, err := s.httpClient.Do(req)

	if resp.StatusCode != 200 {
		return t, handleErrorResponse(resp, "tenant")
	}

	if err != nil {
		return t, fmt.Errorf("error executing API request: %v", err)
	}

	return t, handleTenantResponse(resp, &t)
}

func handleTenantsResponse(resp *http.Response) ([]octopus.Tenant, error) {
	defer resp.Body.Close()

	list := make([]octopus.Tenant, 0)
	dec := json.NewDecoder(resp.Body)

	// throw away opening '['
	dec.Token()

	for dec.More() {
		var t octopus.Tenant
		err := dec.Decode(&t)

		if err != nil {
			return list, fmt.Errorf("Error decoding JSON: %v", err)
		}

		list = append(list, t)
	}

	// throw away closing ']'
	dec.Token()

	return list, nil
}

func handleTenantResponse(resp *http.Response, t *octopus.Tenant) error {
	defer resp.Body.Close()

	err := json.NewDecoder(resp.Body).Decode(t)

	if err != nil {
		return fmt.Errorf("Error decoding JSON: %v", err)
	}

	return nil
}
