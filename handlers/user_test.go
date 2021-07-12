package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func setupUserHandler(t *testing.T) (*User, *httptest.ResponseRecorder) {
	c := &data.MockConnection{}

	c.On("CreateUser").Return(model.User{ID: 1, Username: "User1"}, nil)
	c.On("AuthUser").Return(model.User{ID: 1, Username: "User1"}, nil)
	c.On("CreateToken").Return(model.Token{ID: 2, UserID: 1}, nil)
	c.On("DeleteToken").Return(nil)

	l := hclog.Default()

	return &User{c, l}, httptest.NewRecorder()
}

func setupFailedUserHandler(t *testing.T) (*User, *httptest.ResponseRecorder) {
	c := &data.MockConnection{}

	c.On("CreateUser").Return(nil, errors.New("Unable to create new user"))
	c.On("AuthUser").Return(nil, errors.New("Unable to login with credentials"))
	c.On("CreateToken").Return(nil, errors.New("Unable to create new token"))
	c.On("DeleteToken").Return(nil, errors.New("Unable to delete token"))

	l := hclog.Default()

	return &User{c, l}, httptest.NewRecorder()
}

func TestCreateNewUser(t *testing.T) {
	c, rw := setupUserHandler(t)

	r := httptest.NewRequest("POST", "/signup", nil)

	rb := strings.NewReader(`{"username": "User1", "password": "testPassword"}`)
	r.Body = ioutil.NopCloser(rb)

	c.SignUp(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.User{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)
	assert.NoError(t, err)
}

func TestAuthNewUser(t *testing.T) {
	c, rw := setupUserHandler(t)

	r := httptest.NewRequest("POST", "/signin", nil)

	rb := strings.NewReader(`{"username": "User1", "password": "testPassword"}`)
	r.Body = ioutil.NopCloser(rb)

	c.SignIn(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := model.User{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)
	assert.NoError(t, err)
}

func TestSignOutUser(t *testing.T) {
	c, rw := setupUserHandler(t)

	token, err := c.generateJWTToken(1, "User1")
	fmt.Println(token)
	assert.NoError(t, err)

	r := httptest.NewRequest("POST", "/signout", nil)
	r.Header.Add("Authorization", token)

	c.SignOut(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
}

func TestUnableToCreateNewUser(t *testing.T) {
	c, rw := setupFailedUserHandler(t)

	r := httptest.NewRequest("POST", "/signup", nil)

	username := "User1"

	rb := strings.NewReader(fmt.Sprintf(`{"username": "%+s", "password": "testPassword"}`, username))
	r.Body = ioutil.NopCloser(rb)

	c.SignUp(rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)

	bd := model.User{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("Unable to sign up user: %+s\n", username), string(rw.Body.Bytes()))
}

func TestUnableToCreateNewToken(t *testing.T) {
	c, rw := setupFailedUserHandler(t)

	r := httptest.NewRequest("POST", "/signup", nil)

	username := "User1"

	rb := strings.NewReader(fmt.Sprintf(`{"username": "%+s", "password": "testPassword"}`, username))
	r.Body = ioutil.NopCloser(rb)

	c.SignUp(rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)

	bd := model.User{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("Unable to sign up user: %+s\n", username), string(rw.Body.Bytes()))
}

func TestUnableToAuthNewUser(t *testing.T) {
	c, rw := setupFailedUserHandler(t)

	r := httptest.NewRequest("POST", "/signin", nil)

	rb := strings.NewReader(`{"username": "User1", "password": "testPassword"}`)
	r.Body = ioutil.NopCloser(rb)

	c.SignIn(rw, r)

	assert.Equal(t, http.StatusUnauthorized, rw.Code)

	bd := model.User{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)
	assert.Error(t, err)
	assert.Equal(t, "Invalid Credentials\n", string(rw.Body.Bytes()))
}

func TestUnableToSignOutUser(t *testing.T) {
	c, rw := setupFailedUserHandler(t)

	r := httptest.NewRequest("POST", "/signout", nil)
	r.Header.Add("Authorization", "{}")

	c.SignOut(rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)
	assert.Equal(t, fmt.Sprintf("Unable to sign out user\n"), string(rw.Body.Bytes()))
}
