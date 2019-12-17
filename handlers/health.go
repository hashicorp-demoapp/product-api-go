package handlers

import (
	"fmt"
	"net/http"

	"github.com/hashicorp-demoapp/product-api-go/telemetry"
	"github.com/hashicorp/go-hclog"
)

// Health is a HTTP Handler for health checking
type Health struct {
	logger    hclog.Logger
	telemetry *telemetry.Telemetry
}

// NewHealth creates a new Health handler
func NewHealth(t *telemetry.Telemetry, l hclog.Logger) *Health {
	t.AddMeasure("health.call")

	return &Health{l, t}
}

// Handle the request
func (h *Health) Handle(rw http.ResponseWriter, r *http.Request) {
	done := h.telemetry.NewTiming("health.call")
	defer done()

	fmt.Fprintf(rw, "%s", "ok")
}
