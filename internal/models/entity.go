package models

import (
	"encoding/json"
	"linked-art-reader/internal/utils"
	"strings"
)

type Location struct {
	ID       string
	UUID     string
	Location string
}

// type entityImpl struct {
type Entity struct {
	ID              string
	Type            string
	UUID            string
	DOR_ID          string
	TMS_ID          string
	AccessionNumber string
	Title           string
	Content         interface{}
	WebURL          string
	Location        Location
	Labels          map[string]string
	AltIdentifiers  map[string]string
	References      map[string]string
	AltReferences   map[string]string
	Classifiers     map[string]string
}

func New(object TMSObjectIf, jsonb []byte) *Entity {
	entity := &Entity{
		ID:             "",
		Type:           "",
		UUID:           "",
		DOR_ID:         "",
		TMS_ID:         "",
		Title:          "",
		Labels:         map[string]string{},
		Content:        nil,
		WebURL:         "",
		Location:       Location{ID: "", UUID: "", Location: ""},
		AltIdentifiers: map[string]string{},
		References:     map[string]string{},
		AltReferences:  map[string]string{},
		Classifiers:    map[string]string{},
	}
	entity.ID = object.GetID()
	entity.Type = object.GetType()
	getIDs(object, entity)
	getNamesAndTitles(object, entity)
	getReferences(object, entity)
	getClassifiers(object, entity)

	switch entity.Type {
	case "HumanMadeObject":
		{
			var _object HumanMadeObject
			if err := json.Unmarshal(jsonb, &_object); err != nil {
				break
			}
			entity.Location.ID = _object.CurrentLocation.ID
			entity.Location.UUID = utils.GetLastPathComponent(_object.CurrentLocation.ID)
			entity.Location.Location = _object.CurrentLocation.Label
		}
	}
	return entity
}

func getIDs(object TMSObjectIf, entity *Entity) {
	setIDs(object, entity)
}

func setIDs(object TMSObjectIf, entity *Entity) {
	for _, identifier := range object.GetIdentifiedBy() {
		if identifier.Type == "Identifier" {
			if strings.Contains(identifier.Label, "(TMS) ID") {
				(*entity).TMS_ID = identifier.Content
			}
			if strings.Contains(identifier.Label, "(DOR) ID") {
				(*entity).DOR_ID = identifier.Content
			}
			if strings.Contains(identifier.Label, "(DOR) UUID") {
				(*entity).UUID = identifier.Content
			}
			if strings.Contains(identifier.Label, "Accession Number") {
				(*entity).AccessionNumber = identifier.Content
			}
			(*entity).AltIdentifiers[identifier.Label] = identifier.Content
		}
	}
}

func getClassifiers(object TMSObjectIf, entity *Entity) {
	for _, classifier := range object.GetClassifiedAs() {
		(*entity).Classifiers[classifier.Type] = classifier.Label
	}
}

func getReferences(object TMSObjectIf, entity *Entity) {
	for _, reference := range object.GetReferredToBy() {
		for _, classifier := range reference.ClassifiedAs {
			if classifier.Label != "Brief Text" {
				(*entity).AltReferences[classifier.Label] = reference.Content
			}
		}
	}
}

func getNamesAndTitles(object TMSObjectIf, entity *Entity) {
	for _, identifier := range object.GetIdentifiedBy() {
		if identifier.Type == "Name" {
			if identifier.Label == "" {
				(*entity).AltIdentifiers[identifier.Type] = identifier.Content
				continue
			}
			if identifier.Label == "Primary Title" {
				(*entity).Title = identifier.Content
				continue
			}
			(*entity).Labels[identifier.Label] = identifier.Content
			continue
		}
		(*entity).AltIdentifiers[identifier.Label] = identifier.Content
	}
	return
}
