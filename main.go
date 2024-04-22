package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	port := ":8080"
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	r := gin.Default()
	r.LoadHTMLGlob("pages/*")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home", gin.H{
			"title": "test",
		})
	})
	r.GET("/about", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "about", gin.H{
			"title": "test",
		})
	})

	r.GET("/exit", func(ctx *gin.Context) {
		quit <- syscall.SIGINT
	})

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		log.Println("Listening and serving on port: ", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	// will block the app until receives a signal from the O.S.
	<-quit

	// cancel will release all resources associated with the context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	// gracefully stop accepting new requests and waits for the active ones to be handled
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}
