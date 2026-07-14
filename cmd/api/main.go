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
	configDir := flag.String(`basedir`, `/config/`, `Base dir for configurations`)
	flag.Parse()

	env.LoadEnvironment(*configDir)
}

func main() {
	loglevel, err := logrus.ParseLevel(env.GetString(`LOG_LEVEL`, `info`))
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
		Host:     env.GetString(`DB_HOST`, ``),
		Port:     env.GetString(`DB_PORT`, ``),
		User:     env.GetString(`DB_USER`, ``),
		Password: env.GetString(`DB_PASSWORD`, ``),
		DB:       env.GetString(`DB_NAME`, ``),
	}, logger)
	if err != nil {
		logger.Fatal(`failed to connect to datastore`, err)
	}

	apiHost := env.GetString(`API_HOST`, ``)
	a := api.New(auth, logger, store)
	srv := server.New(context.Background(), a, apiHost)

	logger.Println("starting http listener on", aurora.Cyan(apiHost))
	logger.Printf("allowed origins %s", aurora.Yellow(env.GetString(`ALLOWED_ORIGINS`, `*`)))

	logger.Fatal(srv.Run())
}
