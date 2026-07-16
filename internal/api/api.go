package api

import (
	"context"
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	chiCORS "github.com/go-chi/cors"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
	"github.com/sirupsen/logrus"
	"github.com/vpmv/chargepoint-api/internal/storage"
	"github.com/vpmv/chargepoint-api/pkg/authenticator"
	env "github.com/vpmv/chargepoint-api/pkg/dotenv"
)

type API struct {
	auth     authenticator.Authenticator
	log      *logrus.Logger
	store    storage.ChargePointClient
	Security fuego.Security
}

func New(auth authenticator.Authenticator, logger *logrus.Logger, store storage.ChargePointClient) *API {
	return &API{
		auth:  auth,
		store: store,
		log:   logger,
	}
}

// Bootstrap of the API
func (api *API) Bootstrap(server *fuego.Server, ctx context.Context) {
	api.log.Debug(`Migrating DB...`)
	if err := api.store.Migrate(); err != nil {
		api.log.Fatal(`Error migrating models`)
	}

	api.log.Debug(`Setting up router...`)

	// register CORS
	cors := chiCORS.New(chiCORS.Options{
		// AllowedOrigins:   strings.Split(os.Getenv(`ALLOWED_ORIGINS`), `;`),
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowedHeaders:   []string{"Cookie", "Accept", "Authorization", "Content-Type", "Origin"},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           300,
	})
	fuego.Use(server, cors.Handler)
	// register throttling
	fuego.Use(server, chiMiddleware.Throttle(env.GetInt(`THROTTLE_THRESHOLD`, 100)))

	// Health check
	fuego.Get(server, "/_health", api.HealthCheck)

	v1 := fuego.Group(server, `/v1`, option.Header(`Authorization`, `Usage: "Bearer \<token\>"`, param.Required()))
	fuego.Use(v1, api.bearerAuthorization)

	// /v1/chargepoints
	chargePointGroup := fuego.Group(v1, `/chargepoints`)
	fuego.Use(
		chargePointGroup,
		api.hasPermission(PermissionReadCP),
	)

	fuego.Get(chargePointGroup, ``, api.ListChargePoints,
		option.Middleware(api.hasPermission(PermissionReadCP)),
		optionPagination,
		option.Summary(`List all charge points`),
		option.Description(`Returns a list of charge points, within the paginated range.<br> Deleted charge points are not included`),
	)
	fuego.Post(chargePointGroup, ``, api.CreateChargePoints,
		option.Middleware(api.hasPermission(PermissionCreateCP)),
		option.Summary(`Create new charge points`),
		option.Description(`Creates and/or updates charge points<br> Previously deleted charge points will be undeleted.`),
	)
	fuego.Get(chargePointGroup, `/{id}`, api.GetChargePointByID,
		option.Middleware(api.hasPermission(PermissionReadCP)),
		option.Path(`id`, `Vendor ID`, param.Required()),
		option.Summary(`Get a charge point by VendorID`),
		option.Description(`Retrieves a charge point by VendorID`),
	)
	fuego.Delete(chargePointGroup, `/{id}`, api.DeleteChargePoint,
		option.Middleware(api.hasPermission(PermissionCreateCP)),
		option.Summary(`Delete a charge point`),
		option.Description(`Deletes a charge point by VendorID<br> Deletions aren't permanent'`),
	)
	fuego.Get(chargePointGroup, `/nearby`, api.ListChargePointsByLocation,
		option.Middleware(api.hasPermission(PermissionReadCP)),
		option.QueryInt(`radius`, `radius in KM`, param.Required()),
		option.Query(`lat`, `latitude`, param.Required()),
		option.Query(`lon`, `longitude`, param.Required()),
		option.Summary(`List nearby charge points`),
		option.Description(`Retrieves a list of all ChargePoints within a radius (kilometer) from lat/long coordinate`),
	)

	api.log.Debug(`API bootstrapped!`)
}

// HealthCheck is used to do server health checks outside normal api routing
func (api *API) HealthCheck(c fuego.ContextWithBody[struct{}]) (string, error) {
	return `healthy`, nil
}

func (api *API) QueryParamUint(c fuego.ContextNoBody, param string) uint {
	i := c.QueryParamInt(param)
	return uint(i)
}
