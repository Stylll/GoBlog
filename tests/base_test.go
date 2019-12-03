package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/stylll/GoBlog/api/controllers"
	"github.com/stylll/GoBlog/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}
var postInstance = models.Post{}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	setupDatabase()

	os.Exit(m.Run())
}

func setupDatabase() {
	var err error
	DBURL := "host=%s port=%s user=%s dbname=%s sslmode=disable password=%s"
	DBURL = fmt.Sprintf(DBURL, os.Getenv("DB_HOST_TEST"), os.Getenv("DB_PORT_TEST"),
		os.Getenv("DB_USER_TEST"), os.Getenv("DB_NAME_TEST"), os.Getenv("DB_PASSWORD_TEST"))
	server.DB, err = gorm.Open("postgres", DBURL)
	if err != nil {
		fmt.Print("Cannot connect to database")
		log.Fatal("Error occured: ", err)
	} else {
		fmt.Print("Connected to database")
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}

	log.Printf("Users table refreshed successfully")

	return nil
}

func refreshPostTable() error {
	err := server.DB.DropTableIfExists(&models.Post{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.Post{}).Error
	if err != nil {
		return err
	}

	log.Printf("Posts table refreshed successfully")

	return nil
}

func seedSingleUser(user *models.User) error {
	err := server.DB.Debug().Model(&models.User{}).Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

func seedSinglePost(post *models.Post) error {
	err := server.DB.Debug().Model(&models.Post{}).Create(post).Error
	if err != nil {
		return err
	}

	return nil
}
