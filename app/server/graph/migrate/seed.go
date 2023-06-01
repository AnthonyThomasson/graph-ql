package migrate

import "gorm.io/gorm"

func ResetData(db *gorm.DB) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()

	mustTruncateProducts(db)
	return
}

func Seed(db *gorm.DB) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()

	mustSeedProducts(db)
	return
}
