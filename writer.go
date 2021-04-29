package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/moxspec/moxspec-metrics-agent/promcli"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/prompb"
)

type writer interface {
	writeTimeSeries(ctx context.Context, mList []metrics) error
	setAuth(user, pass string)
}

func encodeMetricsList(mList []metrics) ([]prompb.TimeSeries, []prompb.MetricMetadata) {
	tsList := make([]prompb.TimeSeries, len(mList))
	mdList := make([]prompb.MetricMetadata, len(mList))

	for _, m := range mList {
		ts := m.flatten()
		tsList = append(tsList, ts)
		mdList = append(mdList, promMetricMetadata())
	}
	return tsList, mdList
}

func promMetricMetadata() prompb.MetricMetadata {
	return prompb.MetricMetadata{
		Type:             prompb.MetricMetadata_GAUGE,
		MetricFamilyName: prompb.MetricMetadata_GAUGE.String(),
	}
}

type promRemoteWriter struct {
	cli promcli.Client
}

func (p promRemoteWriter) writeTimeSeries(ctx context.Context, mList []metrics) error {
	tsList, _ := encodeMetricsList(mList)

	promReq := prompb.WriteRequest{
		Timeseries: tsList,
		//Metadata:   mdList,
	}

	log.Debugf("writing %d metrics", len(tsList))
	status, err := p.cli.Write(ctx, promReq)
	if err != nil {
		// Policy: No Retry
		return errors.Wrap(err, fmt.Sprintf("status code: %d", status))
	}
	log.Infof("remote write: status code = %d", status)
	return nil
}

func (p *promRemoteWriter) setAuth(user, pass string) {
	p.cli.SetAuth(user, pass)
}

func newPromRemoteWriter(endpoint string) (*promRemoteWriter, error) {
	client, err := promcli.NewClient(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "unable to construct client: %v")
	}
	return &promRemoteWriter{
		cli: client,
	}, nil
}

func metricsConsumer(ctx context.Context, ch <-chan []metrics, w writer, interval time.Duration) {
	var buf []metrics
	c := time.Tick(interval)

	outStream := make(chan []metrics)
	go metricsWriter(ctx, outStream, w)

	for {
		select {
		case <-c:
			outStream <- buf
			buf = []metrics{}
		case m := <-ch:
			buf = append(buf, m...)
		case <-ctx.Done():
			return
		}
	}
}

func metricsWriter(ctx context.Context, ch <-chan []metrics, w writer) {
	for {
		select {
		case mList := <-ch:
			log.Infof("writing %d metrics", len(mList))
			err := w.writeTimeSeries(ctx, mList)
			if err != nil {
				log.Error(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

type localWriter struct {
}

func (l localWriter) writeTimeSeries(ctx context.Context, mList []metrics) error {
	tsList, _ := encodeMetricsList(mList)
	log.Debugf("%+v", tsList)
	return nil
}

func (l localWriter) setAuth(user, pass string) {
}

func newLocalWriter(endpoint string) (*localWriter, error) {
	return &localWriter{}, nil
}

type jsonHTTPWriter struct {
	cli      *http.Client
	endpoint string
}

func (j jsonHTTPWriter) writeTimeSeries(ctx context.Context, mList []metrics) error {
	bod, err := json.Marshal(mList)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", j.endpoint, bytes.NewReader(bod))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}
	res, err := j.cli.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", res.StatusCode)
	}
	return nil
}

func (j jsonHTTPWriter) setAuth(user, pass string) {
}

func newJSONHTTPWriter(endpoint string) (*jsonHTTPWriter, error) {
	return &jsonHTTPWriter{
		cli:      &http.Client{},
		endpoint: endpoint,
	}, nil
}
