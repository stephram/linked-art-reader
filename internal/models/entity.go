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
	AltIDs          map[string]string
	Title           string
	AltLabels       map[string]string
	Content         interface{}
	WebURL          string
	Location        Location
	References      map[string]string
	AltReferences   map[string]string
}

func New(object TMSObjectIf, jsonb []byte) *Entity {
	entity := &Entity{
		ID:            "",
		Type:          "",
		UUID:          "",
		DOR_ID:        "",
		TMS_ID:        "",
		AltIDs:        map[string]string{},
		Title:         "",
		AltLabels:     map[string]string{},
		Content:       nil,
		WebURL:        "",
		Location:      Location{ID: "", UUID: "", Location: ""},
		References:    map[string]string{},
		AltReferences: map[string]string{},
	}
	entity.ID = object.GetID()
	entity.Type = object.GetType()
	getIDs(object, entity)
	getNamesAndTitles(object, entity)
	getReferences(object, entity)

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
			(*entity).AltIDs[identifier.Label] = identifier.Content
		}
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
				(*entity).AltIDs[identifier.Type] = identifier.Content
				continue
			}
			if identifier.Label == "Primary Title" {
				(*entity).Title = identifier.Content
				continue
			}
			(*entity).AltLabels[identifier.Label] = identifier.Content
			continue
		}
		(*entity).AltIDs[identifier.Label] = identifier.Content
	}
	return
}
