package cdc

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/michaelmosher/monitoring/pkg/metricly"
)

const metriclyMaxResults = 100
const metriclyWorkers = 8

type metriclyStatus struct {
	uaid   string
	sample float64
}

func getMetriclyList(service metriclyClient) ([]metricly.Metric, error) {
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

func getMetricStatus(service metriclyClient, metric metricly.Metric) (metriclyStatus, error) {
	val, err := service.FetchMetricValue(metric)

	return metriclyStatus{
		uaid:   getUAIDFromFQN(metric.FQN),
		sample: val,
	}, err
}

func getUAIDFromFQN(fqn string) string {
	return strings.ToUpper(strings.Split(fqn, ".")[1])
}

func (s *Service) getMetriclySamples() (map[string]float64, error) {
	statuses := make(map[string]float64)
	metricChan := make(chan metricly.Metric)
	statusChan := make(chan metriclyStatus)
	done := make(chan bool)

	var workerWaitGroup sync.WaitGroup

	go func() {
		for status := range statusChan {
			statuses[status.uaid] = status.sample
		}
		done <- true
	}()

	for i := 0; i < metriclyWorkers; i++ {
		workerWaitGroup.Add(1)
		go func() {
			defer workerWaitGroup.Done()

			for metric := range metricChan {
				status, err := getMetricStatus(s.Metricly, metric)
				if err != nil {
					log.Printf("metricly.FetchMetricValue(%s) error: %s", metric.FQN, err)
					continue
				}

				statusChan <- status
			}
		}()
	}

	metrics, err := getMetriclyList(s.Metricly)

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
