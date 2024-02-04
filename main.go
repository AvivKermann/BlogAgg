package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	v1Router.Get("/users", cfg.middlewareAuth(cfg.handlerGetUser))
	v1Router.Get("/feeds", cfg.handlerGetAllFeeds)
	v1Router.Post("/users", cfg.handlerCreateUser)
	v1Router.Post("/feeds", cfg.middlewareAuth(cfg.handlerCreateFeed))
	v1Router.Post("/feed_follows", cfg.middlewareAuth(cfg.handlerCreateFeedFollow))
	v1Router.Get("/posts", cfg.middlewareAuth(cfg.handlerGetPosts))
	v1Router.Get("/feed_follows", cfg.middlewareAuth(cfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", cfg.handlerDeleteFeedFollow)

	server := &http.Server{
		Addr:    ":" + cfg.port,
		Handler: router,
	}

	fmt.Printf("server running on localhost:%v\n", cfg.port)
	go startScraping(dbQuries, 10, time.Minute)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("server crashed")
	}
}
