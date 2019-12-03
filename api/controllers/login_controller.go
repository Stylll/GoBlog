package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/stylll/GoBlog/api/auth"
	"github.com/stylll/GoBlog/api/models"
	"github.com/stylll/GoBlog/api/responses"
	"github.com/stylll/GoBlog/api/utils/formaterror"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {
	user := models.User{}

	err := server.DB.Debug().Model(&models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}

	err = models.VerifyPassword(user.Password, password)
	if err != nil {
		return "", err
	}

	return auth.CreateToken(user.ID)
}
