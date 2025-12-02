package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"task_manager/controllers"
	"task_manager/data"
	"task_manager/router"
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// MongoDB config
	mongoURI := getenv("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getenv("MONGODB_DB", "task_manager_db")

	if os.Getenv("JWT_SECRET") == "" {
		log.Println("WARNING: JWT_SECRET is not set. Middleware will use a default dev secret.")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("failed to connect to MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("failed to ping MongoDB:", err)
	}

	log.Println("Connected to MongoDB:", mongoURI)

	// Collections
	db := client.Database(dbName)
	taskCol := db.Collection("tasks")
	userCol := db.Collection("users")

	// Services
	taskSvc := data.NewTaskService(taskCol)
	userSvc := data.NewUserService(userCol)

	// Controllers
	taskController := controllers.NewTaskController(taskSvc)
	userController := controllers.NewUserController(userSvc)

	// Router
	r := router.SetupRouter(userController, taskController)

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
