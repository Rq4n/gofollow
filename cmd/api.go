package main

import (
	"log"
	"net/http"

	"github.com/Rq4n/gofollow/internal/auth"
	"github.com/Rq4n/gofollow/internal/database"
	"github.com/Rq4n/gofollow/internal/handler"
	"github.com/Rq4n/gofollow/internal/mailer"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServer struct {
	listenAddr    string
	dbConfig      *database.Config
	dbPool        *pgxpool.Pool
	userHandler   *handler.UserHandler
	clientHandler *handler.ClientHandler
	mail          mailer.Mailer
}

func (s *APIServer) mount() http.Handler {
	r := chi.NewRouter()

	r.Post("/register", s.userHandler.HandleCreateUser)
	r.Post("/login", s.userHandler.HandleUserLogin)

	r.Route("/admin", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/client", s.clientHandler.HandleCreateClient)
	})

	return r
}

func (s *APIServer) start(mux http.Handler) error {
	srv := &http.Server{
		Addr:    s.listenAddr,
		Handler: mux,
	}
	log.Printf("server started on port: %v", s.listenAddr)

	return srv.ListenAndServe()
}
