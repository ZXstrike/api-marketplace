package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZXstrike/internal/routers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	startServer("8080", nil) // Replace nil with actual DB connection if needed
}

func startServer(port string, db *gorm.DB) {
	router := gin.New()

	// Add middlewares
	router.Use(gin.Recovery(), gin.Logger())

	// Initialize app routes
	routers.InitRoutes(router, db, nil) // Replace nil with actual Redis client if needed

	// Create an HTTP server using the Gin router
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	// Context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// <-ctx.Done()
	// fmt.Println("Timeout of 5 seconds.")

	fmt.Println("Server exited")
}
