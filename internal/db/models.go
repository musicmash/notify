package db

import "github.com/jinzhu/gorm"

var tables = []interface{}{
	&LastAction{},
	&Chat{},
	&Notification{},
}

func CreateTables(db *gorm.DB) error {
	return db.AutoMigrate(tables...).Error
}

func DropAllTables(db *gorm.DB) error {
	return db.DropTable(tables...).Error
}

func CreateAll(db *gorm.DB) error {
	if err := CreateTables(db); err != nil {
		return err
	}

	fkeys := map[interface{}][][2]string{

		Chat{}: {
			{"user_name", "users(name)"},
		},
		Notification{}: {
			{"user_name", "users(name)"},
			{"release_id", "releases(id)"},
		},
	}

	for model, foreignKey := range fkeys {
		for _, fk := range foreignKey {
			if err := db.Debug().Model(model).AddForeignKey(
				fk[0], fk[1], "RESTRICT", "RESTRICT").Error; err != nil {
				return err
			}
		}
	}

	//if err := db.Debug().Model(&Subscription{}).AddUniqueIndex(
	//	"idx_user_id_artist_name",
	//	"user_id", "artist_name").Error; err != nil {
	//	return err
	//}
	//
	//if err := db.Debug().Model(&Release{}).AddIndex(
	//	"idx_store_id", "store_id").Error; err != nil {
	//	return err
	//}
	//
	//if err := db.Debug().Model(&Store{}).AddIndex(
	//	"idx_store_type_release_id",
	//	"store_type", "release_id").Error; err != nil {
	//	return err
	//}

	return nil
}
