package http

import (
	"net/http"
)

// A Handler responds to an HTTP request.
//
// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it
// is not valid to use the ResponseWriter or read from the
// Request.Body after or concurrently with the completion of the
// ServeHTTP call.
//
// Depending on the HTTP client software, HTTP protocol version, and
// any intermediaries between the client and the Go server, it may not
// be possible to read from the Request.Body after writing to the
// ResponseWriter. Cautious handlers should read the Request.Body
// first, and then reply.
//
// Except for reading the body, handlers should not modify the
// provided Request.
//
// If ServeHTTP panics, the server (the caller of ServeHTTP) assumes
// that the effect of the panic was isolated to the active request.
// It recovers the panic, logs a stack trace to the server error log,
// and either closes the network connection or sends an HTTP/2
// RST_STREAM, depending on the HTTP protocol. To abort a handler so
// the client sees an interrupted response but the server doesn't log
// an error, panic with the value ErrAbortHandler.

// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
// type HandlerFunc func(ResponseWriter, *Request)

/*
	This function uses the http.Handler wrapper technique. Functions that take an
	http.handler and return a new one can do things before and/or after the handler
	is called, and even decide whether to call the original handler at all.

	An http.Handler wrapper is a function that has one input argument, both of type
	http.Handler.

	Wrappers have the following signature:
		func(http.Handler) http.Handler
	The idea is that you take in an http.Handler and return a new one that does
	something else before and/or after calling the ServeHTTP method on the
	original.

	For example, a simple logging wrapper might look like this:
	func log(h http.Handler) http.Handler {
  		return http.HandlerFunc(func(w http.ResponseWriter, r
                                          *http.Request) {
    		log.Println("Before")
    		h.ServeHTTP(w, r) // call original
    		log.Println("After")
  		})
	}

	Here, our log function returns a new handler (remember that http.HandlerFunc is a valid http.Handler too)
	that will print the “Before” string, and call the original handler before printing out the “After” string.
	Now, wherever I pass my original http.Handler I can wrap it such that:
		http.Handle("/path", handleThing)
		becomes
		http.Handle("/path", log(handleThing))
	When would you use wrappers?
	This approach can be used to address lots of different situations, including but not limited to:
	- Logging and tracing
	- Validating the request; such as checking authentication credentials
	- Writing common response headers

	Wrappers get to decide whether to call the original handler or not. If they want to, they can
	even intercept the request and response on their own. Say a key URL parameter is mandatory in our API:

		func checkAPIKey(h http.Handler) http.Handler {
  			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    			if len(r.URL.Query().Get("key") == 0) {
      				http.Error(w, "missing key", http.StatusUnauthorized)
      				return // don't call original handler
    			}
    			h.ServeHTTP(w, r)
  			})
		}

	The checkAPIKey wrapper will make sure there is a key, and if there isn’t, it will return with an Unauthorized error.
	It could be extended to validate the key in a datastore, or ensure the caller is within acceptable rate limits etc.

	Deferring
Using Go’s defer statement we can add code that we can be sure will run whatever happens inside our original handler:
func log(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Println("Before")
    defer log.Println("After")
    h.ServeHTTP(w, r)
  })
}
Now, even if the code inside our handler panics, we’ll still see the “After” line printed.
*/

func CorsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("access-control-allow-origin", "*")
			if r.Method == http.MethodOptions {
				w.Header().Set("access-control-allow-headers", "authorization")
				w.Header().Set("access-control-allow-methods", "PATCH,PUT,POST,OPTIONS,GET,DELETE")
				return
			}
			h.ServeHTTP(w, r)
		})
}
