package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZXstrike/marketplace-app/internal/config"
	"github.com/ZXstrike/marketplace-app/internal/database"
	"github.com/ZXstrike/marketplace-app/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {

	// LoadConfig()
	config, er := config.LoadConfig()
	if er != nil {
		log.Fatalf("Error loading config: %v", er)
	}

	// InitializeDatabase()
	db, er := database.PostgresConnect(&config.PostgresConfig)
	if er != nil {
		log.Fatalf("Error connecting to database: %v", er)
	}

	// StartServer()
	StartServer(config.ServerPort, db, config.PrivateKey, config.PublicKey)
}

// StartServer starts the server
func StartServer(Port string, db *gorm.DB, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {
	router := gin.New()

	gin.SetMode(gin.DebugMode)

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://app.test", "https://zxsttms.tech", "http://zxsttms.tech"}, // The origin of your Vue app
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Add middlewares
	router.Use(gin.Recovery(), gin.Logger(), cors.New(config))

	// Initialize app routes
	routes.InitRoutes(router, db, privateKey, publicKey)

	// Create an HTTP server using the Gin router
	srv := &http.Server{
		Addr:    ":" + Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("Starting server on port %s\n", Port)
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
