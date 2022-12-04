package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/admarc/users/cmd/server/config"
	"github.com/admarc/users/internal/handlers"
	storageUsers "github.com/admarc/users/internal/storage/users"
	"github.com/admarc/users/internal/users"
	"github.com/admarc/users/pkg/dbcollector"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	fmt.Println("starting")
	defer fmt.Println("shutdown")

	cfg, help, err := config.New()

	if err != nil {
		if help != "" {
			log.Fatal(help)
		}
		log.Fatal(err)
	}

	db := getDB(cfg)

	repo := storageUsers.NewStorage(&db)
	us := users.NewService(repo)
	uh := handlers.NewUsers(us)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", uh.Create)
		r.Get("/{id}", uh.Get)
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(15 * time.Second)
		})
	})

	prometheus.MustRegister(dbcollector.NewSQLDatabaseCollector("general", "main", "sqlite", &db))
	r.Mount("/metrics", promhttp.Handler())

	s := http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      r,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	go func() {
		fmt.Println(s.ListenAndServe())
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	<-ctx.Done()
	fmt.Println("signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulTimeout)
	defer cancel()

	fmt.Println("shutting down")
	s.Shutdown(ctx)
}

func getDB(cfg config.Config) sql.DB {
	db, err := sql.Open("sqlite3", cfg.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdleTime)

	return *db
}
