package cdc

import (
	"log"
	"net/http"
	"time"
)

type Status struct {
	Name        string
	Online      bool
	Replicating bool
	AOS         bool
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
