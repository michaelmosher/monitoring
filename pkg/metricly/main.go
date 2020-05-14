package metricly

// Metric is a structure that defines a "metric"; used to look up a metric "result".
type Metric struct {
	ID        string
	ElementID string
	FQN       string
}

type client interface {
	FetchMetrics(MetricQuery) ([]Metric, error)
	FetchMetricValue(Metric) (float64, error)
}

type Service struct {
	client client
}

func New(client client) Service {
	return Service{client: client}
}

func (s Service) FetchMetrics(query MetricQuery) ([]Metric, error) {
	return s.client.FetchMetrics(query)
}

func (s Service) FetchMetricValue(metric Metric) (float64, error) {
	return s.client.FetchMetricValue(metric)
}
