package v2

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	wisdomType "github.com/wisdom-oss/commonTypes"
	utils "github.com/wisdom-oss/microservice-utils"
)

var ErrorMissingUserInformation = wisdomType.WISdoMError{
	ErrorCode:        "MISSING_AUTHORIZATION_INFORMATION_USER",
	ErrorTitle:       "Missing appropriate authorization information",
	ErrorDescription: "The request is missing the required authorization information about the user",
	HttpStatusCode:   400,
	HttpStatusText:   http.StatusText(400),
}

var ErrorMissingGroupsInformation = wisdomType.WISdoMError{
	ErrorCode:        "MISSING_AUTHORIZATION_INFORMATION_GROUPS",
	ErrorTitle:       "Missing appropriate authorization information",
	ErrorDescription: "The request is missing the required authorization information about the groups of the user",
	HttpStatusCode:   400,
	HttpStatusText:   http.StatusText(400),
}

var ErrorNotInRequiredGroup = wisdomType.WISdoMError{
	ErrorCode:        "USER_NOT_IN_REQUIRED_GROUP",
	ErrorTitle:       "Forbidden",
	ErrorDescription: "The supplied user is not in the required user group to access this service",
	HttpStatusCode:   403,
	HttpStatusText:   http.StatusText(403),
}

func Authorization(c wisdomType.AuthorizationConfiguration, serviceName string) func(http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// create a context which contains the following introspection
				// results
				ctx := context.Background()

				// first add the status of the authorization to the context
				ctx = context.WithValue(ctx, "authEnabled", c.Enabled)

				// before trying any other authorization measures, check if
				// the authorization is enabled and bypass it accordingly
				if !c.Enabled {
					// let the next handler handle the request since the
					// authorization is disabled
					nextHandler.ServeHTTP(w, r)
					return
				}

				// now get the values for the authorization information header
				rawIsAdmin := strings.TrimSpace(r.Header.Get("X-Superuser"))
				rawGroups := strings.TrimSpace(r.Header.Get("X-Authenticated-Groups"))
				user := strings.TrimSpace(r.Header.Get("X-Authenticated-User"))

				// now try to parse the admin header
				isAdmin, err := strconv.ParseBool(rawIsAdmin)
				if err != nil {
					// if the value is not parsable, then default to a non-admin
					// user
					isAdmin = false
				}
				// append this result to the context
				ctx = context.WithValue(ctx, "isAdmin", isAdmin)

				// now try to parse the groups that the user is a member of
				groups := strings.Split(rawGroups, ",")
				// append the collected groups to the context of the request
				ctx = context.WithValue(ctx, "groups", groups)

				// now append the user to the request context
				ctx = context.WithValue(ctx, "user", user)

				// first check if the service requires a username to process the
				// request and a username is set/not empty
				if c.RequireUserIdentification && user == "" {
					// since the user is empty (therefore, not set) reject the
					// request with an error message
					_ = ErrorMissingUserInformation.Send(w)
					return
				}

				// now check if the user is an administrator and therefore is
				// allowed to bypass the group restrictions
				if isAdmin {
					nextHandler.ServeHTTP(w, r)
					return
				}

				// now check if the group list contains any entries if not, the
				// request is lacking this kind of information and is rejected
				if len(groups) == 0 {
					// since the groups are empty (therefore, not set) reject
					// the request with an error message
					_ = ErrorMissingGroupsInformation.Send(w)
					return
				}

				// now check if the group list contains the group required by
				// the configuration
				if utils.ArrayContains(groups, c.RequiredUserGroup) {
					nextHandler.ServeHTTP(w, r)
					return
				}

				// since the group list did not contain the required group,
				// the access is forbidden.
				_ = ErrorNotInRequiredGroup.Send(w)
				return
			},
		)
	}
}
