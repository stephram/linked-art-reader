package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"linkedart-reader-golang/internal/models"
	activity_stream "linkedart-reader-golang/internal/readers/activity-stream"
	"linkedart-reader-golang/internal/utils"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"
)

var (
	enProf *bool
	stPage *int
	enPage *int
	asHost *string
	asPath *string
	asSche *string
)

func init() {
	utils.GetLogger()
}

func main() {
	enProf = flag.Bool("prof", false, "enable the pprof package. Listening on port 8080")
	stPage = flag.Int("start", 1, "start at page")
	enPage = flag.Int("end", -1, "stop at page")
	asHost = flag.String("host", "mart.getty.edu", "activity stream host")
	asPath = flag.String("path", "activity-stream", "path to the activity stream")
	asSche = flag.String("scheme", "http", "http(s)")
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
		Scheme: *asSche,
	}

	var orderedCollection *models.OrderedCollection
	var err error

	streamReader := activity_stream.New(endpoint)
	orderedCollection, err = streamReader.GetOrderedCollection("")
	if err != nil {
		log.WithError(err).Errorf("error reading orderedCollection from '%s'", endpoint.String())
		os.Exit(1)
	}
	log.Infof("%T : %+v", orderedCollection, orderedCollection)
	processPageParams(stPage, enPage, orderedCollection)

	for url := orderedCollection.First.ID; url != ""; {
		orderedCollection, err = streamReader.GetOrderedCollection(url)
		if err != nil {
			log.WithError(err).Errorf("error reading orderedCollection '%s'", url)
			time.Sleep(1)
			continue
		}
		for _, orderedItem := range orderedCollection.OrderedItems {
			fmt.Printf("%s: %s: %s\n", orderedItem.ID, orderedItem.Object.Type, orderedItem.Object.ID)

			if len(orderedItem.Object.ID) > 0 {
				object, objErr := streamReader.GetObject(orderedItem.Object.ID)
				if objErr != nil {
					log.WithError(objErr).Errorf("error reading object in %+v", orderedItem)
					time.Sleep(1)
					continue
				}
				fmt.Printf("  %s\n", object.Type)
				for _, identifiedBy := range object.IdentifiedBy {
					switch identifiedBy.RawContent[0] {
					case '"':
						{
							var value string
							if err := json.Unmarshal(identifiedBy.RawContent, &value); err != nil {
								log.WithError(err).Errorf("error unmarshaling '%s', to string", string(identifiedBy.RawContent))
								continue
							}
							fmt.Printf("    %s, %s, ", identifiedBy.Type, value)
						}
					default:
						{
							var value int
							if err := json.Unmarshal(identifiedBy.RawContent, &value); err != nil {
								log.WithError(err).Errorf("error unmarshaling '%s', to int", string(identifiedBy.RawContent))
								continue
							}
							fmt.Printf("    %s, %d, ", identifiedBy.Type, value)
						}
					}
					fmt.Printf("%s\n", identifiedBy.Label)
				}
			}
		}
		url = orderedCollection.Next.ID
	}
}

func processPageParams(stPage *int, enPage *int, orderedCollection *models.OrderedCollection) (int, int) {
	_stPage := orderedCollection.First.GetID()
	_enPage := orderedCollection.Last.GetID()

	if *stPage != _stPage {
		if *stPage > _enPage {
			*stPage = _stPage
		}
	}

	if *enPage != _enPage {
		if *enPage < *stPage {
			*enPage = _enPage
		}
	}

	if *stPage > *enPage {
		*stPage = _stPage
		*enPage = _enPage
	}
	return *stPage, *enPage
}
