// package telemetry contains the otel telemetry setup for the service
package telemetry

import (
	"context"
	"errors"

	hostMetrics "go.opentelemetry.io/contrib/instrumentation/host"
	runtimeMetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/contrib/processors/minsev"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func Setup(ctx context.Context, environment string, logLevel minsev.Severity) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error

	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	res, err := resource.New(ctx, resource.WithHost(), resource.WithAttributes(
		semconv.ServiceName("joke-service"),
		semconv.DeploymentEnvironmentName(environment),
	))
	if err != nil {
		return shutdown, err
	}

	tp, err := newTraceProvider(ctx, res)
	if err != nil {
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)

	mp, err := newMeterProvider(ctx, res)
	if err != nil {
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, mp.Shutdown)

	lp, err := newLoggerProvider(ctx, res, logLevel)
	if err != nil {
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, lp.Shutdown)

	return shutdown, err
}

func newTraceProvider(ctx context.Context, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(res))

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
	))

	return tp, nil
}

func newMeterProvider(ctx context.Context, res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter),
		),
		sdkmetric.WithResource(res),
	)

	err = runtimeMetrics.Start(
		runtimeMetrics.WithMeterProvider(mp),
	)
	if err != nil {
		return nil, err
	}

	err = hostMetrics.Start(
		hostMetrics.WithMeterProvider(mp),
	)
	if err != nil {
		return nil, err
	}

	otel.SetMeterProvider(mp)
	return mp, nil
}

func newLoggerProvider(ctx context.Context, res *resource.Resource, logLevel minsev.Severity) (*sdklog.LoggerProvider, error) {
	exporter, err := otlploghttp.New(ctx)
	if err != nil {
		return nil, err
	}

	p := sdklog.NewBatchProcessor(exporter)
	processor := minsev.NewLogProcessor(p, logLevel)
	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
		sdklog.WithResource(res),
	)

	global.SetLoggerProvider(lp)

	return lp, nil
}
