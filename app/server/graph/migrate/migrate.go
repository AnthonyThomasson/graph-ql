package migrate

import "gorm.io/gorm"

func Migrate(db *gorm.DB) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()

	mustMigrateProducts(db)
	return
}
