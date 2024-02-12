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
