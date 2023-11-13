package v3

import (
	"context"
	"net/http"

	wisdomType "github.com/wisdom-oss/commonTypes"
)

const ERROR_CHANNEL_NAME = "error-channel"
const STATUS_CHANNEL_NAME = "status-channel"

// convertInputChannel converts the bidirectional channel into a send only
// channel
func convertInputChannel(c chan interface{}) chan<- interface{} {
	return c
}

func convertOutputChannel(c chan bool) <-chan bool {
	return c
}

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
			output := make(chan bool)
			// now attach the two channels to the request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, ERROR_CHANNEL_NAME, convertInputChannel(input))
			ctx = context.WithValue(ctx, STATUS_CHANNEL_NAME, convertOutputChannel(output))
			// now use a goroutine to make the error handling code asynchronous
			go func() {
				for {
					select {
					case data := <-input:
						switch data.(type) {
						case string:
							errorCode := data.(string)
							e, errorPresent := errors[errorCode]
							if !errorPresent {
								panic("unregistered error used")
							}
							_ = e.Send(w)
							output <- true
							return
						case error:
							err := data.(error)
							e := wisdomType.WISdoMError{}
							e.WrapError(err)
							_ = e.Send(w)
							output <- true
							return
						default:
							panic("unexpected type found end")
						}
					}
				}
			}()
			// now send the request to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
