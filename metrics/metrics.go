package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricValue struct {
	Fields map[string]string
	Values []float64
}

type Metric struct {
	Namespace string
	Name string
	Help string
	Fields []string

	Values []MetricValue
	PrometheusGauge *prometheus.GaugeVec
}

var metrics = []Metric{}

func init() {
	metrics = append(metrics, sun2000Metrics...)
	metrics = append(metrics, luna2000Metrics...)
	metrics = append(metrics, powerMeterMetrics...)

	for index, metric := range metrics {
		metric.Values = []MetricValue{}
		metric.PrometheusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: metric.Namespace,
				Name: metric.Name,
				Help: metric.Help,
			},
			metric.Fields,
		)

		metrics[index] = metric
	}
}

func AddMetric(namespace string, name string, help string, fields []string) *Metric {
	metrics = append(metrics, Metric{
		Namespace: namespace,
		Name: name,
		Help: help,
		Fields: fields,
		Values: []MetricValue{},
		PrometheusGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: name,
			Help: help,
		}, fields),
	})

	return &metrics[len(metrics) - 1]
}

func GetMetric(namespace string, name string) *Metric {
	for _, metric := range metrics {
		if metric.Namespace == namespace && metric.Name == name {
			return &metric
		}
	}

	return nil
}

func GetMetricValueAverageLastMin(namespace string, name string, labels map[string]string, minutes uint) float64 {
	var newMetric *Metric
	for _, metric := range metrics {
		if metric.Namespace == namespace && metric.Name == name {
			newMetric = &metric
			break
		}
	}

	if newMetric == nil {
		return 0
	}

	var matches *MetricValue
	for _, metricValue := range newMetric.Values {
		match := true
		for key, value := range labels {
			if metricValue.Fields[key] != value {
				match = false
				break
			}
		}

		if !match {
			continue
		}

		matches = &metricValue
	}

	if matches == nil {
		return 0
	}

	var total float64
	for _, value := range matches.Values {
		total += value
	}

	return total / float64(len(matches.Values))
}

func SetMetricValue(namespace string, name string, labels map[string]string, value float64) {
	var newMetric *Metric
	var metricIndex int
	for i, metric := range metrics {
		if metric.Namespace == namespace && metric.Name == name {
			newMetric = &metric
			metricIndex = i
			break
		}
	}

	if newMetric == nil {
		return
	}

	var matches *MetricValue
	for _, metricValue := range newMetric.Values {
		match := true
		for key, value := range labels {
			if metricValue.Fields[key] != value {
				match = false
				break
			}
		}

		if !match {
			continue
		}

		matches = &metricValue
	}

	if matches != nil {
		matches.Values = append(matches.Values, value)
		if len(matches.Values) > 5760 { // 1 days worth of data if we have a value every 15 seconds
			matches.Values = matches.Values[1:]
		}
	} else {
		newMetric.Values = append(newMetric.Values, MetricValue{
			Fields: labels,
			Values: []float64{value},
		})
	}

	metrics[metricIndex] = *newMetric
	newMetric.PrometheusGauge.With(labels).Set(value)
}
