package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwt"

	wisdomType "github.com/wisdom-oss/commonTypes/v2"
)

const issuer = "api-gateway"

func Authorization(config wisdomType.AuthorizationConfiguration, serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get the current context from the request and set the state of the
			// authorization middleware to the context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "auth.enabled", config.Enabled)

			// check if the authorization has been disabled
			if !config.Enabled {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// now check and get the value of the "Authorization" header
			authHeaders, authHeaderSet := r.Header["Authorization"]
			if !authHeaderSet {
				_ = ErrMissingAuthorizationHeader.Send(w)
				return
			}

			// use the first header sent, ignore the others
			authHeader := strings.TrimSpace(authHeaders[0])

			// check if the authorization header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				_ = ErrUnsupportedTokenScheme.Send(w)
				return
			}

			// extract the token from the authorization header
			rawToken := strings.TrimPrefix(authHeader, "Bearer ")

			// now check the jwt to be sure that the token is still alive and
			// does not contain any errors
			serviceToken, err := jwt.ParseString(rawToken, jwt.WithValidate(true), jwt.WithIssuer(issuer))
			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrInvalidJWT()):
					_ = ErrJWTMalformed.Send(w)
					return
				case errors.Is(err, jwt.ErrTokenExpired()):
					_ = ErrJWTExpired.Send(w)
					return
				case errors.Is(err, jwt.ErrTokenNotYetValid()):
					_ = ErrJWTNotYetValid.Send(w)
					return
				case errors.Is(err, jwt.ErrInvalidIssuedAt()):
					_ = ErrJWTNotCreatedYet.Send(w)
					return
				case errors.Is(err, jwt.ErrInvalidIssuer()):
					_ = ErrJWTInvalidIssuer.Send(w)
					return
				default:
					e := wisdomType.WISdoMError{}
					e.WrapNativeError(err)
					_ = e.Send(w)
					return
				}
			}

			groups, groupsSet := serviceToken.PrivateClaims()["groups"].([]string)
			if !groupsSet {
				_ = ErrJWTNoGroups.Send(w)
				return
			}

			allowUser := false
			for _, group := range groups {
				if group == serviceName {
					allowUser = true
					break
				}
			}

			// now check if the user is an administrator and may bypass the
			// group restrictions
			staff, staffSet := serviceToken.PrivateClaims()["staff"].(string)
			if staffSet && staff == "true" {
				ctx = context.WithValue(ctx, "auth.admin", true)
				ctx = context.WithValue(ctx, "auth.group", serviceName)
				allowUser = true
			}

			if !allowUser {
				_ = Forbidden.Send(w)
				return
			}

			// since the user is allowed to access the resource the request will
			// be sent to the next handler with some additional context
			ctx = context.WithValue(ctx, "auth.admin", false)
			ctx = context.WithValue(ctx, "auth.group", serviceName)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
