package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	jwt "github.com/andrskhrchk/myapp/pkg/jwt"

	"github.com/andrskhrchk/myapp/internal/repository/postgres"
	"github.com/andrskhrchk/myapp/internal/services/auth"
	"github.com/andrskhrchk/myapp/internal/transport/httpserver"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=govno2 dbname=postgres sslmode=disable"
	db, err := initPostgres(connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	jwtMgr := jwt.NewTokenManager("123")
	authSvc := auth.NewAuthService(*userRepo, jwtMgr)

	handler := httpserver.NewHandler(authSvc, jwtMgr)

	srv := httpserver.NewServer("8080", handler.InitRoutes())

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()
	log.Println("server started: 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

}

func initPostgres(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
