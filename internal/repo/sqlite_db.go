package repo

import (
	"fmt"
	"linked-art-reader/internal/models"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

type DbLocation struct {
	gorm.Model
	LocationID string `gorm:"primary_key"`
	UUID       string
	Location   string
	EntityRef  uint
}

type DbLabel struct {
	gorm.Model
	LabelID   string `gorm:"primary_key"`
	Label     string
	Value     string
	EntityRef uint
}

type DbIdentifier struct {
	gorm.Model
	IdentifierID string `gorm:"primary_key"`
	Label        string `gorm:"type:varchar(256)"`
	Value        string `gorm:"type:varchar(256)"`
	EntityRef    uint
}

type DbReference struct {
	gorm.Model
	ReferenceID    string `gorm:"primary_key"`
	ReferenceName  string
	ReferenceValue string
	EntityRef      uint
}

type DbClassifier struct {
	gorm.Model
	ClassifierID    string `gorm:"primary_key"`
	ClassifierName  string
	ClassifierValue string
	EntityRef       uint
}

type DbEntity struct {
	gorm.Model
	EntityID        string `gorm:"primary_key"`
	Type            string
	UUID            string
	DOR_ID          string
	TMS_ID          string
	AccessionNumber string
	Title           string
	Content         string
	WebURL          string

	Labels      []DbLabel      `gorm:"foreignkey:EntityRef"`
	References  []DbReference  `gorm:"foreignkey:EntityRef"`
	Identifiers []DbIdentifier `gorm:"foreignkey:EntityRef"`
	Classifiers []DbClassifier `gorm:"foreignkey:EntityRef"`
	Location    DbLocation     `gorm:"foreignkey:EntityRef"`
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
		Content:         fmt.Sprint(entity.Content),
		WebURL:          entity.WebURL,
		Labels:          *(u.makeDbLabels(entity.Labels)),
		References:      *(u.makeDbReferences(entity.AltReferences)),
		Identifiers:     *(u.makeDbIdentifiers(entity.AltIdentifiers)),
		Location:        *(u.makeDbLocation(entity.Location)),
		Classifiers:     *(u.makeDbClassifiers(entity.Classifiers)),
	}
}

func (u *LinkedArtReaderRepoImpl) makeDbLocation(location models.Location) *DbLocation {
	return &DbLocation{
		LocationID: location.ID,
		UUID:       location.UUID,
		Location:   location.Location,
	}
}

func (u *LinkedArtReaderRepoImpl) makeDbClassifiers(classifiers map[string]string) *[]DbClassifier {
	_classifiers := []DbClassifier{}

	for k, v := range classifiers {
		_classifiers = append(_classifiers, DbClassifier{
			ClassifierID:    k,
			ClassifierName:  v,
			ClassifierValue: v,
		})
	}
	return &_classifiers
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
			ReferenceID:    k,
			ReferenceName:  v,
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

func createDB(create bool) *gorm.DB {
	return func() *gorm.DB {
		dbPath := "./larDB.sqlite"
		basePath := os.Getenv("BASE_PATH")
		if basePath != "" {
			dbPath = basePath + "/larDB.sqlite"
		}

		// db, err := gorm.Open("sqlite3", "file:larDB?mode=memory&cache=shared")
		db, err := gorm.Open("sqlite3", dbPath)
		if err != nil {
			log.WithError(err).Errorf("unable to open database")
			return nil
		}
		db.LogMode(true)

		if create {
			dbErr := db.DropTableIfExists(&DbEntity{}, &DbLocation{}, &DbReference{}, &DbIdentifier{}, &DbLabel{}, &DbClassifier{}).Error
			if dbErr != nil {
				log.Errorf("DropTableIfExists failed: %s", dbErr.Error())
			}
			dbErr = db.AutoMigrate(&DbEntity{}, &DbLocation{}, &DbReference{}, &DbIdentifier{}, &DbLabel{}, &DbClassifier{}).Error
			if dbErr != nil {
				log.Errorf("AutoMigrate failed: %s", dbErr.Error())
			}
		}
		db.BlockGlobalUpdate(true)

		log.Infof("opened sqlite database")
		return db
	}()
}
