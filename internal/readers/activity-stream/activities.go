package activity_stream

import (
	"encoding/json"
	"io/ioutil"
	"linked-art-reader/internal/models"
	"net/http"
	"net/url"
)

type ActivityStreamReader interface {
}

type activityStreamReader struct {
	endpoint *url.URL
}

func New(endpoint *url.URL) *activityStreamReader {
	return &activityStreamReader{
		endpoint: endpoint,
	}
}

func (r *activityStreamReader) GetOrderedCollection(id string) (*models.OrderedCollection, error) {
	url := r.endpoint.String()
	if id != "" {
		url = id
	}
	jsonArr, err := getBytes(url)
	if err != nil {
		return nil, err
	}

	var orderedCollection models.OrderedCollection

	if err := json.Unmarshal(jsonArr, &orderedCollection); err != nil {
		return nil, err
	}
	return &orderedCollection, err
}

func (r *activityStreamReader) GetOrderedCollectionItem(id string) (*models.OrderedItem, error) {
	jsonArr, err := getBytes(id)
	if err != nil {
		return nil, err
	}

	var orderedItem models.OrderedItem

	if err := json.Unmarshal(jsonArr, &orderedItem); err != nil {
		return nil, err
	}
	return &orderedItem, err
}

func (r *activityStreamReader) GetTMSObject(id string) (*models.TMSObject, []byte, error) {
	jsonArr, err := getBytes(id)
	if err != nil {
		return nil, jsonArr, err
	}

	var object models.TMSObject

	if err := json.Unmarshal(jsonArr, &object); err != nil {
		return nil, jsonArr, err
	}
	return &object, jsonArr, err
}

func (r *activityStreamReader) GetTypedObject(id string) (*models.Object, []byte, error) {
	jsonArr, err := getBytes(id)
	if err != nil {
		return nil, jsonArr, err
	}

	var object models.Object

	if err := json.Unmarshal(jsonArr, &object); err != nil {
		return nil, jsonArr, err
	}
	return &object, jsonArr, err
}

func (r *activityStreamReader) GetObject(id string) (*models.Object, error) {
	jsonArr, err := getBytes(id)
	if err != nil {
		return nil, err
	}

	var object models.Object

	if err := json.Unmarshal(jsonArr, &object); err != nil {
		return nil, err
	}
	return &object, err
}

func getBytes(id string) ([]byte, error) {
	res, err := http.Get(id)
	if err != nil {
		return nil, err
	}

	jsonArr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return jsonArr, err
}
