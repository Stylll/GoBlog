package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/stylll/GoBlog/api/models"
)

var users = []models.User{
	models.User{
		Email:     "john@yah.com",
		Firstname: "John",
		Lastname:  "Andrew",
		Password:  "JAndrew",
	},
	models.User{
		Email:     "mark@yah.com",
		Firstname: "Mark",
		Lastname:  "Donalds",
		Password:  "MDonalds",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "My Adventure",
		Content: "I travelled the world in 60 days",
	},
	models.Post{
		Title:   "How to Write Code",
		Content: "Software programming is an interesting profession",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("Cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("Cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("Attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed user table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed post table: %v", err)
		}
	}
}
