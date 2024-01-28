package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AvivKermann/BlogAgg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB   *database.Queries
	port string
}

func main() {
	godotenv.Load()

	// set up for database connection
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("cannot connect to database")
	}
	dbQuries := database.New(db)

	// set up for cfg variable
	cfg := apiConfig{
		DB:   dbQuries,
		port: os.Getenv("PORT"),
	}

	//all the different routers
	router := chi.NewRouter()
	router.Use(cors.Handler(middleware))
	v1Router := chi.NewRouter()

	// all the different mounts
	router.Mount("/v1", v1Router)

	// all the different handlers
	v1Router.HandleFunc("/readiness", handlerReadiness)
	v1Router.HandleFunc("/err", handlerErr)
	v1Router.Get("/users", cfg.handlerGetUser)
	v1Router.Post("/users", cfg.handlerCreateUser)
	server := &http.Server{
		Addr:    ":" + cfg.port,
		Handler: router,
	}

	fmt.Printf("server running on localhost:%v\n", cfg.port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("server crashed")
	}
}
