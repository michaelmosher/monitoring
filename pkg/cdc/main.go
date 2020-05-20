package cdc

import (
	"log"
	"net/http"
	"time"

	"github.com/michaelmosher/monitoring/pkg/octopus"
)

const nucOctopusRole = "side-server-appliance"
const vmOctopusRole = "linux-server"

type Status struct {
	Name        string
	Online      bool
	Replicating bool
	AOS         bool
}

type octopusClient interface {
	FetchMachines() ([]octopus.Machine, error)
	FetchTenants() ([]octopus.Tenant, error)
	FetchProject(string) (octopus.Project, error)
}

type Service struct {
	Octo octopusClient
}

// CheckOfflineNUCs returns a slice of Tenant names.
// To get this data, first it must determine which Tenants are associated with
// the configured Projects. Then, determine which Machines are associated with
// those Tenants. Machines with the correct role with status == offline are
// returned.
func (s Service) CheckOfflineNUCs(projectNames ...string) ([]string, error) {
	offline := []string{}

	offlineNUCs, err := s.getOfflineNUCs()

	if err != nil {
		return nil, err
	}

	tenants, err := s.getOctopusTenants()

	if err != nil {
		return nil, err
	}

	projects, err := s.getOctopusProjectIDs(projectNames...)

	if err != nil {
		return nil, err
	}

	for _, nuc := range offlineNUCs {
		for id := range nuc.TenantIDs {
			tenant := tenants[id]
			for _, p := range projects {
				if _, ok := tenant.ProjectIDs[p]; ok == true {
					offline = append(offline, tenant.Name)
				}
			}
		}
	}

	return offline, nil
}

// CheckIdleNUCs returns a slice of NUC ID strings.
// To get this data, first it must determine which Metrics Metricly has
// that match a particular query. Then, it must get a Sample of each Metric.
//
// CheckIdleNUCs must also determine which Tenants are associated with the
// configured Projects. Then, determine which Machines are associated with
// those Tenants. Machines with the correct role, with status != offline, that
// match high Metricly Samples are returned.
func CheckIdleNUCs() (map[string]Status, error) {
	nucs := make(map[string]Status)

	// octopus machines with status != offline,
	// 		role contains "side-server-appliances" or "linux-server",
	//		and metricly replica lag > 600

	return nucs, nil
}

// CheckIdleAOS returns a slice of NUC ID strings.
func CheckIdleAOS() (map[string]Status, error) {
	nucs := make(map[string]Status)
	// octopus machines with status != offline,
	//		role contains "sql-server",
	//		metricly replica lag exists,
	//		and metricly lag > 600
	return nucs, nil
}

func CheckStatus(metriclyUser string, metriclyPassword string, octopusURL string, octopusSpace string, octopusAPIKey string) (map[string]Status, error) {
	statuses := make(map[string]Status)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	octopusService := octopusHTTPService(httpClient, octopusURL, octopusSpace, octopusAPIKey)
	oStatuses, err := getOctopusStatuses(octopusService)

	if err != nil {
		log.Fatalf("getOctopusStatuses error: %s", err)
	}

	metriclyService := metriclyHTTPService(httpClient, metriclyUser, metriclyPassword)
	mStatuses, err := getMetriclyStatuses(metriclyService)

	if err != nil {
		log.Fatalf("getMetriclyStatuses error: %s", err)
	}

	for uaid, status := range oStatuses {
		s, ok := statuses[uaid]

		if ok != true {
			s = Status{}
		}

		s.Online = status.online
		s.Name = status.name
		statuses[uaid] = s
	}

	for uaid, replicating := range mStatuses {
		s, ok := statuses[uaid]

		if ok != true {
			// things that exist in Metricly but not Octopus are AOS.
			s = Status{AOS: true}
		}

		s.Replicating = replicating
		statuses[uaid] = s
	}

	return statuses, err
}
