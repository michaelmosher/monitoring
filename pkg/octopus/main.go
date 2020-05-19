package octopus

import "encoding/json"

type Machine struct {
	Name      string
	Status    string
	Roles     []string
	TenantIds []string
}

type Project struct {
	ID   string
	Name string
}

type Tenant struct {
	ID         string
	Name       string
	ProjectIDs []string
}

func (t *Tenant) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	m := v.(map[string]interface{})

	t.ID = m["Id"].(string)
	t.Name = m["Name"].(string)
	t.ProjectIDs = []string{}

	for k := range m["ProjectEnvironments"].(map[string]interface{}) {
		t.ProjectIDs = append(t.ProjectIDs, k)
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
