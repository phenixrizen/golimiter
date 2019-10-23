package golimiter

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimitHTTP(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)

	limiter := New(1, 2)

	// wrap the servemux with the limiter middleware.
	go http.ListenAndServe(":42280", limiter.LimitHTTP(mux))

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	tm := time.Now()
	for i := 0; i < 10; i++ {
		resp, err := client.Get("http://localhost:42280/")
		since := time.Since(tm)
		fmt.Println(resp.StatusCode, since)
		if resp.StatusCode == http.StatusOK {
			assert.Nil(err, "error should be nil")
			assert.NotNil(resp, "resp should not be empty")
		} else {
			assert.Nil(err, "error should be nil")
			assert.NotNil(resp, "error response should not be nil")
			assert.Equal(http.StatusTooManyRequests, resp.StatusCode, "error response status code should be 429 Too Many Requests")
		}
	}
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
