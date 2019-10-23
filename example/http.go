package main

import (
	"context"
	"net/http"

	"github.com/phenixrizen/golimiter"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)

	limiter := golimiter.New(1, 2)

	ctx := context.Background()

	// Wrap the servemux with the limit middleware.
	http.ListenAndServe(":4000", limiter.LimitHTTP(ctx, mux))
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
