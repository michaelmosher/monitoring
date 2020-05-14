package cdc

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/michaelmosher/monitoring/pkg/metricly"
	metricly_http "github.com/michaelmosher/monitoring/pkg/metricly/http"
)

const metriclyMaxResults = 60
const metriclyWorkers = 8

type metriclyStatus struct {
	uaid        string
	replicating bool
}

func metriclyHTTPService(httpClient *http.Client, username string, password string) metricly.Service {
	return metricly.New(
		metricly_http.Service{
			HTTPClient: httpClient,
			Username:   username,
			Password:   password,
		},
	)
}

func getMetriclyList(service metricly.Service) ([]metricly.Metric, error) {
	metricsQuery := new(metricly.MetricQuery).
		SetStartDate(time.Now().Add(-1*time.Hour)).
		SetEndDate(time.Now()).
		AddElement("prod-hvr-hub-001").
		AddMetric("hvr_latency").
		SetSourceIncludes("fqn", "id", "element").
		SetSort("fqn", "asc")

	metricsQuery.PageSize = metriclyMaxResults

	return service.FetchMetrics(*metricsQuery)
}

func getMetricStatus(service metricly.Service, metric metricly.Metric) (metriclyStatus, error) {
	val, err := service.FetchMetricValue(metric)

	// log.Printf("%s latency: %d\n", metric.FQN, int(val))

	return metriclyStatus{
		uaid:        getUAIDFromFQN(metric.FQN),
		replicating: int(val) < 600,
	}, err
}

func getUAIDFromFQN(fqn string) string {
	return strings.ToUpper(strings.Split(fqn, ".")[1])
}

func getMetriclyStatuses(service metricly.Service) (map[string]bool, error) {
	statuses := make(map[string]bool)
	metricChan := make(chan metricly.Metric)
	statusChan := make(chan metriclyStatus)
	done := make(chan bool)

	var workerWaitGroup sync.WaitGroup

	go func() {
		for status := range statusChan {
			statuses[status.uaid] = status.replicating
		}
		done <- true
	}()

	for i := 0; i < metriclyWorkers; i++ {
		workerWaitGroup.Add(1)
		go func() {
			defer workerWaitGroup.Done()

			for metric := range metricChan {
				status, err := getMetricStatus(service, metric)
				if err != nil {
					log.Printf("metricly.FetchMetricValue(%s) error: %s", metric.FQN, err)
					continue
				}

				statusChan <- status
			}
		}()
	}

	metrics, err := getMetriclyList(service)

	if err != nil {
		return nil, fmt.Errorf("metricly.FetchMetrics error: %s", err)
	}

	for _, metric := range metrics {
		metricChan <- metric
	}

	close(metricChan)

	workerWaitGroup.Wait()
	close(statusChan)

	<-done
	return statuses, nil
}
