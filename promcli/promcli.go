package promcli

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

const (
	useragent = "moxspec-metrics-agent"
)

// Client represents remote writer
type Client struct {
	endPoint string
	timeout  time.Duration
	httpCli  *http.Client
}

// NewClient returns a new client
func NewClient(ep string) (Client, error) {
	return Client{
		endPoint: ep,
		httpCli:  http.DefaultClient,
	}, nil
}

// Write writes prometheus request to the remote write endpoint
func (c Client) Write(ctx context.Context, promReq prompb.WriteRequest) (int, error) {
	data, err := proto.Marshal(&promReq)
	if err != nil {
		return 0, err
	}
	compressed := snappy.Encode(nil, data)

	req, err := http.NewRequestWithContext(ctx, "POST", c.endPoint, bytes.NewBuffer(compressed))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Encoding", "snappy")
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("User-Agent", useragent)
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	// post data
	resp, err := c.httpCli.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode/100 == 2 {
		return resp.StatusCode, nil
	}

	return resp.StatusCode, fmt.Errorf("http error: %d", resp.StatusCode)
}
