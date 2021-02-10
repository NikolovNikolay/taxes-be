package main

import (
	"database/sql"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/caarlos0/env"
	sdformatter "github.com/joonix/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"taxes-be/internal/atleastonce"
	"taxes-be/internal/atleastonce/atleastoncedao"
	"taxes-be/internal/cron"
	"taxes-be/internal/statements/statementsview"
	"taxes-be/utils/echoaddons"
	logrusctx "taxes-be/utils/logrus"
	"time"
)

type Config struct {
	FacebookKey string
	LogLevel    string `env:"LOG_LEVEL" envDefault:"debug"`

	DB               string `env:"DB_CONNECT_STRING" envDefault:"user=postgres password=password dbname=postgres sslmode=disable""` // nolint:lll // let it stay on one line
	UseInternalTimer bool   `env:"USE_INTERNAL_TIMER" envDefault:"true"`
	AloLimit         int    `env:"ALO_LIMIT" envDefault:"15"`

	AWSRegion              string `env:"AWS_REGION" envDefault:"eu-west-1"`
	AWSAccessKey           string `env:"AWS_ACCESS_KEY" envDefault:""`
	AWSSecretKey           string `env:"AWS_SECRET_KEY" envDefault:""`
	S3StatementsBucketName string `env:"S3_STATEMENTS_BUCKET_NAME" envDefault:""`
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

	db := openConnection(cfg.DB)
	defer func() {
		_ = db.Close()
	}()
	logrus.Debug("connected to DB")

	aloStore := atleastoncedao.NewStore(db)
	aloDoer := atleastonce.New(aloStore)
	defer aloDoer.Close()

	awsSession := connectAws(cfg)

	e := echo.New()
	e.Validator = echoaddons.NewValidator()
	e.HTTPErrorHandler = echoaddons.CustomHTTPErrorHandler
	e.Use(middleware.CORS())

	statementsview.RegisterStatementsEndpoints(e.Group("/api/statements"), awsSession, cfg.S3StatementsBucketName)

	http.Handle("/", e)

	if cfg.UseInternalTimer {
		logrus.Info("starting internal tickers...")
		go cron.RunInternalTimer(aloDoer, cfg.AloLimit)
	}

	srv := &http.Server{
		Addr:              ":8080",
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		Handler: &echoaddons.MaxBytesHandler{
			Handler:      http.DefaultServeMux,
			MaxReqLength: 5 << 20, // 5 MB
		},
	}

	err = srv.ListenAndServe()
	if err != nil {
		logrus.WithError(err).Fatalln("could not start server")
	}
}

func openConnection(connectionStr string) *sql.DB {
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"driver":     "postgres",
			"connection": connectionStr,
		}).WithError(err).Fatalln("failed to connect to the database")
	}

	return db
}

func connectAws(cfg *Config) *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(cfg.AWSRegion),
			Credentials: credentials.NewStaticCredentials(
				cfg.AWSAccessKey,
				cfg.AWSSecretKey,
				"",
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}