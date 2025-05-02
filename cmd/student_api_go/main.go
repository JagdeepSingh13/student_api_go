package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JagdeepSingh13/student_api_go/internal/config"
	"github.com/JagdeepSingh13/student_api_go/internal/http/handlers/student"
	"github.com/JagdeepSingh13/student_api_go/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// db set-up
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()

	// so that we can pass this storage, made func CreateStudent in sqlite
	router.HandleFunc("POST /api/students", student.New(storage))

	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server Started:", slog.String("address", cfg.Addr))

	// just to make server shutdown graceful
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e := server.Shutdown(ctx)
	if e != nil {
		slog.Error("failed to shut-down server", slog.String("error", e.Error()))
	}

	slog.Info("server shut-down successfull")

}
