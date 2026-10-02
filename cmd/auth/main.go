package main

import (
	"context"
	"github.com/Holocron1/authservice/configs"
	grpcserver "github.com/Holocron1/authservice/internal/grpc"
	"github.com/Holocron1/authservice/internal/handler"
	"github.com/Holocron1/authservice/internal/service"
	"github.com/Holocron1/authservice/internal/store"
	"github.com/Holocron1/authservice/pkg/authpb"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
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
		log.Fatal("database error: ", err)
	}

	postgresStore := store.NewPostgresStore(pool)
	userService := service.NewUserServiceImpl(postgresStore, c.JWTSecret, c.JWTTTL)

	grpcServer := grpc.NewServer()
	authServer := grpcserver.NewServer(userService)
	authpb.RegisterAuthServiceServer(grpcServer, authServer)

	listener, err := net.Listen("tcp", ":"+c.GRPCPort)
	if err != nil {
		log.Fatal("failed to listen on grpc port: ", err)
	}

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("grpc error:", err)
		}
	}()

	loginHandler := handler.NewLoginHandler(userService)

	http.HandleFunc("/login", loginHandler.Handle)

	verifyHandler := handler.NewVerifyHandler(userService)
	http.HandleFunc("/verify", verifyHandler.Verify)

	err = http.ListenAndServe(":"+c.HTTPPort, nil)
	if err != nil {
		log.Fatal("http server error: ", err)
	}
}
