package repo

import (
	"linked-art-reader/internal/models"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

type DbLocation struct {
	LocationID string `gorm:"primary_key"`
	UUID       string `gorm:"type:varchar(256)"`
	Location   string `gorm:"type:varchar(256)"`
}

type DbLabel struct {
	LabelID string `gorm:"primary_key"`
	Label   string
	Value   string
}

type DbIdentifier struct {
	IdentifierID string `gorm:"primary_key"`
	Label        string `gorm:"type:varchar(256)"`
	Value        string `gorm:"type:varchar(256)"`
}

type DbReference struct {
	ReferenceID    string `gorm:"primary_key"`
	ReferenceName  string
	ReferenceValue string
}

type DbEntity struct {
	EntityID        string `gorm:"primary_key"`
	Type            string
	UUID            string
	DOR_ID          string
	TMS_ID          string
	AccessionNumber string
	Title           string
	Content         string
	WebURL          string

	Labels      []DbLabel      `gorm:"ForeignKey:LabelID"`
	References  []DbReference  `gorm:"ForeignKey:ReferenceID"`
	Identifiers []DbIdentifier `gorm:"ForeignKey:IdentifierID"`
	// Location DbLocation `gorm:"ForeignKey:LocationID"`

	UpdatedAt time.Time
	DeletedAt time.Time
}

type LinkedArtReaderRepoImpl struct {
	db *gorm.DB
}

func (u *LinkedArtReaderRepoImpl) FindEntity(entityID string) (*models.Entity, error) {
	var dbEntity DbEntity

	query := u.db.Where("entity_id = ?", entityID)
	err := query.Find(&dbEntity).Error

	return u.toModelEntity(&dbEntity), err
}

func (u *LinkedArtReaderRepoImpl) StoreEntity(entity *models.Entity) (*models.Entity, error) {
	dbEntity := u.toDbEntity(entity)
	err := u.db.Where(DbEntity{EntityID: entity.ID}).Assign(u.toDbEntity(entity)).FirstOrCreate(&dbEntity).Error
	return u.toModelEntity(dbEntity), err
}

func (u *LinkedArtReaderRepoImpl) toModelEntity(dbEntity *DbEntity) *models.Entity {
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

func (u *LinkedArtReaderRepoImpl) toDbEntity(entity *models.Entity) *DbEntity {
	return &DbEntity{
		EntityID:        entity.ID,
		Type:            entity.Type,
		UUID:            entity.UUID,
		DOR_ID:          entity.DOR_ID,
		AccessionNumber: entity.AccessionNumber,
		TMS_ID:          entity.TMS_ID,
		Title:           entity.Title,
		// Content: entity.Content,
		WebURL:      entity.WebURL,
		Labels:      *(u.makeDbLabels(entity.Labels)),
		References:  *(u.makeDbReferences(entity.AltReferences)),
		Identifiers: *(u.makeDbIdentifiers(entity.AltIdentifiers)),
	}
}

func (u *LinkedArtReaderRepoImpl) makeDbLabels(labels map[string]string) *[]DbLabel {
	_labels := []DbLabel{}

	for k, v := range labels {
		_labels = append(_labels, DbLabel{
			Label: k,
			Value: v,
		})
	}
	return &_labels
}

func (u *LinkedArtReaderRepoImpl) makeDbReferences(references map[string]string) *[]DbReference {
	refs := []DbReference{}

	for k, v := range references {
		refs = append(refs, DbReference{
			ReferenceID:    v,
			ReferenceName:  k,
			ReferenceValue: v,
		})
	}
	return &refs
}

func (u *LinkedArtReaderRepoImpl) makeDbIdentifiers(identifiers map[string]string) *[]DbIdentifier {
	ids := []DbIdentifier{}

	for k, v := range identifiers {
		ids = append(ids, DbIdentifier{
			IdentifierID: k,
			Label:        v,
			Value:        v,
		})
	}
	return &ids
}

func New(db *gorm.DB) LinkedArtReaderRepo {
	larRepo := &LinkedArtReaderRepoImpl{
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
			dbErr := db.DropTableIfExists(&DbEntity{}, &DbLocation{}, &DbReference{}, &DbIdentifier{}, &DbLabel{}).Error
			if dbErr != nil {
				log.Errorf("DropTableIfExists failed: %s", dbErr.Error())
			}
			dbErr = db.AutoMigrate(&DbEntity{}, &DbLocation{}, &DbReference{}, &DbIdentifier{}, &DbLabel{}).Error
			if dbErr != nil {
				log.Errorf("AutoMigrate failed: %s", dbErr.Error())
			}
			dbErr = db.Model(&DbLocation{}).AddForeignKey("location_id", "db_entities(entity_id)", "CASCADE", "CASCADE").Error
			if dbErr != nil {
				log.Errorf("AddForeignKey failed: %s", dbErr.Error())
			}
			dbErr = db.Model(&DbIdentifier{}).AddForeignKey("identifier_id", "db_entities(entity_id)", "CASCADE", "CASCADE").Error
			if dbErr != nil {
				log.Errorf("AddForeignKey failed: %s", dbErr.Error())
			}
			dbErr = db.Model(&DbReference{}).AddForeignKey("reference_id", "db_entities(entity_id)", "CASCADE", "CASCADE").Error
			if dbErr != nil {
				log.Errorf("AddForeignKey failed: %s", dbErr.Error())
			}
			dbErr = db.Model(&DbLabel{}).AddForeignKey("label_id", "db_entities(entity_id)", "CASCADE", "CASCADE").Error
			if dbErr != nil {
				log.Errorf("AddForeignKey failed: %s", dbErr.Error())
			}
		}
		db.BlockGlobalUpdate(true)

		log.Infof("opened sqlite database")
		return db
	}()
}
