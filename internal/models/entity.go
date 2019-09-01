package models

type Entity struct {
	ID       string
	UID      string
	UUID     string
	DORID    string
	Label    string
	Title    string
	Name     string
	AltNames []string
	Content  interface{}
	WebURL   string
}
