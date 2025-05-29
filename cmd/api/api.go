package main

import (
	"Url-Shortener/internal/auth"
	"Url-Shortener/internal/ratelimiter"
	"Url-Shortener/internal/store"
	"Url-Shortener/internal/store/cache"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	authenticator auth.Authenticator
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	auth        authConfig
	redisCfg    redisConfig
	rateLimiter ratelimiter.Config
}

type dbConfig struct {
	host     string
	port     string
	dbName   string
	username string
	password string
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	user string
	pass string
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

func (app *application) mount() http.Handler {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Optional Rate Limiting Middleware (if enabled in config)
	if app.config.rateLimiter.Enabled {
		e.Use(app.RateLimiterMiddleware)
	}

	// -----------------------------
	// API V1 Routes
	// -----------------------------
	v1 := e.Group("api/v1")

	// Health check endpoint (no auth)
	v1.GET("/health", app.healthCheckHandler)

	// URL resource routes (with token auth)
	url := v1.Group("/urls")
	url.POST("", app.createUrlHandler, app.AuthTokenMiddleware()) // Create short URL
	url.GET("/:shortCode", app.getUrlHandler)                     // Resolve short URL

	// -----------------------------
	// Authentication Routes
	// -----------------------------
	authentication := v1.Group("/auth")

	// Register and login
	authentication.POST("/login", app.loginUserHandler)
	authentication.POST("/register", app.registerUserHandle)

	return e
}

func (app *application) run(mux http.Handler) error {

	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go func() {
		// Create context that listens for the interrupt signal from the OS.
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		// Listen for the interrupt signal.
		<-ctx.Done()

		log.Println("shutting down gracefully, press Ctrl+C again to force")

		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown with error: %v", err)
		}

		log.Println("Server exiting")

		// Notify the main goroutine that the shutdown is complete
		done <- true
	}()

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %s", err))
	}
	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")

	return nil
}
