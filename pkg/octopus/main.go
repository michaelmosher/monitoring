package octopus

type Machine struct {
	Name      string
	Status    string
	Roles     []string
	TenantIds []string
}

type Tenant struct {
	ID         string
	Name       string
	ProjectIDs []string
}

type client interface {
	FetchMachines() ([]Machine, error)
	FetchMachine(machineID string) (Machine, error)

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

func (s Service) FetchTenants() ([]Tenant, error) {
	return s.client.FetchTenants()
}

func (s Service) FetchTenant(tenantID string) (Tenant, error) {
	return s.client.FetchTenant(tenantID)
}
