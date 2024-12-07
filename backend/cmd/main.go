package main

import (
	"auth/backend/internal/config"
	"auth/backend/internal/database"
	"auth/backend/internal/handlers"
	"auth/backend/internal/middleware"
	"auth/backend/internal/services"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetCfg()

	storage := database.NewStorage(cfg)
	err := storage.DB.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully connected to database")

	tokenRepo := database.NewTokenRepo(storage.DB)

	s := services.NewService(tokenRepo, cfg)

	h := handlers.NewHandler(s)

	r := h.NewRouter()
	r.Use(middleware.ContentTypeApplicationJsonMiddleware)

	srv := &http.Server{
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	log.Println("Server started at", srv.Addr)
}
