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

	baseURL := fmt.Sprintf("%s/page/%d", endpoint.String(), *stPage)

	for url := baseURL; url != ""; {
		orderedCollection, err = streamReader.GetOrderedCollection(url)
		if err != nil {
			log.WithError(err).Errorf("error reading orderedCollection '%s'", url)
			time.Sleep(1)
			continue
		}
		for _, orderedItem := range orderedCollection.OrderedItems {
			// fmt.Printf("%s: %s: %s\n", orderedItem.ID, orderedItem.Object.Type, orderedItem.Object.ID)

			if len(orderedItem.Object.ID) > 0 {
				_object, _jsonb, err := streamReader.GetTypedObject(orderedItem.Object.ID)
				if err != nil {
					log.WithError(err).Errorf("error reading '%+v'", orderedItem.Object)
					continue
				}

				switch _object.Type {
				case "Person":
					{
						object, objErr := streamReader.HydratePerson(_jsonb)
						if objErr != nil {
							log.WithError(objErr).Errorf("error reading object in %+v", orderedItem)
							time.Sleep(1)
							continue
						}
						fmt.Printf("%s : %s\n", object.Type, object.ID)
						resolveIdentifiedBy(&object.IdentifiedBy)
					}
				case "Group":
					{
						object, objErr := streamReader.HydrateGroup(_jsonb)
						if objErr != nil {
							log.WithError(objErr).Errorf("error reading object in %+v", orderedItem)
							time.Sleep(1)
							continue
						}
						fmt.Printf("%s : %s\n", object.Type, object.ID)
						resolveIdentifiedBy(&object.IdentifiedBy)
					}
				case "HumanMadeObject":
					{
						object, objErr := streamReader.HydrateHumanMadeObject(_jsonb)
						if objErr != nil {
							log.WithError(objErr).Errorf("error reading object in %+v", orderedItem)
							time.Sleep(1)
							continue
						}
						fmt.Printf("%s : %s\n", object.Type, object.ID)
						resolveIdentifiedBy(&object.IdentifiedBy)
						fmt.Print("\n")
					}
				default:
					{
						fmt.Printf("%s : %s\n", _object.Type, _object.ID)
						resolveIdentifiedBy(&_object.IdentifiedBy)
					}
				}
			}
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
				fmt.Printf("    %s, %s, ", _identifiedBy.Type, value)
			}
		default:
			{
				var value int
				if err := json.Unmarshal(_identifiedBy.RawContent, &value); err != nil {
					log.WithError(err).Errorf("error unmarshaling '%s', to int", string(_identifiedBy.RawContent))
					continue
				}
				_identifiedBy.Content = strconv.Itoa(value)
				fmt.Printf("    %s, %d, ", _identifiedBy.Type, value)
			}
		}
		fmt.Printf("%s\n", _identifiedBy.Label)
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
