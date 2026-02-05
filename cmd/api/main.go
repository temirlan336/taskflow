package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"taskflow/internal/api"
	"taskflow/internal/app/server"
	"taskflow/internal/middleware"
	"taskflow/internal/repository/postgres"
	"taskflow/internal/service"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// ---------- PostgreSQL ----------
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://localhost:5432/taskflow?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	ctxPing, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctxPing); err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	fmt.Println("PostgreSQL connected")

	// ---------- Redis ----------
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctxRedis, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctxRedis).Err(); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Redis connected")
	fmt.Println("API_KEY =", os.Getenv("API_KEY"))

	// ---------- Rate Limiter ----------
	limiter := middleware.NewRateLimiter(redisClient, 10*time.Second, 5)

	// ---------- Application layer ----------
	// taskStorage := memory.NewMemoryStorage()
	taskStorage := postgres.NewTaskStorage(db)
	eventStorage := postgres.NewEventStorage(db)
	txManager := postgres.NewTxManagerImpl(db)
	taskService := service.NewTaskService(taskStorage, eventStorage, txManager)
	handler := api.NewHandler(taskService, limiter)

	// ---------- HTTP ----------
	router := http.NewServeMux()

	tasksHandler := http.HandlerFunc(handler.HandleTasks)
	tasksByIDHandler := http.HandlerFunc(handler.HandleTasksByID)

	// router.HandleFunc("/tasks", handler.HandleTasks)
	// router.HandleFunc("/tasks/", handler.HandleTasksByID)

	router.Handle("/tasks", middleware.AuthMiddleware(tasksHandler))
	router.Handle("/tasks/", middleware.AuthMiddleware(tasksByIDHandler))

	loggedRouter := middleware.Logging(router)
	taskFlowServer := server.NewServer(":8080", loggedRouter)

	// ---------- Graceful shutdown ----------
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("server started on 8080")

	go func() {
		if err := taskFlowServer.Run(); err != nil {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	<-signalChan
	log.Println("Shutting down server...")

	ctxShut, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := taskFlowServer.Shutdown(ctxShut); err != nil {
		log.Fatalf("Error during shutdown: %v\n", err)
	}

	log.Println("Server stopped")

}
