package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AnthonyThomasson/graph-ql/api"
	m "github.com/AnthonyThomasson/graph-ql/graph/migrate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := connectToDB()
	migrate(db)
	seed(db)
	api.NewServer(db).Start(fmt.Sprintf(":%v", os.Getenv("SERV_PORT")))
}

func connectToDB() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	var db *gorm.DB
	var err error
	maxRetries := 5
	retryInterval := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			return db
		}

		log.Printf("Failed to connect to the database. Retrying in %v...", retryInterval)
		time.Sleep(retryInterval)
	}

	log.Fatal("Could not establish a database connection within the maximum retry limit.")
	return nil
}

func migrate(db *gorm.DB) {
	err := m.Migrate(db)
	if err != nil {
		log.Fatal("Failed to migrate:", err)
	}
}

func seed(db *gorm.DB) {
	shouldForceSeed := flag.Bool("f-seed", false, "The database will be seeded with test data regardless of whether it has already been seeded")
	shouldSeed := flag.Bool("seed", *shouldForceSeed, "The database will be seeded with test data if not already seeded")
	flag.Parse()
	if *shouldSeed {
		if *shouldForceSeed {
			err := m.ResetData(db)
			if err != nil {
				log.Fatal("Failed resetting the DB: ", err)
			}
		}
		err := m.Seed(db)
		if err != nil {
			log.Fatal("Failed seeding the DB: ", err)
		}
	}
}
