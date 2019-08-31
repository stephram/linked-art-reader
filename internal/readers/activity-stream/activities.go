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

func (r *activityStreamReader) GetPerson(id string) (*models.Person, error) {
	jsonArr, err := getBytes(id)
	if err != nil {
		return nil, err
	}
	return r.HydratePerson(jsonArr)
}

func (r *activityStreamReader) HydratePerson(jsonb []byte) (*models.Person, error) {
	var person models.Person

	if err := json.Unmarshal(jsonb, &person); err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *activityStreamReader) GetGroup(id string) (*models.Group, error) {
	jsonArr, err := getBytes(id)
	if err != nil {
		return nil, err
	}
	return r.HydrateGroup(jsonArr)
}

func (r *activityStreamReader) HydrateGroup(jsonb []byte) (*models.Group, error) {
	var group models.Group

	if err := json.Unmarshal(jsonb, &group); err != nil {
		return nil, err
	}
	return &group, nil

}

func (r *activityStreamReader) GetHumanMadeObject(id string) (*models.HumanMadeObject, error) {
	jsonArr, err := getBytes(id)
	if err != nil {
		return nil, err
	}
	return r.HydrateHumanMadeObject(jsonArr)
}

func (r *activityStreamReader) HydrateHumanMadeObject(jsonb []byte) (*models.HumanMadeObject, error) {
	var humanMadeObject models.HumanMadeObject

	if err := json.Unmarshal(jsonb, &humanMadeObject); err != nil {
		return nil, err
	}
	return &humanMadeObject, nil
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
