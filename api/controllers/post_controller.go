package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stylll/GoBlog/api/auth"
	"github.com/stylll/GoBlog/api/models"
	"github.com/stylll/GoBlog/api/responses"
	"github.com/stylll/GoBlog/api/utils/formaterror"
)

func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil || tokenID == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized Access"))
		return
	}

	if tokenID != int64(post.AuthorID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	post.Prepare()
	err = post.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	newPost, err := post.SavePost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, newPost.ID))
	responses.JSON(w, http.StatusCreated, newPost)
}

func (server *Server) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}

	allPosts, err := post.FindAllPosts(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, allPosts)
}

func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}
	vars := mux.Vars(r)

	postId, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	retrievedPost, err := post.FindPostByID(server.DB, int(postId))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, retrievedPost)
}

func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	post := models.Post{}
	foundPost, err := post.FindPostByID(server.DB, int(postId))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post Not Found"))
		return
	}

	if tokenID != int64(foundPost.AuthorID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// check that userId is the same as AuthorID in the new post to update
	if tokenID != int64(post.AuthorID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	post.Prepare()
	err = post.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	post.ID = int(postId) // set the post ID : not sure if this is necessary since post is retrieved from the db at first

	updatedPost, err := post.UpdateAPost(server.DB, int(postId))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, updatedPost)
}

func (server *Server) DeleteAPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	post := models.Post{}
	foundPost, err := post.FindPostByID(server.DB, int(postID))
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post Not Found"))
		return
	}

	if tokenID != int64(foundPost.AuthorID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = post.DeleteAPost(server.DB, int(postID), int(tokenID))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", postID))
	responses.JSON(w, http.StatusNoContent, "")
}
