package handlers

import (
	"fmt"
	"net/http"

	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/telemetry"
	"github.com/hashicorp/go-hclog"
)

// Health is a HTTP Handler for health checking
type Health struct {
	logger    hclog.Logger
	telemetry *telemetry.Telemetry
	db        data.Connection
}

// NewHealth creates a new Health handler
func NewHealth(t *telemetry.Telemetry, l hclog.Logger, db data.Connection) *Health {
	t.AddMeasure("health.call")
	t.AddMeasure("health.livez")
	t.AddMeasure("health.readyz")

	return &Health{l, t, db}
}

// ServeHTTP implements the handler interface
// Deprecated: Use Liveness and Readiness handlers instead
func (h *Health) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	done := h.telemetry.NewTiming("health.call")
	defer done()

	_, err := h.db.IsConnected()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "error %s", err)
	}

	fmt.Fprintf(rw, "%s", "ok")
}

// Liveness endpoint for health checks indicates server has started
func (h *Health) Liveness(rw http.ResponseWriter, r *http.Request) {
	done := h.telemetry.NewTiming("health.livez")
	defer done()

	fmt.Fprintf(rw, "%s", "ok")
}

// Readiness endpoint for health checks indicates server is ready to serve traffic
func (h *Health) Readiness(rw http.ResponseWriter, r *http.Request) {
	done := h.telemetry.NewTiming("health.readyz")
	defer done()

	_, err := h.db.IsConnected()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "error %s", err)
	}

	fmt.Fprintf(rw, "%s", "ok")
}
