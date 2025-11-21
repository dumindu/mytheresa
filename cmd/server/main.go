package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/mytheresa/go-hiring-challenge/internal/config"
	"github.com/mytheresa/go-hiring-challenge/internal/router"
	"github.com/mytheresa/go-hiring-challenge/internal/util/logger"
	"github.com/mytheresa/go-hiring-challenge/internal/util/pgutil"
	"github.com/mytheresa/go-hiring-challenge/internal/util/validatorutil"
)

//	@title			MYTHERESA.DEV API
//	@version		1.0
//	@description	This is a sample RESTful API application

//	@contact.name	Dumindu Madunuwan
//	@contact.url	https://www.linkedin.com/in/dumindunuwan

//	@license.name	MIT License
//	@license.url	https://github.com/dumindu/mytheresa/blob/main/LICENSE

// @servers.url	localhost:8080/v1
func main() {
	// Initialize configuration, logger, and validator
	c := config.New()
	l := logger.New(c.Server.Debug)
	v := validatorutil.New()

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	db, close := pgutil.New(&c.DB)
	defer close()

	// Initialize router
	r := router.New(db, l, v)

	// Set up the HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	// Start the server
	go func() {
		log.Printf("Starting server on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}

		log.Println("Server stopped gracefully")
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
	srv.Shutdown(ctx)
	stop()
}
