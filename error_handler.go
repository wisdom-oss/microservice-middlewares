//go:generate gomarkdoc --output v4.md .
package middleware

import (
	"context"
	"net/http"

	wisdomType "github.com/wisdom-oss/commonTypes/v2"
)

const ERROR_CHANNEL_NAME = "error-channel"
const STATUS_CHANNEL_NAME = "status-channel"

// ErrorHandler allows the global handling and wrapping errors
// occurring in API calls. The function needs the service name as a parameter
// to correctly generate the error code used in the wisdomType.WISdoMError.
// Furthermore, it also accepts the usage of preregistered errors
//
// To access the channel added to the request context in an http handler use
// the following call:
//
//	errorHandler := r.Context().Value("error-channel").(chan<- interface{})
//
// To watch for the handling to be completed, use the following channel from
// the handler
//
//	errorHandled :=  r.Context().Value("status.channel").(<-chan bool)
//
// To handle an error just send it into the error handler channel and listen on
// the statusChannel for a boolean return.
//
//	errorHandler <- errors.New("test error")
//	<-errorHandled
//
// After handling the error it is recommended to exit the handler to hide errors
// and warnings about the http response writer being called when closed or tying
// to write headers again
func ErrorHandler(serviceName string, errors map[string]wisdomType.WISdoMError) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// create a channel for receiving errors and strings by using
			// interface as type
			input := make(chan interface{})
			// create a channel for notifying the sender that the error has
			// been handled
			statusChannel := make(chan bool)
			// now attach the two channels to the request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, ERROR_CHANNEL_NAME, prepareInputChannel(input))
			ctx = context.WithValue(ctx, STATUS_CHANNEL_NAME, prepareStatusChannel(statusChannel))
			// now use a goroutine to make the error handling code asynchronous
			go func() {
				for {
					select {
					case data := <-input:
						switch data.(type) {
						case wisdomType.WISdoMError:
							e := data.(wisdomType.WISdoMError)
							_ = e.Send(w)
							statusChannel <- true
							return
						case error:
							err := data.(error)
							e := wisdomType.WISdoMError{}
							e.WrapNativeError(err)
							_ = e.Send(w)
							statusChannel <- true
							return
						default:
							_ = InvalidTypeProvided.Send(w)
							statusChannel <- true
							return
						}
					}
				}
			}()
			// now send the request to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
