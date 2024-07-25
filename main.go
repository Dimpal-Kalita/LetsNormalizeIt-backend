package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/config"
	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/internal/database"
	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/internal/router"
)

func main() {
	cfg := config.Loadconfig()

	dbClient, err := database.NewClient(cfg.MongoURI)
	defer database.CloseClient(dbClient)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB:", err)
	}

	r := router.NewRouter(dbClient, cfg.JWT_SECRET_KEY)

	srv := &http.Server{
		Addr:    ":" + cfg.APP_PORT,
		Handler: r,
	}

	// Gracefull shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :"+cfg.APP_PORT+": %v\n", err)
		}
	}()

	<-stop
	log.Println("Server is Shutting Down... ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited\n")
}
