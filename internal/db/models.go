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

	fkeys := map[interface{}][][2]string{}

	for model, foreignKey := range fkeys {
		for _, fk := range foreignKey {
			if err := db.Debug().Model(model).AddForeignKey(
				fk[0], fk[1], "RESTRICT", "RESTRICT").Error; err != nil {
				return err
			}
		}
	}

	if err := db.Debug().Model(&Chat{}).AddUniqueIndex(
		"idx_chat_id_user_name",
		"id", "user_name").Error; err != nil {
		return err
	}
	if err := db.Debug().Model(&Notification{}).AddUniqueIndex(
		"idx_user_name_release_id_is_coming",
		"user_name", "release_id", "is_coming").Error; err != nil {
		return err
	}
	return nil
}
