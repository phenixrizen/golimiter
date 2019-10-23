# golimiter

[![Build Status](https://travis-ci.org/phenixrizen/golimiter.svg?branch=master)](https://travis-ci.org/phenixrizen/golimiter) [![GoDoc](https://godoc.org/github.com/phenixrizen/golimiter?status.svg)](https://godoc.org/github.com/phenixrizen/golimiter) [![Coverage Status](https://coveralls.io/repos/github/phenixrizen/golimiter/badge.svg?branch=master)](https://coveralls.io/github/phenixrizen/golimiter?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/phenixrizen/golimiter)](https://goreportcard.com/report/github.com/phenixrizen/golimiter)

Provides a nice limiting API for Golang.

### Get Started

#### Installation
```bash
$ go get github.com/phenixrizen/golimiter
```

#### Usage

```go
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
```