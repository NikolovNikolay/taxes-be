package main

import (
	"database/sql"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/caarlos0/env"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	sdformatter "github.com/joonix/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v71"
	"net/http"
	_ "net/http/pprof"
	"os"
	"taxes-be/internal/atleastonce"
	"taxes-be/internal/atleastonce/atleastoncedao"
	aws2 "taxes-be/internal/aws"
	"taxes-be/internal/conversion"
	"taxes-be/internal/coupons/couponsdao"
	"taxes-be/internal/cron"
	"taxes-be/internal/inquiries/inquiriesdao"
	"taxes-be/internal/payments/paymentsview"
	"taxes-be/internal/sendgrid"
	"taxes-be/internal/statements"
	"taxes-be/internal/statements/statementsview"
	"taxes-be/utils/echoaddons"
	logrusctx "taxes-be/utils/logrus"
	"time"
)

type Config struct {
	FacebookKey string
	LogLevel    string `env:"LOG_LEVEL" envDefault:"debug"`

	DB               string `env:"DB_CONNECT_STRING" envDefault:"user=postgres password=password dbname=postgres sslmode=disable"` // nolint:lll // let it stay on one line
	UseInternalTimer bool   `env:"USE_INTERNAL_TIMER" envDefault:"true"`
	AloLimit         int    `env:"ALO_LIMIT" envDefault:"15"`
	ServerAddress    string `env:"SERVER_ADDRESS" envDefault:":8080"`
	WebsiteBaseURL   string `env:"WEBSITE_BASE_URL" envDefault:"http://localhost:3000/#"`

	ClientToken string `env:"CLIENT_TOKEN" envDefault:""`

	SendgridAPIKey string `env:"SENDGRID_API_KEY" envDefault:""`
	StripeAPIKey   string `env:"STRIPE_API_KEY" envDefault:""`

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

	stripe.Key = cfg.StripeAPIKey

	db := openConnection(cfg.DB)
	defer func() {
		_ = db.Close()
	}()
	logrus.Debug("connected to DB")

	inquiryStore := inquiriesdao.NewStore(db)
	couponStore := couponsdao.NewStore(db)

	aloStore := atleastoncedao.NewStore(db)
	aloDoer := atleastonce.New(aloStore)
	defer aloDoer.Close()

	awsSession := connectAws(cfg)
	s3Manager := aws2.NewS3Manager(awsSession)
	mailer := sendgrid.NewMailer(cfg.SendgridAPIKey)

	rs := conversion.NewExchangeRateService()
	statements.NewStatementManager(
		*aloDoer,
		s3Manager,
		cfg.S3StatementsBucketName,
		mailer,
		inquiryStore,
		couponStore,
		aloStore,
		rs,
	)

	rlm := tollbooth.NewLimiter(3, &limiter.ExpirableOptions{
		DefaultExpirationTTL: 1 * time.Second,
		ExpireJobInterval:    1 * time.Minute,
	})

	e := echo.New()
	e.Validator = echoaddons.NewValidator()
	e.HTTPErrorHandler = echoaddons.CustomHTTPErrorHandler
	e.Use(middleware.CORS())

	statementsview.RegisterStatementsEndpoints(
		e.Group("/api/statements",
			echoaddons.AuthHandler(cfg.ClientToken),
			echoaddons.RateLimitHandler(rlm),
		),
		s3Manager,
		cfg.S3StatementsBucketName,
		inquiryStore,
		couponStore,
		aloStore,
	)

	paymentsview.RegisterEndpoints(
		e.Group("/api/payments",
			echoaddons.AuthHandler(cfg.ClientToken),
			echoaddons.RateLimitHandler(rlm),
		),
		cfg.WebsiteBaseURL,
		couponStore,
		inquiryStore,
	)

	e.GET("/", func(context echo.Context) error {
		return context.NoContent(http.StatusOK)
	})

	http.Handle("/", e)

	if cfg.UseInternalTimer {
		logrus.Info("starting internal tickers...")
		go cron.RunInternalTimer(aloDoer, cfg.AloLimit)
	}

	srv := &http.Server{
		Addr:              cfg.ServerAddress,
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
