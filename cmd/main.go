package main

import (
	"github.com/caarlos0/env"
	sdformatter "github.com/joonix/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"taxes-be/utils/echoaddons"
	logrusctx "taxes-be/utils/logrus"
)

type Config struct {
	FacebookKey string
	LogLevel    string `env:"LOG_LEVEL" envDefault:"debug"`
}

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(sdformatter.NewFormatter())
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(sdformatter.NewFormatter())
	logrus.AddHook(&logrusctx.LogCtxHook{})

	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		logrus.WithError(err).Fatalln("failed to read config, exiting")
	}

	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err == nil {
		logrus.SetLevel(lvl)
	} else {
		logrus.WithError(err).Fatalln("failed to parse log level from cfg")
	}

	// httpClient := client.NewHttpClient()

	e := echo.New()
	e.Validator = echoaddons.NewValidator()
	e.HTTPErrorHandler = echoaddons.CustomHTTPErrorHandler

	e.Use(middleware.CORS())

	http.Handle("/", e)
	err = http.ListenAndServe(":8080", http.DefaultServeMux)
	if err != nil {
		logrus.WithError(err).Fatalln("could not start server")
	}

}
