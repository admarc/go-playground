package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/admarc/users/internal/handlers"
	storageUsers "github.com/admarc/users/internal/storage/users"
	"github.com/admarc/users/internal/users"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/mattn/go-sqlite3"
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
	r.Post("/users", uh.Create)

	s := http.Server{
		Addr:    ":8083",
		Handler: r,
	}

	fmt.Println(s.ListenAndServe())
}

func getDB() sql.DB {
	const dbPath = "./users.db"

	if _, err := os.Stat(dbPath); err != nil {
		if _, err = os.Create(dbPath); err != nil {
			panic(err)
		}
	}

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		panic(err)
	}

	createDB := "CREATE TABLE IF NOT EXISTS `users` (`id` VARCHAR(36) PRIMARY KEY,`name` VARCHAR(64) NULL);"

	_, err = db.Exec(createDB)

	if err != nil {
		panic(err)
	}

	return *db
}
