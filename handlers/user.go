package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/telemetry"
	"github.com/hashicorp/go-hclog"
)

const jwtSecret = "test"

// User -
type User struct {
	log       hclog.Logger
	telemetry *telemetry.Telemetry
	con       data.Connection
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// AuthResponse -
type AuthResponse struct {
	UserID   int    `json:"user_id,omitempty`
	Username string `json:"username,omitempty`
	Token    string `json:"token,omitempty"`
}

// NewUser -
func NewUser(t *telemetry.Telemetry, l hclog.Logger, con data.Connection) *User {
	t.AddMeasure("user.signup")
	t.AddMeasure("user.signin")
	t.AddMeasure("user.generate_jwt_token")
	t.AddMeasure("user.invalidate_jwt_token")
	t.AddMeasure("user.signout")

	return &User{l, t, con}
}

func (c *User) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c.log.Info("Handle User | unknown", "path", r.URL.Path)
	http.NotFound(rw, r)
}

// SignUp registers a new user and returns a JWT token
// only restriction is username must be unique
func (c *User) SignUp(rw http.ResponseWriter, r *http.Request) {
	done := c.telemetry.NewTiming("user.signup")
	defer done()

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

	tokenString, err := c.generateJWTToken(u.ID, u.Username)
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
}

// SignIn signs in a user and returns a JWT token
func (c *User) SignIn(rw http.ResponseWriter, r *http.Request) {
	done := c.telemetry.NewTiming("user.signin")
	defer done()

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

	tokenString, err := c.generateJWTToken(u.ID, u.Username)
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
}

func (c *User) generateJWTToken(userID int, username string) (string, error) {
	done := c.telemetry.NewTiming("user.generate_jwt_token")
	defer done()

	t, err := c.con.CreateToken(userID)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token_id": t.ID,
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}

func (c *User) invalidateJWTToken(authToken string) error {
	done := c.telemetry.NewTiming("user.invalidate_jwt_token")
	defer done()

	tokenID, userID, err := ExtractJWT(authToken)
	if err != nil {
		return err
	}
	if err = c.con.DeleteToken(tokenID, userID); err != nil {
		return err
	}
	return nil
}

// SignOut signs out a user and invalidates a JWT token
func (c *User) SignOut(rw http.ResponseWriter, r *http.Request) {
	done := c.telemetry.NewTiming("user.signout")
	defer done()

	c.log.Info("Handle User | signout")

	authToken := r.Header.Get("Authorization")

	if err := c.invalidateJWTToken(authToken); err != nil {
		c.log.Error("Unable to sign out user", "error", err)
		http.Error(rw, "Unable to sign out user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "%s", "Signed out user")
}
