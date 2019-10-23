package main

import (
	"net/http"

	"github.com/phenixrizen/golimiter"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)

	limiter := golimiter.New(1, 2)

	// Wrap the servemux with the limit middleware.
	http.ListenAndServe(":42280", limiter.LimitHTTP(mux))
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
