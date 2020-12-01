package esid

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"jbndlr/example/api"
	"jbndlr/example/api/grpc"
	"jbndlr/example/api/rest"
	"jbndlr/example/conf"

	"github.com/gin-gonic/gin"
)

var (
	// Status : This API's status information
	Status = api.NewStatus()
)

// Serve : Start serving ESID API (Exposed Status Information & Diagnostics).
func Serve(port int16) error {
	gin.SetMode("release")
	router := gin.Default()
	router.GET("/healthy", handleHealthy)
	router.GET("/ready", handleReady)
	router.GET("/status", handleStatus)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  time.Duration(conf.P.API.HTTPReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.P.API.HTTPWriteTimeout) * time.Second,
	}

	start := func() error {
		Status.Start()
		log.Printf("Serving ESID :%d\n", port)
		return server.ListenAndServe()
	}

	stop := func() error {
		log.Printf("Stopping ESID")
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

func allHealthy() bool {
	return Status.Healthy() && grpc.Status.Healthy() && rest.Status.Healthy()
}

func allReady() bool {
	return Status.Ready() && grpc.Status.Ready() && rest.Status.Ready()
}

func handleHealthy(ctx *gin.Context) {
	statusCode := func() int {
		if allHealthy() {
			return http.StatusOK
		}
		return http.StatusInternalServerError
	}()
	ctx.IndentedJSON(statusCode, gin.H{
		"_self":   "healthy",
		"healthy": fmt.Sprintf("%t", allHealthy()),
	})
}

func handleReady(ctx *gin.Context) {
	statusCode := func() int {
		if allReady() {
			return http.StatusOK
		}
		return http.StatusInternalServerError
	}()
	ctx.IndentedJSON(statusCode, gin.H{
		"_self": "ready",
		"ready": fmt.Sprintf("%t", allReady()),
	})
}

func handleStatus(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"_self": "status",
		"service": gin.H{
			"name":    conf.P.Service.Name,
			"version": conf.P.Service.Version,
		},
		"healthy": fmt.Sprintf("%t", allHealthy()),
		"ready":   fmt.Sprintf("%t", allReady()),
		"api": gin.H{
			"esid": Status.Format(),
			"rest": rest.Status.Format(),
			"grpc": grpc.Status.Format(),
		},
	})
}
