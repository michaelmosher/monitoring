package octopus

import (
	"encoding/json"
	"time"
)

type Machine struct {
	ID        string
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

	m.ID = v["Id"].(string)
	m.Name = v["Name"].(string)
	m.Status = v["HealthStatus"].(string)
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
	Variables  map[string]string
}

func (t *Tenant) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if id, ok := v["TenantId"]; ok {
		t.ID = id.(string)
	} else {
		t.ID = v["Id"].(string)
	}

	if name, ok := v["TenantName"]; ok {
		t.Name = name.(string)
	} else {
		t.Name = v["Name"].(string)
	}

	t.ProjectIDs = make(map[string]struct{})

	var projects map[string]interface{}

	if _, ok := v["ProjectEnvironments"]; ok {
		projects = v["ProjectEnvironments"].(map[string]interface{})
	} else {
		projects = v["ProjectVariables"].(map[string]interface{})
	}

	for k := range projects {
		t.ProjectIDs[k] = struct{}{}
	}

	t.Variables = make(map[string]string)

	if libraries, ok := v["LibraryVariables"]; ok {
		libraries := libraries.(map[string]interface{})

		for _, variableSet := range libraries {
			variableSet := variableSet.(map[string]interface{})

			if templates, ok := variableSet["Templates"]; ok {
				templates := templates.([]interface{})

				for _, template := range templates {
					template := template.(map[string]interface{})

					varName := template["Name"].(string)
					varID := template["Id"].(string)
					t.Variables[varID] = varName
				}
			}

			if variables, ok := variableSet["Variables"]; ok {
				variables := variables.(map[string]interface{})

				for id, v := range variables {
					varName := t.Variables[id]

					switch varValue := v.(type) {
					case string:
						t.Variables[varName] = varValue
					default:
						continue
					}

					delete(t.Variables, id)
				}
			}
		}
	}

	return nil
}

type Event struct {
	ID       string
	Category string
	Occurred time.Time
}

type client interface {
	FetchMachines() ([]Machine, error)
	FetchMachine(machineID string) (Machine, error)

	FetchProjects() ([]Project, error)
	FetchProject(projectID string) (Project, error)

	FetchTenants() ([]Tenant, error)
	FetchTenant(tenantID string) (Tenant, error)

	FetchEvents(filter map[string]string) ([]Event, error)
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

func (s Service) FetchEvents(filter map[string]string) ([]Event, error) {
	return s.client.FetchEvents(filter)
}
