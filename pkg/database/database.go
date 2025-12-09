package database

import (
	"SpeechAnalytics/pkg/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDB(dsn string) {
	// user := "questions"
	// pass := "0505"
	// host := "192.168.88.6"
	// name := "steamdb"
	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Moscow",
	// 	host,
	// 	user,
	// 	pass,
	// 	name,
	// )
	// dsn := fmt.Sprintf(
	// 	"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Moscow",
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_NAME"),
	// )

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database.\n", err)
		os.Exit(1)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migration")
	err = db.AutoMigrate(&models.Call{}, &models.SpeakerStatistics{})
	if err != nil {
		log.Fatal("migrations failed \n", err)
	}

	DB = Dbinstance{
		Db: db,
	}
}
