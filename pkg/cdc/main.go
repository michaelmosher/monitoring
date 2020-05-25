package cdc

import (
	"fmt"
	"sync"

	"github.com/michaelmosher/monitoring/pkg/metricly"
	"github.com/michaelmosher/monitoring/pkg/octopus"
)

type metriclyClient interface {
	FetchMetrics(metricly.MetricQuery) ([]metricly.Metric, error)
	FetchMetricValue(metric metricly.Metric) (float64, error)
}

type octopusClient interface {
	FetchMachines() ([]octopus.Machine, error)
	FetchTenants() ([]octopus.Tenant, error)
	FetchProject(string) (octopus.Project, error)
}

type metricCache struct {
	justOnce     sync.Once
	doneFetching sync.RWMutex
	samples      map[string]float64
}

type Service struct {
	Metricly    metriclyClient
	Octopus     octopusClient
	metricCache metricCache
}

// CheckOfflineNUCs returns a slice of Tenant names.
func (s *Service) CheckOfflineNUCs(projectNames ...string) ([]string, error) {
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

// CheckIdleMachines returns a slice of Tenant names.
func (s *Service) CheckIdleMachines(projectNames ...string) ([]string, error) {
	idle := []string{}

	onlineMachines, err := s.getOnlineMachines()

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

	s.metricCache.justOnce.Do(func() {
		s.metricCache.doneFetching.Lock()
		samples, err := s.getMetriclySamples()

		if err != nil {
			fmt.Printf("error getting Metricly Samples: %s\n", err)
		}
		s.metricCache.samples = samples
		s.metricCache.doneFetching.Unlock()
	})

	for _, nuc := range onlineMachines {
		for id := range nuc.TenantIDs {
			tenant := tenants[id]

			s.metricCache.doneFetching.RLock()
			latency := s.metricCache.samples[tenant.Variables["UAID"]]
			s.metricCache.doneFetching.RUnlock()

			if latency > 600 {
				for _, p := range projects {
					if _, ok := tenant.ProjectIDs[p]; ok == true {
						idle = append(idle, tenant.Name)
					}
				}
			}
		}
	}

	return idle, nil
}

// CheckIdleAOS returns a slice of NUC ID strings.
func CheckIdleAOS(projectNames ...string) ([]string, error) {
	idle := []string{}
	// octopus machines with status != offline,
	//		role contains "sql-server",
	//		metricly replica lag exists,
	//		and metricly lag > 600
	return idle, nil
}
