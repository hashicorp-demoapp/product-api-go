package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp/go-hclog"
)

const jwtSecret = "test"

type User struct {
	con data.Connection
	log hclog.Logger
}

type AuthStruct struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthResponse struct {
	UserID   int    `json:"user_id,omitempty`
	Username string `json:"username,omitempty`
	Token    string `json:"token,omitempty"`
}

func NewUser(con data.Connection, l hclog.Logger) *User {
	return &User{con, l}
}

func (c *User) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/signup":
		c.log.Info("Handle User | signup")

		body := AuthStruct{}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			c.log.Error("Unable to decode JSON", "error", err)
			http.Error(rw, "Unable to parse request body", http.StatusInternalServerError)
			return
		}

		u, err := c.con.CreateUser(body.Username, body.Password)
		if err != nil {
			c.log.Error("Unable to create new user", "error", err)
			if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
				http.Error(rw, fmt.Sprintf("User already exists: %s", body.Username), http.StatusInternalServerError)
				return
			}
			http.Error(rw, fmt.Sprintf("Unable to sign up user: %s", body.Username), http.StatusInternalServerError)
			return
		}

		tokenString, err := generateJWTToken(u.ID, u.Username)
		if err != nil {
			c.log.Error("Unable to generate JWT token", "error", err)
			http.Error(rw, "Unable to generate JWT token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(rw).Encode(AuthResponse{
			UserID:   u.ID,
			Username: u.Username,
			Token:    tokenString,
		})
	case "/signin":
		c.log.Info("Handle User | signin")

		body := AuthStruct{}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			c.log.Error("Unable to decode JSON", "error", err)
			http.Error(rw, "Unable to parse request body", http.StatusInternalServerError)
			return
		}

		u, err := c.con.AuthUser(body.Username, body.Password)
		if err != nil {
			c.log.Error("Unable to sign in user", "error", err)
			http.Error(rw, "Invalid Credentials", http.StatusUnauthorized)
			return
		}

		tokenString, err := generateJWTToken(u.ID, u.Username)
		if err != nil {
			c.log.Error("Unable to generate JWT token", "error", err)
			http.Error(rw, "Unable to generate JWT token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(rw).Encode(AuthResponse{
			UserID:   u.ID,
			Username: u.Username,
			Token:    tokenString,
		})
	default:
		c.log.Info("Handle User | unknown", "path", r.URL.Path)
	}
}

func generateJWTToken(id int, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  id,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}
