package activity_stream

import (
	"encoding/json"
	"io/ioutil"
	"linkedart-reader-golang/internal/models"
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
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	jsonArr, err := ioutil.ReadAll(res.Body)
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
	res, err := http.Get(id)
	if err != nil {
		return nil, err
	}

	jsonArr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var orderedItem models.OrderedItem

	if err := json.Unmarshal(jsonArr, &orderedItem); err != nil {
		return nil, err
	}
	return &orderedItem, err
}

func (r *activityStreamReader) GetObject(id string) (*models.Object, error) {
	res, err := http.Get(id)
	if err != nil {
		return nil, err
	}

	jsonArr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var object models.Object

	if err := json.Unmarshal(jsonArr, &object); err != nil {
		return nil, err
	}
	return &object, err
}

// func (r *activityStreamReader) GetPage(pageID *string) (*models.ActivityStreamPage, error) {
// 	r.endpoint.Path = fmt.Sprintf("activity-stream/%s", *pageID)
func (r *activityStreamReader) GetPage() (*models.OrderedCollectionPage, error) {
	res, err := http.Get(r.endpoint.String())
	if err != nil {
		return nil, err
	}

	var asp models.OrderedCollectionPage

	jbs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jbs, &asp); err != nil {
		return nil, err
	}
	return &asp, nil
}
