package telemetry

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/exporter/metric/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/batcher/defaultkeys"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"

	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

type Telemetry struct {
	pusher   *controller.Controller
	meter    api.Meter
	measures map[string]*api.Float64Measure
	counters map[string]*api.Float64Counter
}

func New(bind_address string) *Telemetry {
	selector := simple.NewWithExactMeasure()
	exporter, err := prometheus.NewExporter(prometheus.Options{})

	if err != nil {
		log.Panicf("failed to initialize metric stdout exporter %v", err)
	}

	batcher := defaultkeys.New(selector, sdkmetric.NewDefaultLabelEncoder(), false)
	pusher := controller.New(batcher, exporter, time.Second)
	pusher.Start()

	go func() {
		_ = http.ListenAndServe(bind_address, exporter)
	}()

	global.SetMeterProvider(pusher)

	meter := global.MeterProvider().Meter("ex.com/basic")

	m := make(map[string]*api.Float64Measure)
	c := make(map[string]*api.Float64Counter)

	return &Telemetry{pusher, meter, m, c}
}

// AddMeasure to metrics
func (t *Telemetry) AddMeasure(key string) {
	met := t.meter.NewFloat64Measure(key)
	t.measures[key] = &met
}

// AddCounter to metrics collection
func (t *Telemetry) AddCounter(key string) {
	met := t.meter.NewFloat64Counter(key)
	t.counters[key] = &met
}

// NewTiming creates a new timing metric and returns a done function
func (t *Telemetry) NewTiming(key string) func() {
	// record the start time
	st := time.Now()

	return func() {
		dur := time.Now().Sub(st).Nanoseconds()
		handler := t.measures[key].AcquireHandle(nil)
		defer handler.Release()

		t.meter.RecordBatch(
			context.Background(),
			nil,
			t.measures[key].Measurement(float64(dur)),
		)
	}
}
