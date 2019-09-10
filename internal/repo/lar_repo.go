package repo

import "linked-art-reader/internal/models"

type LinkedArtReaderRepo interface {
	FindEntity(entityID string) (*models.Entity, error)
	StoreEntity(entity *models.Entity) (*models.Entity, error)
}
