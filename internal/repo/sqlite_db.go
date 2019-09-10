package repo

import (
	"linked-art-reader/internal/models"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

type DbLocation struct {
	gorm.Model

	LocationID string `gorm:"type:varchar(256);unique_index"`
	UUID       string `gorm:"type:varchar(256);unique_index"`
	Location   string `gorm:"type:varchar(256)"`
}

type DbAltIdentifier struct {
	IdentifierID string `gorm:"type:varchar(128);unique_index"`
	Label        string `gorm:"type:varchar(256)"`
	Value        string `gorm:"type:varchar(256)"`
}

type DbAltReference struct {
	ReferenceID   string `gorm:"type:varchar(256);unique_index"`
	ReferenceName string `gorm:"type:varchar(256);index"`
}

type DbEntity struct {
	gorm.Model

	EntityID        string `gorm:"type:varchar(256);unique_index"`
	Type            string `gorm:"type:varchar(128)"`
	UUID            string
	DOR_ID          string
	TMS_ID          string
	AccessionNumber string `gorm:"type:varchar(128);index"`
	Title           string
	Content         string
	WebURL          string

	// AltLabels       map[string]string
	AltReferences  []DbAltReference  `gorm:"ForeignKey:ReferenceID"`
	AltIdentifiers []DbAltIdentifier `gorm:"ForeignKey:IdentifierID"`

	Location DbLocation `gorm:"ForeignKey:LocationID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type larRepoImpl struct {
	db *gorm.DB
}

func (u *larRepoImpl) FindEntity(entityID string) (*models.Entity, error) {
	var dbEntity DbEntity

	query := u.db.Where("id = ?", entityID)
	err := query.Find(&dbEntity).Error

	return u.toModelEntity(&dbEntity), err
}

func (u *larRepoImpl) StoreEntity(entity *models.Entity) (*models.Entity, error) {

	dbEntity := u.toDbEntity(entity)

	err := u.db.Where(DbEntity{EntityID: entity.ID}).Assign(u.toDbEntity(entity)).FirstOrCreate(&dbEntity).Error
	return u.toModelEntity(dbEntity), err
}

func (u *larRepoImpl) toModelEntity(dbEntity *DbEntity) *models.Entity {
	return &models.Entity{
		ID:              dbEntity.EntityID,
		UUID:            dbEntity.UUID,
		DOR_ID:          dbEntity.DOR_ID,
		TMS_ID:          dbEntity.TMS_ID,
		Type:            dbEntity.Type,
		AccessionNumber: dbEntity.AccessionNumber,
		Title:           dbEntity.Title,
	}
}

func (u *larRepoImpl) toDbEntity(entity *models.Entity) *DbEntity {
	return &DbEntity{
		EntityID:        entity.ID,
		Type:            entity.Type,
		UUID:            entity.UUID,
		DOR_ID:          entity.DOR_ID,
		AccessionNumber: entity.AccessionNumber,
		TMS_ID:          entity.TMS_ID,
		Title:           entity.Title,
		// Content: entity.Content,
		WebURL:         entity.WebURL,
		AltReferences:  *(u.makeDbReferences(entity.AltReferences)),
		AltIdentifiers: *(u.makeDbIdentifiers(entity.AltIDs)),
	}
}

func (u *larRepoImpl) makeDbReferences(references map[string]string) *[]DbAltReference {
	refs := []DbAltReference{}

	for k, v := range references {
		refs = append(refs, DbAltReference{
			ReferenceID:   v,
			ReferenceName: k,
		})
	}
	return &refs
}

func (u *larRepoImpl) makeDbIdentifiers(identifiers map[string]string) *[]DbAltIdentifier {
	ids := []DbAltIdentifier{}

	for k, v := range identifiers {
		ids = append(ids, DbAltIdentifier{
			IdentifierID: k,
			Label:        v,
			Value:        v,
		})
	}
	return &ids
}

func New(db *gorm.DB) LinkedArtReaderRepo {
	larRepo := &larRepoImpl{
		db: func() *gorm.DB {
			if db == nil {
				return createDB(true)
			}
			return db
		}(),
	}
	return larRepo
}

func createDB(testMode bool) *gorm.DB {
	return func() *gorm.DB {
		// db, err := gorm.Open("sqlite3", "file:larDB?mode=memory&cache=shared")
		db, err := gorm.Open("sqlite3", "./larDB.sqlite")
		if err != nil {
			log.WithError(err).Errorf("unable to open database")
			return nil
		}
		db.LogMode(true)
		if testMode {
			db.AutoMigrate(&DbEntity{})
			db.AutoMigrate(&DbLocation{})
			db.AutoMigrate(&DbAltReference{})
			db.AutoMigrate(&DbAltIdentifier{})

			db.Model(&DbAltIdentifier{}).AddForeignKey("entity_id", "db_entities(entity_id)", "CASCADE", "CASCADE")
			db.Model(&DbAltReference{}).AddForeignKey("entity_id", "db_entities(entity_id)", "CASCADE", "CASCADE")
		}
		db.BlockGlobalUpdate(true)

		log.Infof("opened sqlite database")
		return db
	}()
}
