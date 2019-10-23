package golimiter

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimitHTTP(t *testing.T) {
	assert := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)

	limiter := New(1, 2)
	t.Log("Limiting requests to 1/sec, with burst of 2")

	// wrap the servemux with the limiter middleware.
	go http.ListenAndServe(":42280", limiter.LimitHTTP(mux))

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	tm := time.Now()
	for i := 0; i < 10; i++ {
		resp, err := client.Get("http://localhost:42280/")
		assert.NotNil(resp, "resp should not be empty")
		t.Log("Code:", resp.StatusCode, "Retry-After:", resp.Header.Get(RetryAfterHeader), "Runtime:", time.Since(tm))
		if resp.StatusCode == http.StatusOK {
			assert.Nil(err, "error should be nil")
		} else {
			assert.Nil(err, "error should be nil")
			assert.Equal(http.StatusTooManyRequests, resp.StatusCode, "error response status code should be 429 Too Many Requests")
		}
	}

	err := limiter.Wait(context.Background())
	assert.Nil(err)
	time.Sleep(1 * time.Second)

	t.Log("Testing Retry-After")
	tm = time.Now()
	retryAfter := 0
	for i := 0; i < 10; i++ {
		resp, err := client.Get("http://localhost:42280/")
		assert.NotNil(resp, "resp should not be empty")
		t.Log("Code:", resp.StatusCode, "Retry-After:", resp.Header.Get(RetryAfterHeader), "Runtime:", time.Since(tm))
		if resp.StatusCode == http.StatusOK {
			assert.Nil(err, "error should be nil")
		} else {
			retryAfter, err = strconv.Atoi(resp.Header.Get(RetryAfterHeader))
			assert.Nil(err, "could not cast retry after to int")
			break
		}
	}
	t.Log("Retry-After:", retryAfter)
	time.Sleep(time.Duration(retryAfter) * time.Second)
	resp, err := client.Get("http://localhost:42280/")
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Nil(err, "error should be nil")
	assert.NotNil(resp, "resp should not be empty")
	t.Log(resp.StatusCode, time.Since(tm))
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
