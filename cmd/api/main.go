package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Rq4n/gofollow/cmd/worker"
	"github.com/Rq4n/gofollow/internal/config"
	"github.com/Rq4n/gofollow/internal/database"
	"github.com/Rq4n/gofollow/internal/handler"
	"github.com/Rq4n/gofollow/internal/mailer"
	"github.com/Rq4n/gofollow/internal/repository"
	"github.com/Rq4n/gofollow/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mailClient, err := mailer.NewMailTrapClient(
		cfg.Mailer.FromEmail,
		cfg.Mailer.APIKey,
	)
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/gofollow?sslmode=disable",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.URL,
		cfg.DB.Port,
	)

	dbConfig := &database.Config{
		DSN:          dsn,
		MaxIdleTime:  cfg.DB.MaxIdleTime,
		MaxIdleConns: cfg.DB.MinIdleConn,
		MaxOpenConns: cfg.DB.MaxOpenConn,
	}

	dbPool, err := database.NewConn(dbConfig)
	if err != nil {
		log.Fatalf("database error, %v", err)
	}

	defer dbPool.Close()
	log.Print("connection pool established")

	repo := repository.New(dbPool)

	services := service.NewAppService(repo)
	emailService := service.NewEmailJobService(repo)

	userHandler := handler.NewUserHandler(*services.UserService)
	clientHandler := handler.NewClientHandler(*services.ClientService)

	workerInstance := &worker.Worker{
		Email:  emailService,
		Client: services.ClientService,
		Mail:   mailClient,
	}

	go func() {
		for {
			jobs, err := emailService.GetPendingEmailJobs(context.Background())
			if err != nil {
				log.Println("worker error:", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, job := range jobs {
				go workerInstance.ProcessSingleJob(context.Background(), job.ID)
			}

			time.Sleep(5 * time.Second)
		}
	}()

	app := &APIServer{
		listenAddr:    cfg.Port,
		dbConfig:      dbConfig,
		dbPool:        dbPool,
		userHandler:   userHandler,
		clientHandler: clientHandler,
		mail:          mailClient,
	}

	mux := app.mount()
	log.Fatal(app.start(mux))
}
