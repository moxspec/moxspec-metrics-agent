package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/m3db/prometheus_remote_client_golang/promremote"
	"github.com/pkg/errors"
)

type writer interface {
	writeTimeSeries(ctx context.Context, mList []metrics) error
}

func encodeMetricsList(mList []metrics) []promremote.TimeSeries {
	var tsList []promremote.TimeSeries
	for _, m := range mList {
		tsList = append(tsList, m.flatten())
	}
	return tsList
}

type promRemoteWriter struct {
	cli promremote.Client
}

func (p promRemoteWriter) writeTimeSeries(ctx context.Context, mList []metrics) error {
	tsList := encodeMetricsList(mList)
	headers := make(map[string]string)

	log.Debugf("writing %d metrics", len(tsList))
	result, err := p.cli.WriteTimeSeries(ctx, tsList, promremote.WriteOptions{Headers: headers})
	if err != nil {
		// Policy: No Retry
		return errors.Wrap(err, fmt.Sprintf("status code: %d", result.StatusCode))
	}
	log.Infof("remote write: status code = %d", result.StatusCode)
	return nil
}

func newPromRemoteWriter(endpoint string) (*promRemoteWriter, error) {
	cfg := promremote.NewConfig(
		promremote.WriteURLOption(endpoint),
		promremote.HTTPClientTimeoutOption(30*time.Second),
	)

	client, err := promremote.NewClient(cfg)
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
	tsList := encodeMetricsList(mList)
	log.Debugf("%+v", tsList)
	return nil
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

func newJSONHTTPWriter(endpoint string) (*jsonHTTPWriter, error) {
	return &jsonHTTPWriter{
		cli:      &http.Client{},
		endpoint: endpoint,
	}, nil
}
