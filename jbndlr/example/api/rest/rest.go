package rest

import (
	"context"
	"fmt"
	"jbndlr/example/api"
	"jbndlr/example/conf"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// Status : This API's status information
	Status = api.NewStatus()
)

// Serve : Start serving REST API (via HTTP/1.1).
func Serve(port int16) error {
	gin.SetMode("release")
	router := gin.Default()
	router.Use(LimitConcurrentRequests(conf.P.API.LimitMaxConcurrentRequests))
	router.Use(LimitRate(conf.P.API.LimitMaxWindowRequests, time.Duration(conf.P.API.LimitWindowSeconds)*time.Second))

	router.GET("/", handleRoot)
	router.POST("/auth", handleAuth)
	router.GET("/secret", RequireAuthentication(), handleSecret)
	router.GET("/slow", LimitRate(2, 10*time.Second), handleSlow)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  time.Duration(conf.P.API.HTTPReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.P.API.HTTPWriteTimeout) * time.Second,
	}

	start := func() error {
		Status.Start()
		log.Printf("Serving REST :%d\n", port)
		return server.ListenAndServe()
	}

	stop := func() error {
		log.Printf("Stopping REST")
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(conf.P.API.GracefulSeconds)*time.Second,
		)
		defer cancel()
		return server.Shutdown(ctx)
	}

	err := api.ServeGracefully(start, stop)
	Status.Stop(err)
	return err
}

func handleRoot(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusServiceUnavailable, gin.H{
		"error": "Service Unavailable",
	})
}

func handleAuth(ctx *gin.Context) {
	subject, err := Authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.IndentedJSON(http.StatusOK, subject)
}

func handleSecret(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"secret": ":3",
	})
}

func handleSlow(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"speed": "turtle",
	})
}
