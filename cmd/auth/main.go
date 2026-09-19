package main

import (
	"context"
	"github.com/Holocron1/authservice/configs"
	"github.com/Holocron1/authservice/internal/handler"
	"github.com/Holocron1/authservice/internal/service"
	"github.com/Holocron1/authservice/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	c := configs.LoadConfig()

	pool, err := pgxpool.New(context.Background(), c.DatabaseURL)
	if err != nil {
		log.Println("database error", err)
	}

	postgresStore := store.NewPostgresStore(pool)
	userService := service.NewUserServiceImpl(postgresStore, c.JWTSecret, c.JWTTTL)
	loginHandler := handler.NewLoginHandler(userService)

	http.HandleFunc("/login", loginHandler.Handle)

	verifyHandler := handler.NewVerifyHandler(userService)
	http.HandleFunc("/verify", verifyHandler.Verify)

	err = http.ListenAndServe(":"+c.HTTPPort, nil)
	if err != nil {
		log.Println("env file not found")
	}
}
