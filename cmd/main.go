package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ps_backend/api"
	"ps_backend/db"
)

func main() {
	// Load config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8100"
	}
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost user=psdb password=jmung002 dbname=ps_db port=5432 sslmode=disable"
	}
	corsOrigins := os.Getenv("CORS_ORIGINS") // comma-separated

	// Initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	// Connect database and run migrations
	if _, err := db.ConnectAndMigrateWithDSN(dsn); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.LoggerWithWriter(log.Writer()), gin.Recovery())

	// CORS
	origins := strings.Split(corsOrigins, ",")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup API routes
	router = api.SetupRouter()

	// HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Infof("Server running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Info("Server exiting")
}
