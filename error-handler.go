//go:generate gomarkdoc --output v4.md .
package middleware

import (
	"context"
	"net/http"

	wisdomType "github.com/wisdom-oss/commonTypes/v2"
)

const (
	ErrorChannelName = string(rune(iota))
	StatusChannelName
)

// ErrorHandler is a function that returns an http.Handler middleware function
// for handling errors. It creates channels for receiving errors and notifying
// the sender that the error has been handled. It attaches these channels to
// the request context and uses a goroutine to handle errors asynchronously.
// The middleware function sends the request to the next handler after
// attaching the channels to the request context.
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a channel for receiving errors and strings by using
		// interface as type
		input := make(chan interface{})
		// create a channel for notifying the sender that the error has
		// been handled
		statusChannel := make(chan bool)
		// now attach the two channels to the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, ErrorChannelName, prepareInputChannel(input))
		ctx = context.WithValue(ctx, StatusChannelName, prepareStatusChannel(statusChannel))
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
