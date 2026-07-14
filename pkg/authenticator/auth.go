package authenticator

import (
	"context"
)

const AuthContextKey = `Authorization`

type Authenticator interface {
	AuthenticateBearer(bearer string) (*Authorization, bool, error)
}

// GetAuthorization retrieves authorization DTO from context
func GetAuthorization(ctx context.Context) Authorization {
	return ctx.Value(AuthContextKey).(Authorization)
}

// SetAuthorizationContext sets authorization DTO
func SetAuthorizationContext(ctx context.Context, auth Authorization) context.Context {
	return context.WithValue(ctx, AuthContextKey, auth)
}
