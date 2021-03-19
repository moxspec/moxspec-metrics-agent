package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/m3db/prometheus_remote_client_golang/promremote"
	"github.com/moxspec/moxspec/util"
)

const (
	jobName = "moxspec"
	prefix  = "mox"
)

type metrics struct {
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
	Timestamp time.Time         `json:"timestamp"`
	Value     float64           `json:"value"`
}

func (m metrics) flatten() promremote.TimeSeries {
	var lbl []promremote.Label
	for k, v := range m.Labels {
		lbl = append(lbl, promremote.Label{Name: k, Value: v})
	}

	hostname, err := os.Hostname()
	if err != nil {
		// TODO: try to find a better way
		hostname = "localhost"
	}

	lbl = append(lbl, []promremote.Label{
		{Name: "__name__", Value: m.Name},
		{Name: "instance", Value: hostname},
		{Name: "job", Value: jobName},
	}...)

	return promremote.TimeSeries{
		Labels: lbl,
		Datapoint: promremote.Datapoint{
			Timestamp: m.Timestamp,
			Value:     m.Value,
		},
	}
}

func newMetrics(name string, labels map[string]string, val interface{}) (metrics, error) {
	if name == "" {
		return metrics{}, fmt.Errorf("empty name given")
	}
	if len(labels) == 0 {
		return metrics{}, fmt.Errorf("empty labels given")
	}

	v, err := util.CastToFloat64(val)
	if err != nil {
		return metrics{}, err
	}

	return metrics{
		Name:      getMetricsNameWithPrefix(name),
		Labels:    labels,
		Timestamp: time.Now(),
		Value:     v,
	}, nil
}

func getMetricsNameWithPrefix(name string) string {
	return prefix + "_" + name
}

type metricsReporter func() ([]metrics, error)

type profile struct {
	reporter metricsReporter
	interval time.Duration
}

func metricsProducer(ctx context.Context, ch chan<- []metrics, p profile) {
	for {
		m, err := p.reporter()
		if err != nil {
			log.Error(err.Error())
			return
		}

		select {
		case ch <- m:
		case <-ctx.Done():
			return
		}
		time.Sleep(p.interval)
	}
}

func createMetricsList(labels map[string]string, points map[string]interface{}) ([]metrics, error) {
	if len(labels) == 0 {
		return nil, fmt.Errorf("empty labels given")
	}

	if len(points) == 0 {
		return nil, fmt.Errorf("empty points given")
	}

	var mList []metrics
	for k, v := range points {
		m, err := newMetrics(k, labels, v)
		if err != nil {
			return nil, err
		}
		mList = append(mList, m)
	}
	return mList, nil
}

func copyLabels(src map[string]string) map[string]string {
	if src == nil {
		return nil
	}
	dst := make(map[string]string)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
