package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/moxspec/moxspec/loglet"
)

var (
	log *loglet.Logger
)

const (
	bufMargin = 2.0
	bufMax    = 65535
)

func init() {
	log = loglet.NewLogger("main")
}

func main() {
	var (
		output    string
		endpoint  string
		retIntSec int
		subIntSec int
		authUser  string
		authPass  string
		debug     bool
	)

	flag.StringVar(&output, "o", "promRemote", "output selection (promRemote, jsonHttp, stdout)")
	flag.StringVar(&endpoint, "e", "http://localhost:3030/remote/write", "a remote-write endpoint")
	flag.IntVar(&retIntSec, "r", 5, "default retrieval interval (sec)")
	flag.IntVar(&subIntSec, "s", 10, "default submission interval (sec)")
	flag.BoolVar(&debug, "d", false, "enable debug logging")
	flag.Parse()

	authUser = os.Getenv("MOXSPEC_AUTH_USER")
	authPass = os.Getenv("MOXSPEC_AUTH_PASS")

	var (
		w   writer
		err error
	)
	switch output {
	case "promRemote":
		w, err = newPromRemoteWriter(endpoint)
	case "jsonHttp":
		w, err = newJSONHTTPWriter(endpoint)
	case "stdout":
		w, err = newLocalWriter(endpoint)
	default:
		w, err = newLocalWriter(endpoint)
	}
	if err != nil {
		log.Fatal(err)
	}

	if authUser != "" && authPass != "" {
		log.Info("auth info set")
		w.setAuth(authUser, authPass)
	}

	if retIntSec < 1 {
		log.Fatal("invalid default retrieval interval given")
	}

	if subIntSec < 1 {
		log.Fatal("invalid default submission interval given")
	}

	loglet.SetLevel(loglet.INFO)
	if debug {
		loglet.SetLevel(loglet.DEBUG)
	}

	// dry run to find invalid profiles and determine the buffer size
	defRetInterval := time.Second * time.Duration(retIntSec)
	profs, mps, err := dryrun(scanDevices(defRetInterval))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("expected mps: %.1f", mps)

	defSubInterval := time.Second * time.Duration(subIntSec)
	log.Infof("report interval: %ds", defSubInterval/time.Second)

	bufSize, err := calcBufferSize(defSubInterval, mps)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("internal buffer: %d", bufSize)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricsStream := make(chan []metrics, bufSize)
	for _, p := range profs { // TODO: add a system information producer
		go metricsProducer(ctx, metricsStream, p)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsConsumer(ctx, metricsStream, w, defSubInterval)
	}()
	wg.Wait()
}

func dryrun(in []profile) ([]profile, float64, error) {
	if len(in) < 1 {
		return nil, 0, fmt.Errorf("no profile given")
	}

	var profs []profile
	var mps float64 // metrics per second
	for _, p := range in {
		ms, err := p.reporter()
		if err != nil {
			log.Warnf("dryrun error: %s", err)
			continue
		}
		if p.interval <= 0 {
			log.Warnf("zero interval found")
			continue
		}

		mps += float64(len(ms)) / float64(p.interval/time.Second)
		profs = append(profs, p)
	}

	if len(profs) < 1 {
		return nil, 0, fmt.Errorf("no valid profile found")
	}

	if mps == 0 {
		return nil, 0, fmt.Errorf("no valid mps calcurated")
	}

	return profs, mps, nil
}

func calcBufferSize(repInterval time.Duration, mps float64) (int, error) {
	bufSize := int(float64(repInterval)/float64(time.Second)*math.Ceil(mps)) * bufMargin
	if bufSize < 1 {
		return 0, fmt.Errorf("invalid buffer size")
	}
	if bufSize > bufMax {
		log.Warnf("buffer size is truncated from %d to %d", bufSize, bufMax)
		bufSize = bufMax
	}
	return bufSize, nil
}
