package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type OrderedCollection struct {
	Summary      string                `json:"summary"`
	Type         string                `json:"type"`
	ID           string                `json:"id"`
	StartIndex   int                   `json:"startIndex"`
	TotalItems   int                   `json:"totalItems"`
	TotalPages   int                   `json:"totalPages"`
	MaxPerPage   int                   `json:"maxPerPage"`
	First        OrderedCollectionPage `json:"first"`
	Last         OrderedCollectionPage `json:"last"`
	Next         OrderedCollectionPage `json:"next,omitempty"`
	OrderedItems []OrderedItem         `json:"orderedItems,omitempty"`
}

type OrderedCollectionPage struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (p *OrderedCollectionPage) GetID() int {
	ss := strings.Split(p.ID, "/")
	if len(ss) == 0 {
		panic(fmt.Sprintf("invalid ID '%s'", p.ID))
	}
	pg := ss[len(ss)-1]

	val, err := strconv.Atoi(pg)
	if err != nil {
		return -1
	}
	return val
}

type OrderedItem struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Actor     string    `json:"actor"`
	Object    Object    `json:"object"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Published time.Time `json:"published"`
}

type Object struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type LinguisticObject struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Label        string `json:"_label,omitempty"`
	Content      string `json:"content,omitempty"`
	Format       string `json:"format,omitempty"`
	ClassifiedAs []Type `json:"classified_as,omitempty"`
}

type InformationObject struct {
	ID     string   `json:"id"`
	Type   string   `json:"type"`
	Format []string `json:"format,omitempty"`
}

type Place struct {
	ID    string `json:"id"`
	Type  string `json:"type,omitempty"`
	Label string `json:"_label,omitempty"`
}

type Type struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Label string `json:"_label,omitempty"`
}

type Identifier struct {
	ID         string          `json:"id"`
	Type       string          `json:"type,omitempty"`
	Label      string          `json:"_label,omitempty"`
	RawContent json.RawMessage `json:"content,omitempty"`
	Content    string          `json:"-"`
}

type TMSObjectIf interface {
	GetID() string
	GetType() string
	GetIdentifiedBy() []Identifier
	GetClassifiedAs() []Type
	GetReferredToBy() []LinguisticObject
}

type TMSObject struct {
	ID           string             `json:"id"`
	Type         string             `json:"type"`
	IdentifiedBy []Identifier       `json:"identified_by,omitempty"`
	ClassifiedAs []Type             `json:"classified_as,omitempty"`
	ReferredToBy []LinguisticObject `json:"referred_to_by,omitempty"`
}

func (t *TMSObject) GetID() string {
	return t.ID
}

func (t *TMSObject) GetType() string {
	return t.Type
}

func (t *TMSObject) GetIdentifiedBy() []Identifier {
	return t.IdentifiedBy
}

func (t *TMSObject) GetClassifiedAs() []Type {
	return t.ClassifiedAs
}

func (t *TMSObject) GetReferredToBy() []LinguisticObject {
	return t.ReferredToBy
}

type Person struct {
	TMSObject
}

type Group struct {
	TMSObject
}

type HumanMadeObject struct {
	TMSObject
	Label           string          `json:"_label,omitempty"`
	CurrentKeeper   []Object        `json:"current_keeper,omitempty"`
	CurrentOwner    []Object        `json:"current_owner,omitempty"`
	CurrentLocation Place           `json:"current_location,omitempty"`
	RawSubjectOf    json.RawMessage `json:"subject_of,omitempty"`
	SubjectOf       []Object        `json:"-"`
	Reperesentation []Object        `json:"representation,omitempty"`
}
