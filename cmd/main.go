package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Rq4n/gofollow/cmd/worker"
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

	cfg := &APIServer{
		listenAddr: os.Getenv("PORT"),
	}

	dbConfig := &database.Config{
		DSN:          os.Getenv("DSN"),
		MaxIdleTime:  os.Getenv("MAX_IDLE_TIME"),
		MaxIdleConns: 15,
		MaxOpenConns: 15,
	}

	dbPool, err := database.NewConn(dbConfig)
	if err != nil {
		log.Fatalf("database error, %v", err)
	}

	defer dbPool.Close()
	log.Print("connection pool established")

	mailClient, err := mailer.NewMailTrapClient(
		os.Getenv("MAILTRAP_FROM_EMAIL"),
		os.Getenv("MAILTRAP_API_KEY"),
	)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.New(dbPool)

	// initialize user service and handler
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	// initialize client service and handler
	clientService := service.NewClientService(repo)
	clientHandler := handler.NewClientHandler(clientService)

	// initialize email service and handler
	emailService := service.NewEmailJobService(repo)

	workerInstance := &worker.Worker{
		Email:  emailService,
		Client: clientService,
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
		listenAddr:    cfg.listenAddr,
		dbConfig:      cfg.dbConfig,
		dbPool:        dbPool,
		userHandler:   userHandler,
		clientHandler: clientHandler,
		mail:          mailClient,
	}

	mux := app.mount()
	log.Fatal(app.start(mux))
}
