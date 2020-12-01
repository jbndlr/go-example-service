package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Status : Struct carrying information on API status.
type Status struct {
	Serving     bool
	Error       error
	TimeInit    time.Time
	TimeStarted time.Time
	TimeStopped time.Time

	Ready   func() bool
	Healthy func() bool
}

func alwaysReady() bool {
	return true
}

func alwaysHealthy() bool {
	return true
}

// NewStatus : Construct a new (initial) API status.
func NewStatus() *Status {
	return &Status{
		false, nil, time.Now(), time.Time{}, time.Time{},
		alwaysReady, alwaysHealthy,
	}
}

// Format : Format Status fields as strings in gin.H output.
func (as *Status) Format() gin.H {
	timeOrEmpty := func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format(time.RFC3339)
	}
	errorOrEmpty := func(e error) string {
		if e == nil {
			return ""
		}
		return fmt.Sprintf("%v", e)
	}

	return gin.H{
		"serving":         fmt.Sprintf("%t", as.Serving),
		"error":           errorOrEmpty(as.Error),
		"timeInitialized": timeOrEmpty(as.TimeInit),
		"timeStarted":     timeOrEmpty(as.TimeStarted),
		"timeStopped":     timeOrEmpty(as.TimeStopped),
		"ready":           fmt.Sprintf("%t", as.Ready()),
		"healthy":         fmt.Sprintf("%t", as.Healthy()),
	}
}

// Stop : Reflect a stop.
func (as *Status) Stop(e error) {
	as.Serving = false
	as.Error = e
	as.TimeStopped = time.Now()
}

// Start : Reflect a successful start.
func (as *Status) Start() {
	as.Serving = true
	as.Error = nil
	as.TimeStarted = time.Now()
	as.TimeStopped = time.Time{}
}
