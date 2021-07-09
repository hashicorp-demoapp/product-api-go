package handlers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp/go-hclog"
)

// Middleware -
type AuthMiddleware struct {
	con data.Connection
	log hclog.Logger
}

// NewMiddleware -
func NewAuthMiddleware(con data.Connection, l hclog.Logger) *AuthMiddleware {
	return &AuthMiddleware{con, l}
}

// ExtractJWT retrieves the token and user ID from the JWT
func ExtractJWT(authToken string) (int, int, error) {
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return -1, -1, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenID := int(claims["token_id"].(float64))
		userID := int(claims["user_id"].(float64))
		return tokenID, userID, nil
	}
	return -1, -1, nil
}

func (c *AuthMiddleware) VerifyJWT(authToken string) (int, error) {
	tokenID, userID, err := ExtractJWT(authToken)
	if err != nil {
		return userID, err
	}
	if _, err := c.con.GetToken(tokenID, userID); err != nil {
		return userID, err
	}
	return userID, nil
}

// IsAuthorized
func (c *AuthMiddleware) IsAuthorized(next func(userID int, w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		userID, err := c.VerifyJWT(authToken)
		if err == nil {
			next(userID, w, r)
			return
		}
		c.log.Error("Unauthorized", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	})
}
