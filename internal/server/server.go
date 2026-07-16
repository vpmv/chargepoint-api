package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/getkin/kin-openapi/openapi3"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-fuego/fuego"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/vpmv/chargepoint-api/internal/api"
	env "github.com/vpmv/goenv"
)

func Recover() func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stackTrace := debug.Stack()

					logEntry := chiMiddleware.GetLogEntry(r)
					if logEntry != nil {
						logEntry.Panic(rec, stackTrace)
					} else {
						_, _ = fmt.Fprintf(os.Stderr, "Panic: %+v\n", rec)
						debug.PrintStack()
					}

					errorText := http.StatusText(http.StatusInternalServerError)

					if os.Getenv(`ENV`) == `development` {
						errorText = fmt.Sprintf("%s\n %+v\n", errorText, rec)
					}

					http.Error(w, errorText, http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func openApiHandler(specURL string) http.Handler {
	return httpSwagger.Handler(
		httpSwagger.Layout(httpSwagger.BaseLayout),
		httpSwagger.PersistAuthorization(true),
		httpSwagger.URL(specURL), // The url pointing to API definition
	)
}

func New(ctx context.Context, api *api.API, hostAddr string, options ...func(*fuego.Server)) *fuego.Server {
	serverOptions := []func(*fuego.Server){
		fuego.WithAddr(hostAddr),

		// disabled AutoAuth in lieu of custom API Bearer tokens
		// enabling this would generate an /auth/ group with user/password login
		// returning secure signed JWT
		// fuego.WithAutoAuth(api.Oauth),

		fuego.WithRouteOptions(
			fuego.OptionAddResponse(http.StatusForbidden, "Forbidden", fuego.Response{Type: fuego.HTTPError{}}),
		),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				SwaggerURL: `/openapi`,
				SpecURL:    `/openapi/openapi.json`,
				Info: &openapi3.Info{
					Title:       `ChargePoint API`,
					Description: `Demo application`,
					Contact: &openapi3.Contact{
						Name:  `Valentijn V.`,
						URL:   `https://valentijn.co`,
						Email: `contact@valentijn.co`,
					},

					Version: "v1",
				},
				// UIHandler: openApiHandler,
			}),
		),
	}

	options = append(serverOptions, options...)

	app := fuego.NewServer(options...)
	api.Security = app.Security

	if openapiHost := env.GetString(`OPENAPI_HOST`, ``); len(openapiHost) > 0 {
		app.OpenAPI.Description().Servers = append(app.OpenAPI.Description().Servers, &openapi3.Server{URL: openapiHost})
	}

	fuego.Use(app, Recover())
	fuego.Use(app, chiMiddleware.Compress(5, "application/json"))

	api.Bootstrap(fuego.Group(app, "/api"), ctx)

	return app
}
