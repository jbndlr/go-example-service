package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LimitConcurrentRequests : Set a hard upper limit for simultaneous requests.
func LimitConcurrentRequests(n int) gin.HandlerFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(c *gin.Context) {
		acquire()       // before request
		defer release() // after request
		c.Next()
	}
}

type window struct {
	n        int
	current  int
	previous int
	baseTime time.Time
	duration time.Duration
}

func (w *window) hit() error {
	if time.Now().Truncate(w.duration) != w.baseTime {
		w.baseTime = time.Now().Truncate(w.duration)
		w.previous = w.current
		w.current = 0
	}

	frac := float64(time.Now().Sub(w.baseTime)) / float64(w.duration)
	rate := (1.0-frac)*float64(w.previous) + float64(w.current)

	if rate >= float64(w.n) {
		return fmt.Errorf("Rate limit exceeded; max %d, was %.2f", w.n, rate)
	}
	w.current++ // Only add this request if it is handled at all.
	return nil
}

// LimitRate : Limit requests per time window (sliding window algorithm).
func LimitRate(n int, duration time.Duration) gin.HandlerFunc {
	wnd := &window{n, 0, 0, time.Now().Truncate(duration), duration}
	return func(c *gin.Context) {
		err := wnd.hit()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too Many Requests",
			})
		}
		c.Next()
	}
}
