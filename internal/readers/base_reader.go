package readers

import "linkedart-reader-golang/internal/models"

type BaseReader interface {
	GetObject(objectID *string) (*models.MuseumObject, error)
	GetNextObject(museumObject *models.MuseumObject) (*models.MuseumObject, error)
}
