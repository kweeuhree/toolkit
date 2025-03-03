package toolkit

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (t *Tools) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if t.InfoLog != nil {
			t.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method,
				r.URL.RequestURI()) // Use provided logger
		} else {
			log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method,
				r.URL.RequestURI()) // Fallback to default log package
		}

		next.ServeHTTP(w, r)
	})
}

func (t *Tools) GetClientIP(r *http.Request) string {
	// Check for the X-Forwarded-For header
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// If multiple IPs are present, split by comma and return the first part
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// If X-Forwarded-For is not present, return RemoteAddr
	// This will return the IP address of the immediate connection,
	// which might be the client or a proxy
	return strings.Split(r.RemoteAddr, ":")[0]
}

func (t *Tools) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				t.ServerError(w, fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
