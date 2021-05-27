package main

import (
	"go-restful/api"
	"go-restful/app"
	"go-restful/conf"
	"go-restful/service"

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

	// API Routes.
	api.Load(router)

	// init service
	svc := service.New()

	// set global service
	service.Svc = svc

	app.Run()
}
