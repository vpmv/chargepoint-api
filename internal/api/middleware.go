package api

import (
	"github.com/go-fuego/fuego"
	"github.com/vpmv/chargepoint-api/pkg/authenticator"

	"net/http"
)

func (api *API) bearerAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, ok, err := api.auth.AuthenticateBearer(fuego.TokenFromHeader(r))

		if err != nil {
			fuego.SendJSONError(w, r, new(fuego.InternalServerError))
			return
		}

		if !ok {
			fuego.SendJSONError(w, r, new(fuego.UnauthorizedError))
			return
		}

		next.ServeHTTP(w, r.WithContext(authenticator.SetAuthorizationContext(r.Context(), *auth)))
	})
}

func (api *API) hasPermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := authenticator.GetAuthorization(r.Context())

			if !auth.HasPermission(permission) {
				fuego.SendError(w, r, new(fuego.UnauthorizedError))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
