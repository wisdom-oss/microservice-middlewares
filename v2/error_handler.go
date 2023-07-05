package v2

import (
	"context"
	"net/http"

	wisdomType "github.com/wisdom-oss/commonTypes"
)

// NativeErrorHandler allows the global handling and wrapping of native errors
// occurring in API calls. The function needs the service name as a parameter
// to correctly generate the error code used in the wisdomType.WISdoMError
//
// To access the channel added to the request context in a http handler use
// the following call:
//
//	nativeErrorChannel := r.Context().Value("nativeErrorChannel").(chan error)
//
// To render an error, just send it into the channel using the following syntax:
//
//	nativeErrorChannel<-err
//
// Due to the way channels work in golang to stop the handler from returning too
// early and triggering errors with the transmitted content length, the
// NativeErrorHandler adds a boolean channel called "nativeErrorHandled". This
// channel needs to be listened to using the following syntax to stop the
// handler sending the error from returning too early
//
//		nativeErrorHandled := r.Context().Value("nativeErrorHandled").(chan bool)
//	 <-nativeErrorHandled
func NativeErrorHandler(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// create a new channel
			c := make(chan error)
			b := make(chan bool)
			// now access the request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "nativeErrorChannel", c)
			ctx = context.WithValue(ctx, "nativeErrorHandled", b)
			// use a go function to listen to the channel and output the
			// request error to the client using json
			go func() {
				for {
					select {
					case err := <-c:
						e := wisdomType.WISdoMError{}
						e.WrapError(err, serviceName)
						_ = e.Send(w)
						b <- true
						return
					}
				}
			}()
			// now let the next handler handle the request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// WISdoMErrorHandler allows the global handling and wrapping of WISdoM errors
// occurring in API calls. The function needs the mapping of predefined error
// messages using the wisdomType.WISdoMError
//
// To access the channel added to the request context in a http handler use
// the following call:
//
//	wisdomErrorChannel := r.Context().Value("wisdomErrorChannel").(chan string)
//
// To render the error just send the error code into the channel:
//
//	wisdomErrorChannel <- "test"
//
// If an unregistered error code is used a panic will be released in the go
// function handling the actual rendering
//
// Due to the way channels work in golang to stop the handler from returning too
// early and triggering errors with the transmitted content length, the
// WISdoMErrorHandler adds a boolean channel called "wisdomErrorHandled". This
// channel needs to be listened to using the following syntax to stop the
// handler sending the error from returning too early
//
//		wisdomErrorHandled := r.Context().Value("wisdomErrorChannel").(chan bool)
//	 <-wisdomErrorHandled
func WISdoMErrorHandler(errors map[string]wisdomType.WISdoMError) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// create a new channel
			c := make(chan string)
			b := make(chan bool)
			// now access the request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "wisdomErrorChannel", c)
			ctx = context.WithValue(ctx, "wisdomErrorHandled", b)
			// use a go function to listen to the channel and output the
			// request error to the client using json
			go func() {
				for {
					select {
					case errorCode := <-c:
						e, errorPresent := errors[errorCode]
						if !errorPresent {
							panic("using unregistered error")
						}
						_ = e.Send(w)
						return
					}
				}
			}()
			// now let the next handler handle the request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
