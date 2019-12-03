package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stylll/GoBlog/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestSignIn(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Firstname: "Andrew",
		Lastname:  "Benard",
		Email:     "andrew.benard@dundermifflin.com",
		Password:  "NardDog!",
	}

	err = seedSingleUser(&user)
	if err != nil {
		log.Fatal(err)
	}

	testCases := []struct {
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        user.Email,
			password:     "NardDog!",
			errorMessage: "",
		},
		{
			email:        user.Email,
			password:     "Bernard",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			email:        "wrong@email.com",
			password:     "Bernard",
			errorMessage: "record not found",
		},
	}

	for _, i := range testCases {
		token, err := server.SignIn(i.email, i.password)
		if err != nil {
			assert.Equal(t, errors.New(i.errorMessage), err)
		} else {
			assert.NotEqual(t, token, "")
		}
	}
}

func TestLogin(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Firstname: "Pam",
		Lastname:  "Beesly",
		Email:     "pam.beesly@dundermifflin.com",
		Password:  "Pamela20",
	}

	err = seedSingleUser(&user)
	if err != nil {
		log.Fatal(err)
	}

	testCases := []struct {
		inputJSON    string
		statusCode   int
		email        string
		password     string
		errorMessage string
	}{
		{
			inputJSON:    `{"email": "pam.beesly@dundermifflin.com", "password": "Pamela20"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "pam.beesly@dundermifflin.com", "password": "wrong password"}`,
			statusCode:   422,
			errorMessage: "Incorrect Password",
		},
		{
			inputJSON:    `{"email": "pam@dundermifflin.com", "password": "Pamela20"}`,
			statusCode:   422,
			errorMessage: "Incorrect details",
		},
		{
			inputJSON:    `{"email": "dundermifflin.com", "password": "Pamela20"}`,
			statusCode:   400,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"email": "", "password": "Pamela20"}`,
			statusCode:   400,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"email": "pam.best@dundermifflin.com", "password": ""}`,
			statusCode:   400,
			errorMessage: "Required Password",
		},
	}

	for _, i := range testCases {
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(i.inputJSON))
		if err != nil {
			t.Errorf("error occured: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, i.statusCode)
		if i.statusCode == 200 {
			assert.NotEqual(t, rr.Body.String(), "")
		}

		if i.statusCode == 422 && i.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err := json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert response to json: %v", err)
			}

			assert.Equal(t, responseMap["error"], i.errorMessage)
		}
	}
}
