package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stylll/GoBlog/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	DBURL := "host=%s port=%s user=%s dbname=%s sslmode=disable password=%s"
	DBURL = fmt.Sprintf(DBURL, DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open("postgres", DBURL)

	if err != nil {
		fmt.Print("Cannot connect to database")
		log.Fatal("Error occured: ", err)
	} else {
		fmt.Print("Connected to database")
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{})

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(address string) {
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(address, server.Router))
}
