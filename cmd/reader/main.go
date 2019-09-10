package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"linked-art-reader/internal/models"
	activity_stream "linked-art-reader/internal/readers/activity-stream"
	"linked-art-reader/internal/repo"
	"linked-art-reader/internal/utils"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	pretty *bool
)

func init() {
	utils.GetLogger()
}

func main() {
	enProf = flag.Bool("prof", false, "enable the pprof package. Listening on port 8080")
	stPage = flag.Int("start", 1, "start at page")
	enPage = flag.Int("end", -1, "stop at page")
	pretty = flag.Bool("pretty", false, "pretty print JSON output")
	asHost = flag.String("host", "", "activity stream host")
	asPath = flag.String("path", "activity-stream", "path to the activity stream")
	asSche = flag.String("scheme", "http", "http(s)")
	flag.Parse()

	larDB := repo.New(nil)
	if larDB == nil {
		log.Errorf("failed to open repository")
		os.Exit(100)
	}

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
	log.Infof("pages %d - %d", *stPage, *enPage)

	baseURL := fmt.Sprintf("%s/page/%d", endpoint.String(), *stPage)

	for url := baseURL; url != ""; {
		orderedCollection, err = streamReader.GetOrderedCollection(url)
		if err != nil {
			log.WithError(err).Errorf("error reading orderedCollection '%s'", url)
			time.Sleep(1)
			continue
		}
		for _, orderedItem := range orderedCollection.OrderedItems {
			if len(orderedItem.Object.ID) > 0 {
				_object, _jsonb, err := streamReader.GetTMSObject(orderedItem.Object.ID)
				if err != nil {
					log.WithError(err).Errorf("error reading '%+v'", orderedItem.Object)
					continue
				}
				resolveIdentifiedBy(&_object.IdentifiedBy)
				resolveClassifiedAs(&_object.ClassifiedAs)
				resolveReferredToBy(&_object.ReferredToBy)

				entity := models.New(_object, _jsonb)
				_, newErr := larDB.StoreEntity(entity)
				if newErr != nil {
					log.WithError(err).Errorf("error storing entity")
				}
				if *pretty {
					fmt.Printf("%s\n", utils.ConvertToPrettyJSON(entity))
					continue
				}
				fmt.Printf("%s\n", utils.ConvertToJSON(entity))
			}
		}
		nextPage := orderedCollection.Next.GetID()
		if nextPage < 0 || nextPage >= *enPage {
			url = ""
			continue
		}
		url = orderedCollection.Next.ID
	}
}

func resolveIdentifiedBy(identifiedByArray *[]models.Identifier) {
	for i, _ := range *identifiedByArray {
		_identifiedBy := &(*identifiedByArray)[i]

		switch _identifiedBy.RawContent[0] {
		case '"':
			{
				var value string
				if err := json.Unmarshal(_identifiedBy.RawContent, &value); err != nil {
					log.WithError(err).Errorf("error unmarshaling '%s', to string", string(_identifiedBy.RawContent))
					continue
				}
				_identifiedBy.Content = value
			}
		default:
			{
				var value int
				if err := json.Unmarshal(_identifiedBy.RawContent, &value); err != nil {
					log.WithError(err).Errorf("error unmarshaling '%s', to int", string(_identifiedBy.RawContent))
					continue
				}
				_identifiedBy.Content = strconv.Itoa(value)
			}
		}
	}
}

func resolveReferredToBy(referredToByArray *[]models.LinguisticObject) {
	for i, _ := range *referredToByArray {
		_referredToBy := &(*referredToByArray)[i]
		resolveClassifiedAs(&_referredToBy.ClassifiedAs)
	}
}

func resolveClassifiedAs(classifiedAsArray *[]models.Type) {
	// for i, _ := range *classifiedAsArray {
	// 	// _classifiedAs := &(*classifiedAsArray)[i]
	// }
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
