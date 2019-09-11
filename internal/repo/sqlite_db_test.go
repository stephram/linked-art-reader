package repo

import (
	"linked-art-reader/internal/utils/ulid"
	"os"
	"testing"

	"linked-art-reader/internal/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/stretchr/testify/assert"
)

var (
	db       *gorm.DB
	larRepo  LinkedArtReaderRepo
	entityID string
)

func TestMain(m *testing.M) {
	db = createDB(true)
	defer db.Close()

	larRepo = New(db)
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	t.Run("store and retrieve ok", func(t *testing.T) {
		id := ulid.New()

		var entity *models.Entity
		var err error

		if entity, err = larRepo.StoreEntity(&models.Entity{
			ID: id,
		}); err != nil {
			assert.FailNow(t, err.Error())
		}

		assert.NotNil(t, entity)
		assert.Equal(t, id, entity.ID)

		if entity, err = larRepo.FindEntity(id); err != nil {
			assert.FailNow(t, err.Error())
		}

		assert.NotNil(t, entity)
		assert.Equal(t, id, entity.ID)
	})

	// Haven't implemented FindEntity yet.
	// t.Run("test FindEntity", func(t *testing.T) {
	// 	var entity *models.Entity
	// 	var err error
	//
	// 	if entity, err = larRepo.FindEntity(entityID); err != nil {
	// 		assert.FailNow(t, err.Error())
	// 	}
	// 	assert.NotNil(t, entity)
	// 	assert.Equal(t, entityID, entity.ID)
	// })
}
