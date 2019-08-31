package linkedart

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"linkedart-reader-golang/internal/models"
	"linkedart-reader-golang/internal/readers"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

type linkedArtReader struct {
	endpoint *url.URL
}

func (r *linkedArtReader) GetNextObject(museumObject *models.MuseumObject) (*models.MuseumObject, error) {
	return &models.MuseumObject{
		ID:          uuid.New().String(),
		Description: "A museum object",
	}, nil
}

func (r *linkedArtReader) GetObject(objectID *string) (*models.MuseumObject, error) {
	if objectID == nil {
		res, err := http.Get(r.endpoint.String())
		if err != nil {
			return nil, err
		}

		var asi models.OrderedCollection

		jbs, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(jbs, &asi); err != nil {
			return nil, err
		}
		return r.GetObject(aws.String(asi.First.ID))
	}

	r.endpoint.Path = fmt.Sprintf("activity-stream/%s", *objectID)

	res, err := http.Get(r.endpoint.String())
	if err != nil {
		return nil, err
	}

	var asp models.ActivityStreamPage
	jbs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jbs, &asp); err != nil {
		return nil, err
	}

	return &models.MuseumObject{
		ID:          uuid.New().String(),
		Description: fmt.Sprintf("%+v", asp),
	}, nil
}

func New(endpoint *url.URL) readers.BaseReader {
	return &linkedArtReader{
		endpoint: endpoint,
	}
}
