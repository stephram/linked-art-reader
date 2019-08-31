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
		panic(fmt.Sprint("invalid ID '%s'", p.ID))
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
	Published time.Time `json:"published"'`
}

type Object struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	IdentifiedBy []Identifier `json:"identified_by,omitempty"`
}

type LinguisticObject struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Label        string `json:"_label,omitempty"`
	Format       string `json:"format,omitempty"`
	ClassifiedAs []Type `json:"classified_as,omitempty"`
}

type InformationObject struct {
	ID     string   `json:"id"`
	Type   string   `json:"type"`
	Format []string `json:"format,omitempty"`
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

type Person struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	IdentifiedBy []Identifier `json:"identified_by,omitempty"`
}

type Group struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	IdentifiedBy []Identifier `json:"identified_by,omitempty"`
}

type HumanMadeObject struct {
	ID            string             `json:"id"`
	Type          string             `json:"type"`
	Label         string             `json:"_label,omitempty"`
	IdentifiedBy  []Identifier       `json:"identified_by,omitempty"`
	ClassifiedAs  []Type             `json:"classified_as,omitempty"`
	ReferredToBy  []LinguisticObject `json:"referred_to_by,omitempty"`
	CurrentKeeper []Object           `json:"current_keeper,omitempty"`
	CurrentOwner  []Object           `json:"current_owner,omitempty"`
	RawSubjectOf  json.RawMessage    `json:"subject_of",omitempty`
	SubjectOf     []Object           `json:"-"`
}
