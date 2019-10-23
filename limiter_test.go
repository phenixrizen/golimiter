package golimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimiter(t *testing.T) {
	assert := assert.New(t)

	limiter := New(1, 1)

	ch := make(chan int, 5000)

	tm := time.Now()
	go fun(limiter *Limiter) {
		for i := 0; i < 5000; i++ {
		time.Sleep(50 * time.Millisecond)
		limiter.Wait(nil)
		ch <- i
		}
	}(limiter)

	assert.Nil(nil)
}
