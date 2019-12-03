package controllers

import (
	"net/http"

	"github.com/stylll/GoBlog/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to GoBlog API")
}
