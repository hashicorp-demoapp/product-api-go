package handlers

import (
	"math/rand"
	"net/http"

	"github.com/hashicorp/go-hclog"
)

// Middleware -
type ErrorRateMiddleware struct {
	errorRate int
	log       hclog.Logger
}

// NewMiddleware -
func NewErrorMiddleware(errorRate int, l hclog.Logger) *ErrorRateMiddleware {
	return &ErrorRateMiddleware{errorRate, l}
}

// error rate middleware - returns an error based on environment variable
func (emw *ErrorRateMiddleware) Middleware(next http.Handler) http.Handler {
	randInt := rand.Intn(100)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if emw.errorRate > randInt {
			emw.log.Error("Error rate triggered", "error", emw.errorRate)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
