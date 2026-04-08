package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rq4n/gofollow/internal/auth"
	"github.com/Rq4n/gofollow/internal/database"
	"github.com/Rq4n/gofollow/internal/handler"
	"github.com/Rq4n/gofollow/internal/mailer"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Application struct {
	listenAddr string
	mail       mailer.Mailer
	logger     *zap.SugaredLogger
	DBConfig
	Handler
}

type Handler struct {
	handleUser   *handler.UserHandler
	handleClient *handler.ClientHandler
}

type DBConfig struct {
	dbPool   *pgxpool.Pool
	DBConfig *database.Config
}

func (app *Application) mount() http.Handler {
	r := chi.NewRouter()

	r.Post("/register", app.handleUser.HandleCreateUser)
	r.Post("/login", app.handleUser.HandleUserLogin)

	r.Route("/admin", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/client", app.handleClient.HandleCreateClient)
	})

	return r
}

func (app *Application) start(mux http.Handler) error {
	srv := &http.Server{
		Addr:    app.listenAddr,
		Handler: mux,
	}
	log.Printf("server started on port: %v", app.listenAddr)

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Print("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	log.Print("server has stopped")

	return nil
}
