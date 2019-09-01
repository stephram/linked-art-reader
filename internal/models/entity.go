package models

import (
	"reflect"
	"strings"
)

type Entity interface {
	GetID() string
	GetDOR_ID() string
	GetUUID() string
	GetTMS_ID() string
}

type entityImpl struct {
	ID              string
	UUID            string
	DOR_ID          string
	TMS_ID          string
	AccessionNumber string
	AltIDs          map[string]string
	Type            string
	Title           string
	Name            string
	AltLabels       map[string]string
	Content         interface{}
	WebURL          string
	LocationID      string
	Location        string
	References      map[string]string
	AltReferences   map[string]string
}

func New(object interface{}) Entity {
	entity := &entityImpl{
		ID:            "",
		UUID:          "",
		DOR_ID:        "",
		TMS_ID:        "",
		AltIDs:        map[string]string{},
		Type:          "",
		Title:         "",
		Name:          "",
		AltLabels:     map[string]string{},
		Content:       nil,
		WebURL:        "",
		LocationID:    "",
		Location:      "",
		References:    map[string]string{},
		AltReferences: map[string]string{},
	}

	if _object, ok := object.(*HumanMadeObject); ok {
		entity.ID = _object.ID
		entity.Type = reflect.TypeOf(*_object).String()
		entity.LocationID = _object.CurrentLocation.ID
		entity.Location = _object.CurrentLocation.Label
		getIDs(_object, entity)
		getNamesAndTitles(_object, entity)
		getReferences(_object, entity)
	}
	return entity
}

func getIDs(object interface{}, entity *entityImpl) {
	if _object, ok := object.(*HumanMadeObject); ok {
		for _, identifier := range _object.IdentifiedBy {
			if identifier.Type == "Identifier" {
				if strings.Contains(identifier.Label, "(TMS) ID") {
					(*entity).TMS_ID = identifier.Content
					continue
				}
				if strings.Contains(identifier.Label, "(DOR) ID") {
					(*entity).DOR_ID = identifier.Content
					continue
				}
				if strings.Contains(identifier.Label, "(DOR) UUID") {
					(*entity).UUID = identifier.Content
					continue
				}
				if strings.Contains(identifier.Label, "Accession Number") {
					(*entity).AccessionNumber = identifier.Content
					continue
				}
				(*entity).AltIDs[identifier.Label] = identifier.Content
			}
		}
	}
}

func getReferences(object interface{}, entity *entityImpl) {
	if _object, ok := object.(*HumanMadeObject); ok {
		for _, reference := range _object.ReferredToBy {
			for _, classifier := range reference.ClassifiedAs {
				if classifier.Label != "Brief Text" {
					(*entity).AltReferences[classifier.Label] = reference.Content
				}
			}
		}
	}
}

func getNamesAndTitles(object interface{}, entity *entityImpl) {
	if _object, ok := object.(*HumanMadeObject); ok {
		for _, identifier := range _object.IdentifiedBy {
			if identifier.Type == "Name" {
				if identifier.Label == "Primary Title" {
					(*entity).Title = identifier.Content
					continue
				}
				(*entity).AltLabels[identifier.Label] = identifier.Content
			}
		}
	}
}

func (e *entityImpl) GetID() string {
	return e.ID
}

func (e *entityImpl) GetDOR_ID() string {
	return e.DOR_ID
}

func (e *entityImpl) GetTMS_ID() string {
	return e.TMS_ID
}

// func (e *entityImpl) GetUID() string {
// 	return e.UID
// }

func (e *entityImpl) GetUUID() string {
	return e.UUID
}
