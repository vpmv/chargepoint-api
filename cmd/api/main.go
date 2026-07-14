package main

import (
	"context"
	"flag"
	"os"

	"github.com/logrusorgru/aurora"

	"github.com/sirupsen/logrus"
	"github.com/vpmv/chargepoint-api/internal/api"
	"github.com/vpmv/chargepoint-api/internal/server"
	"github.com/vpmv/chargepoint-api/internal/storage/postgres"
	"github.com/vpmv/chargepoint-api/pkg/authenticator"
	env "github.com/vpmv/chargepoint-api/pkg/dotenv"
)

type Config struct {
	Env      string
	LogLevel string
	DB       *postgres.Config
}

type SimpleAuthenticator struct {
	tokens map[string]*authenticator.Authorization
}

func (auth *SimpleAuthenticator) AuthenticateBearer(apiKey string) (*authenticator.Authorization, bool, error) {
	// todo: remove this block
	if env.IsEnv(`development`) {
		return auth.tokens[`secret`], true, nil
	}

	if app, ok := auth.tokens[apiKey]; ok && app.Token == apiKey {
		return app, true, nil
	}

	return nil, false, nil
}

func init() {
	configDir := flag.String(`basedir`, `/config`, `Base dir for configurations`)
	flag.Parse()

	env.LoadEnvironment(*configDir)
}

func main() {
	loglevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(`invalid log level: ` + err.Error())
	}

	logger := logrus.New()
	logger.SetLevel(loglevel)

	// Create the authenticator
	auth := &SimpleAuthenticator{
		tokens: map[string]*authenticator.Authorization{
			"secret": {
				Token:       "secret",
				Description: "Admin",
				Permissions: []authenticator.Permission{
					{api.PermissionCreateCP, "Create Charge Points"},
					{api.PermissionReadCP, "Read Charge Points"},
				},
			},
		},
	}

	store, err := postgres.NewClient(postgres.Config{
		Host:     os.Getenv(`DB_HOST`),
		Port:     os.Getenv(`DB_PORT`),
		User:     os.Getenv(`DB_USER`),
		Password: os.Getenv(`DB_PASSWORD`),
		DB:       os.Getenv(`DB_NAME`),
	}, logger)
	if err != nil {
		logger.Fatal(`failed to connect to datastore`, err)
	}

	a := api.New(auth, logger, store)
	srv := server.New(context.Background(), a, os.Getenv(`API_HOST`))

	logger.Println("starting http listener on", aurora.Cyan(os.Getenv(`API_HOST`)))
	logger.Printf("allowed origins %s", aurora.Yellow(os.Getenv(`ALLOWED_ORIGINS)`)))

	logger.Fatal(srv.Run())
}
