package octopus

import "encoding/json"

type Machine struct {
	Name      string
	Status    string
	Roles     map[string]struct{}
	TenantIDs map[string]struct{}
}

func (m *Machine) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	m.Name = v["Name"].(string)
	m.Status = v["Status"].(string)
	m.Roles = make(map[string]struct{})
	m.TenantIDs = make(map[string]struct{})

	for _, k := range v["Roles"].([]interface{}) {
		m.Roles[k.(string)] = struct{}{}
	}

	for _, k := range v["TenantIds"].([]interface{}) {
		m.TenantIDs[k.(string)] = struct{}{}
	}

	return nil
}

type Project struct {
	ID   string
	Name string
}

type Tenant struct {
	ID         string
	Name       string
	ProjectIDs map[string]struct{}
}

func (t *Tenant) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	t.ID = v["Id"].(string)
	t.Name = v["Name"].(string)
	t.ProjectIDs = make(map[string]struct{})

	for k := range v["ProjectEnvironments"].(map[string]interface{}) {
		t.ProjectIDs[k] = struct{}{}
	}

	return nil
}

type client interface {
	FetchMachines() ([]Machine, error)
	FetchMachine(machineID string) (Machine, error)

	FetchProjects() ([]Project, error)
	FetchProject(projectID string) (Project, error)

	FetchTenants() ([]Tenant, error)
	FetchTenant(tenantID string) (Tenant, error)
}

type Service struct {
	client client
}

func New(client client) Service {
	return Service{
		client: client,
	}
}

func (s Service) FetchMachines() ([]Machine, error) {
	return s.client.FetchMachines()
}

func (s Service) FetchMachine(machineID string) (Machine, error) {
	return s.client.FetchMachine(machineID)
}

func (s Service) FetchProjects() ([]Project, error) {
	return s.client.FetchProjects()
}

func (s Service) FetchProject(projectID string) (Project, error) {
	return s.client.FetchProject(projectID)
}

func (s Service) FetchTenants() ([]Tenant, error) {
	return s.client.FetchTenants()
}

func (s Service) FetchTenant(tenantID string) (Tenant, error) {
	return s.client.FetchTenant(tenantID)
}
