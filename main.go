package main

import (
	"go-restful/api"
	"go-restful/app"
	"go-restful/conf"
	"go-restful/service"

	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/spf13/pflag"
)

var (
	cfg = pflag.StringP("config", "c", "", "config file path.")
)

func main() {
	pflag.Parse()

	// init config, parse yaml
	if err := conf.Init(*cfg); err != nil {
		panic(err)
	}

	// init app
	app := app.New()

	// Initial the Gin engine.
	router := app.Router

	// Health Check
	router.GET("/health", api.Ping)

	// get global Monitor object
	m := ginmetrics.GetMonitor()
	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(router)

	// API Routes.
	api.Load(router)

	// init service
	svc := service.New()

	// set global service
	service.Svc = svc

	app.Run()
}
