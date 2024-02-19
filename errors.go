package middleware

import wisdomType "github.com/wisdom-oss/commonTypes/v2"

// InvalidTypeProvided represents an internal error which is sent if an invalid
// type as been provided to the input channel of the error handler
var InvalidTypeProvided = wisdomType.WISdoMError{
	Type:   "https://pkg.go.dev/github.com/wisdom-oss/microservice-middlewares/v4#InvalidTypeProvided",
	Status: 500,
	Title:  "Invalid Type Provided",
	Detail: "An invalid type has been provided to the error handler. Please contact your administrator",
}

// ErrMissingAuthorizationHeader is returned if the request did not contain
// the `Authorization` header
var ErrMissingAuthorizationHeader = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc6750.html#section-3.1",
	Status: 401,
	Title:  "Missing Authorization Header",
	Detail: "The request did not contain the 'Authorization' header. Please check your request.",
}

// ErrUnsupportedTokenScheme is returned if the request did not utilize the
// Bearer token scheme as documented in [RFC 6750].
//
// [RFC 6750]: https://www.rfc-editor.org/rfc/rfc6750
var ErrUnsupportedTokenScheme = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc6750.html#section-3.1",
	Status: 400,
	Title:  "Unsupported Token Scheme used",
	Detail: "The token scheme used in this request is not supported by the service. Please check your request.",
}

// ErrJWTMalformed is returned if the request did contain a JWT but is malformed
var ErrJWTMalformed = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "JSON Web Token Malformed",
	Detail: "The JSON Web Token presented as Bearer Token is not correctly formatted",
}

// ErrJWTExpired is returned if the JWT in the request is already expired
var ErrJWTExpired = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc6750.html#section-3.1",
	Status: 401,
	Title:  "JSON Web Token Expired",
	Detail: "The JSON Web Token used to access this resource has expired. Access has been denied",
}

// ErrJWTNotYetValid is returned if the field indicating a time before the
// jwt is not valid contains a time in the future
var ErrJWTNotYetValid = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc6750.html#section-3.1",
	Status: 401,
	Title:  "JSON Web Token Used Before Validity",
	Detail: "The JSON Web Token used to access this resource has been used before it is permitted to be used. Access has been denied",
}

// ErrJWTNotCreatedYet is returned if the JWTs iat field indicating at which the
// token has been issued is in the future
var ErrJWTNotCreatedYet = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc6750.html#section-3.1",
	Status: 401,
	Title:  "JSON Web Token Used Before Creation",
	Detail: "The JSON Web Token used to access this resource been created in the future, therefore it is invalid and the access has been denied. Please check your authentication provider.",
}

// ErrJWTInvalidIssuer is returned if the JWTs issuer field indicates that it
// has not been issued by the API Gateway
var ErrJWTInvalidIssuer = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc6750.html#section-3.1",
	Status: 401,
	Title:  "JSON Web Token Issuer Wrong",
	Detail: "The JSON Web Token used to access this resource has not been issued by the correct issuer. Please check your authentication provider.",
}

// ErrJWTNoGroups is returned if the JWT did not contain the group claim and
// therefore is not usable for the service
var ErrJWTNoGroups = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "JSON Web Token No Groups Claim",
	Detail: "The JSON Web Token used to access this resource did not contain the required `groups` claim",
}

// Forbidden is returned if the user is authenticated but not authorized to
// access the resource
var Forbidden = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.4",
	Status: 403,
	Title:  "Access Forbidden",
	Detail: "The user is not in the appropriate user group to access this service",
}
