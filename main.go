package main

import (
	"context"
	"fmt"
	_ "github.com/busragumusel/insider-case/docs" // Import Swagger Docs
	"github.com/busragumusel/insider-case/internal/api"
	"github.com/busragumusel/insider-case/internal/entity"
	"github.com/busragumusel/insider-case/internal/repository"
	"github.com/busragumusel/insider-case/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var db *gorm.DB

func initDB() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Istanbul",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&entity.Message{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Connected to PostgreSQL and migrated schema")
}

func initRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	options := &redis.Options{
		Addr: redisAddr,
		DB:   0,
	}

	if redisPassword != "" {
		options.Password = redisPassword
	}

	client := redis.NewClient(options)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	log.Println("Redis is connected!")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	initDB()
	initRedis()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	messageRouter := api.NewAPI(db)
	messageRouter.RegisterRoutes(router)

	messageRepo := repository.NewMessageRepository(db)
	stopChan := make(chan bool)
	messageService := service.NewMessageService(messageRepo, stopChan)
	go messageService.StartProcess()

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
	fmt.Println("Server is running on port 8080")
}
