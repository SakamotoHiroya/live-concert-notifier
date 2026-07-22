package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/handler"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	databaseURL := os.Getenv("DATABASE_URL")

	ctx := context.Background()
	pool, err := repository.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	artistRepo := repository.NewArtistRepository(pool)
	concertRepo := repository.NewConcertRepository(pool)

	userService := service.NewUserService(userRepo)
	artistService := service.NewArtistService(artistRepo)
	followService := service.NewFollowService(repository.NewFollowRepository(pool), userRepo, artistRepo)
	concertService := service.NewConcertService(concertRepo)
	dashboardService := service.NewDashboardService(concertRepo, userRepo)
	adminScrapeService := service.NewAdminScrapeService(repository.NewScrapeJobRepository(pool), artistRepo, service.LogScraperTrigger{})
	apiHandler := handler.NewAPIHandler(userService, artistService, followService, concertService, dashboardService, adminScrapeService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, apiHandler)

	log.Printf("server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
