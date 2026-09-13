package main

import (
	"AuthService/configs"
	"AuthService/internal/handler"
	"AuthService/internal/service"
	"AuthService/internal/store"
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

	postgresStore := store.NewPostgresStore(c)
	userService := service.UserServiceImpl{UserStore: postgresStore}
	loginHandler := &handler.LoginHandler{UserService: &userService}

	http.HandleFunc("/login", loginHandler.Handle)

	verifyHandler := &handler.VerifyHandler{UserService: &userService}
	http.HandleFunc("/verify", verifyHandler.Verify)

	err = http.ListenAndServe(":"+c.HTTP_PORT, nil)
	if err != nil {
		log.Println("env file not found")
	}
}
