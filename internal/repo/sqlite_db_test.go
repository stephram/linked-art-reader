package repo

import (
	"linked-art-reader/internal/utils/ulid"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"linked-art-reader/internal/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/stretchr/testify/assert"
)

var (
	db      *gorm.DB
	larRepo LinkedArtReaderRepo
)

func TestMain(m *testing.M) {
	db = createDB(true, "testDB.sqlite")
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

	t.Run("test FindEntity", func(t *testing.T) {
		var entity *models.Entity = &models.Entity{
			ID:    ulid.New(),
			Title: "test",
		}
		var err error

		entity, err = storeEntity(entity)
		assert.Nil(t, err)
		require.Nil(t, err)

		entityID := entity.ID

		if entity, err = larRepo.FindEntity(entityID); err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.NotNil(t, entity)
		assert.Equal(t, entityID, entity.ID)
		assert.Equal(t, "test", entity.Title)
	})
}

func storeEntity(entity *models.Entity) (*models.Entity, error) {
	return larRepo.StoreEntity(entity)
}
