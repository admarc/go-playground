package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	db := getDB()

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
		Addr:         ":8083",
		Handler:      r,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  2 * time.Second,
	}

	go func() {
		fmt.Println(s.ListenAndServe())
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	<-ctx.Done()
	fmt.Println("signal received")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	fmt.Println("shutting down")
	s.Shutdown(ctx)
}

func getDB() sql.DB {
	const dbPath = "./users.db"

	if _, err := os.Stat(dbPath); err != nil {
		if _, err = os.Create(dbPath); err != nil {
			panic(err)
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetConnMaxIdleTime(time.Second * 1)

	if err != nil {
		panic(err)
	}

	return *db
}
