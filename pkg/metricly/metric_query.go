package metricly

import (
	"time"
)

// MetricQuery is a data type that encapsulates the information required by
// Metricly's `metrics/elasticsearch/metricQuery` API endpoint.
type MetricQuery struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`

	Sort         sortBlock           `json:"sort"`
	ElementFqns  querySpecifierBlock `json:"elementFqns"`
	MetricFqns   querySpecifierBlock `json:"metricFqns"`
	SourceFilter sourceFilterBlock   `json:"sourceFilter"`
}

type sortBlock struct {
	Field   string `json:"field"`
	Order   string `json:"order"`
	Missing string `json:"missing"`
}

type querySpecifierBlock struct {
	Items []querySpecifier `json:"items,omitempty"`
}

type querySpecifier struct {
	Literal  bool   `json:"literal"`
	Contains bool   `json:"contains"`
	Item     string `json:"item"`
}

type sourceFilterBlock struct {
	Includes []string `json:"includes,omitempty"`
	Excludes []string `json:"excludes,omitempty"`
}

// SetStartDate sets the "startDate" field of a MetricQuery. A client could set
// this using standard struct access, but this method prevents them from
// needing to know the proper time format.
func (mq *MetricQuery) SetStartDate(t time.Time) *MetricQuery {
	mq.StartDate = t.Format(time.RFC3339)

	return mq
}

// SetEndDate sets the "endDate" field of a MetricQuery. A client could set
// this using standard struct access, but this method prevents them from
// needing to know the proper time format.
func (mq *MetricQuery) SetEndDate(t time.Time) *MetricQuery {
	mq.EndDate = t.Format(time.RFC3339)

	return mq
}

// AddElement adds a Metricly "element" to a MetricQuery. Multiple invocations
// add additional elements to the query.
func (mq *MetricQuery) AddElement(element string) *MetricQuery {
	newItem := querySpecifier{
		Item:     element,
		Literal:  false,
		Contains: true,
	}

	mq.ElementFqns.Items = append(mq.ElementFqns.Items, newItem)

	return mq
}

// AddMetric adds a Metricly "metric" to a MetricQuery. Multiple invocations
// add additional metrics to the query.
func (mq *MetricQuery) AddMetric(metric string) *MetricQuery {
	newItem := querySpecifier{
		Item:     metric,
		Literal:  false,
		Contains: true,
	}

	mq.MetricFqns.Items = append(mq.MetricFqns.Items, newItem)
	return mq
}

// SetSourceIncludes sets the field values that will be included in the
// response from a MetricQuery.
// TODO: the Metric type should be expanded so any values passed here are
// properly included in the response from `GetMetricsList`
func (mq *MetricQuery) SetSourceIncludes(fields ...string) *MetricQuery {
	mq.SourceFilter = sourceFilterBlock{
		Includes: fields,
	}

	return mq
}

// SetSort sets the ordering of the response from `GetMetricsList` for this
// MetricQuery. The "field" argument is expencted to be a field in the
// response (see `SetSourceIncludes`), but this is not enforced. The "order"
// argument is expected to be either "asc" or "desc", but this is not enforced.
func (mq *MetricQuery) SetSort(field string, order string) *MetricQuery {
	mq.Sort = sortBlock{
		Field:   field,
		Order:   order,
		Missing: "_last",
	}

	return mq
}
