package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gym-app/models"
)

var DB *gorm.DB

func ConnectDB() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=postgres_db user=abs password=yourpassword dbname=gym port=5432 sslmode=disable TimeZone=UTC"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("connected")
	DB = db

	log.Println("running migrations")
	// The DB is pre-populated by init.sql but we can run AutoMigrate to ensure they align
	db.AutoMigrate(&models.User{}, &models.WorkoutSession{}, &models.Exercise{}, &models.ExerciseLog{}, &models.Set{})
}
