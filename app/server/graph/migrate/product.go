package migrate

import (
	"github.com/AnthonyThomasson/graph-ql/graph/model"
	"gorm.io/gorm"
)

func mustMigrateProducts(db *gorm.DB) {
	err := db.AutoMigrate(model.Product{})
	if err != nil {
		panic(err)
	}
}

func mustTruncateProducts(db *gorm.DB) {
	r := db.Exec("TRUNCATE TABLE products")
	if r.Error != nil {
		panic(r.Error)
	}
}

func mustSeedProducts(db *gorm.DB) {
	var count int64
	if err := db.Model(&model.Product{}).Count(&count).Error; err != nil {
		panic(err)
	}
	if count > 0 {
		return
	}

	products := []model.Product{
		{Name: "iPhone 13", Price: 1299.99, Rank: 1},
		{Name: "Samsung Galaxy S21", Price: 1099.99, Rank: 2},
		{Name: "MacBook Pro", Price: 2399.99, Rank: 3},
		{Name: "Sony PlayStation 5", Price: 499.99, Rank: 4},
		{Name: "Nintendo Switch", Price: 299.99, Rank: 5},
		{Name: "Bose QuietComfort 35 II", Price: 349.99, Rank: 6},
		{Name: "Canon EOS R6", Price: 2499.99, Rank: 7},
		{Name: "Dyson V11 Absolute", Price: 599.99, Rank: 8},
		{Name: "Nike Air Zoom Pegasus 38", Price: 119.99, Rank: 9},
		{Name: "Levis 501 Original Fit Jeans", Price: 59.99, Rank: 10},
	}

	r := db.Create(&products)
	if r.Error != nil {
		panic(r.Error)
	}
}
