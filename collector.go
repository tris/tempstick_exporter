package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

type tempstickMetric struct {
	desc *prometheus.Desc
	value float64
	labels prometheus.Labels
	timestamp time.Time
}

type TempstickCollector struct {
	metrics []*tempstickMetric
}

func (c *TempstickCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m.Desc()
	}
}

func (c *TempstickCollector) Collect(ch chan<- prometheus.Metric) {
	for _, m := range c.metrics {
		ch <- m
	}
}

func (c *tempstickMetric) Desc() *prometheus.Desc {
	return c.desc
}

func (c *tempstickMetric) Write(m *dto.Metric) error {
	m.Label = []*dto.LabelPair{}
	for k, v := range c.labels {
		m.Label = append(m.Label, &dto.LabelPair{
			Name:  proto.String(k),
			Value: proto.String(v),
		})
	}
	m.Gauge = &dto.Gauge{Value: &c.value}
	m.TimestampMs = proto.Int64(c.timestamp.UnixNano() / int64(time.Millisecond))
	return nil
}

func newTempstickMetric(name string, help string, labels prometheus.Labels, value float64, timestamp time.Time) *tempstickMetric {
	return &tempstickMetric{
		desc: prometheus.NewDesc(name, help, nil, labels),
		value: value,
		labels: labels,
		timestamp: timestamp,
	}
}
