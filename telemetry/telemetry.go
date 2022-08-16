package telemetry

import (
	"context"
	"log"

	"github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

type Telemetry struct {
	// pusher *controller.Controller
	meter  metric.Meter
	Tracer opentracing.Tracer
	// measures map[string]*syncfloat64.InstrumentProvider
	// counters map[string]*syncfloat64.InstrumentProvider
}

func New(bind_address string) *Telemetry {
	// selector := simple.NewWithExactMeasure()
	// exporter, err := prometheus.NewExporter(prometheus.Options{})

	// if err != nil {
	// 	log.Panicf("failed to initialize metric stdout exporter %v", err)
	// }

	// batcher := defaultkeys.New(selector, sdkmetric.NewDefaultLabelEncoder(), false)
	// pusher := push.New(batcher, exporter, time.Second)
	// pusher.Start()

	// go func() {
	// 	_ = http.ListenAndServe(bind_address, exporter)
	// }()

	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter: %w", err)
	}

	global.SetMeterProvider(exporter.MeterProvider())

	meter := global.MeterProvider().Meter("ex.com/basic")
	tracer := opentracing.GlobalTracer()

	return &Telemetry{meter, tracer}

	// m := make(map[string]*meter.Float64Measure)
	// c := make(map[string]*meter.Float64Counter)

	// return &Telemetry{pusher, meter, m, c}
}

// AddMeasure to metrics
func (t *Telemetry) AddMeasure(key string) {
	ctx := context.Background()

	counter, _ := t.meter.SyncFloat64().Counter(key)
	counter.Add(ctx, 1)
	// met := t.meter.NewFloat64Measure(key)
	// t.measures[key] = &met
}

// AddCounter to metrics collection
func (t *Telemetry) AddCounter(key string) {
	ctx := context.Background()

	counter, _ := t.meter.SyncFloat64().Counter(key)
	counter.Add(ctx, 1)
	// met := t.meter.NewFloat64Counter(key)
	// t.counters[key] = &met
}

// NewTiming creates a new timing metric and returns a done function
func (t *Telemetry) NewTiming(key string) func() {
	// record the start time
	// st := time.Now()

	return func() {
		// dur := time.Now().Sub(st).Nanoseconds()
		// handler := t.measures[key].AcquireHandle(nil)
		// defer handler.Release()

		// t.meter.RecordBatch(
		// 	context.Background(),
		// 	nil,
		// 	t.measures[key].Measurement(float64(dur)),
		// )
	}
}
