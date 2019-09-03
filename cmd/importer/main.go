package main

import (
	"flag"
	"linked-art-reader/internal/utils"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

var (
	enProf *bool
	asHost *string
	asPath *string
	asProt *string
)

func init() {
	utils.GetLogger()
}

func main() {
	enProf = flag.Bool("prof", false, "enable the pprof package. Listening on port 8080")
	asHost = flag.String("host", "mart.getty.edu", "activity stream host")
	asPath = flag.String("path", "activity-stream", "path to the activity stream")
	asProt = flag.String("scheme", "http", "http(s)")
	flag.Parse()

	if *enProf {
		go func() {
			log.Infof("enabled profiler on port 8080")
			_ = http.ListenAndServe("localhost:8080", nil)
		}()
	}

	endpoint := &url.URL{
		Host:   *asHost,
		Path:   *asPath,
		Scheme: *asProt,
	}
	log.Infof("dbhost: %s\n", endpoint.String())
}
