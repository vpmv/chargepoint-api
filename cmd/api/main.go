package main

import (
	"context"
	"flag"

	"github.com/logrusorgru/aurora"

	"github.com/sirupsen/logrus"
	"github.com/vpmv/chargepoint-api/internal/api"
	"github.com/vpmv/chargepoint-api/internal/server"
	"github.com/vpmv/chargepoint-api/internal/storage/postgres"
	"github.com/vpmv/chargepoint-api/pkg/authenticator"
	env "github.com/vpmv/goenv"
)

var (
	seedDatabase *bool
	logger       *logrus.Logger
)

type SimpleAuthenticator struct {
	tokens map[string]*authenticator.Authorization
}

func (auth *SimpleAuthenticator) AuthenticateBearer(apiKey string) (*authenticator.Authorization, bool, error) {
	if env.IsEnv(`development`) && apiKey == `` {
		logger.Debug(aurora.Red(`No Authorization token supplied. Did you forget to prefix the token with Bearer?`))
	}
	if app, ok := auth.tokens[apiKey]; ok && app.Token == apiKey {
		return app, true, nil
	}

	return nil, false, nil
}

func init() {
	configDir := flag.String(`basedir`, `/config/`, `Base dir for configurations`)
	seedDatabase = flag.Bool(`seed`, false, `Seed database with test data`)
	flag.Parse()

	env.LoadDotEnv(*configDir)
}

func main() {
	loglevel, err := logrus.ParseLevel(env.GetString(`LOG_LEVEL`, `info`))
	if err != nil {
		panic(`invalid log level: ` + err.Error())
	}

	logger = logrus.New()
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
		Host:     env.MustString(`DB_HOST`),
		Port:     env.MustString(`DB_PORT`),
		User:     env.MustString(`DB_USER`),
		Password: env.MustString(`DB_PASSWORD`),
		DB:       env.MustString(`DB_NAME`),
	}, logger)
	if err != nil {
		logger.Fatal(`failed to connect to datastore`, err)
	}

	if seedDatabase != nil && *seedDatabase {
		store.MustSeed()
	}

	apiHost := env.GetString(`API_HOST`, ``)
	a := api.New(auth, logger, store)
	srv := server.New(context.Background(), a, apiHost)

	logger.Println("starting http listener on", aurora.Cyan(apiHost))
	logger.Printf("allowed origins %s", aurora.Yellow(env.GetString(`ALLOWED_ORIGINS`, `*`)))

	logger.Fatal(srv.Run())
}
