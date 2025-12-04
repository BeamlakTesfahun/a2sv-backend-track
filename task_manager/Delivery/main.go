package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	deliveryControllers "task_manager/Delivery/controllers"
	"task_manager/Infrastructure"
	"task_manager/Repositories"
	"task_manager/Usecases"
	"task_manager/Delivery/routers"
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	_ = godotenv.Load()

	mongoURI := getenv("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getenv("MONGODB_DB", "task_manager_db")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("mongo connect error:", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("mongo ping error:", err)
	}

	db := client.Database(dbName)
	taskCol := db.Collection("tasks")
	userCol := db.Collection("users")

	// Repositories
	taskRepo := Repositories.NewMongoTaskRepository(taskCol)
	userRepo := Repositories.NewMongoUserRepository(userCol)

	// Infrastructure services
	passSvc := Infrastructure.NewBcryptPasswordService()
	jwtSvc := Infrastructure.NewJWTService()

	// Usecases
	taskUC := Usecases.NewTaskUsecases(taskRepo)
	userUC := Usecases.NewUserUsecases(userRepo, passSvc, jwtSvc)

	// Controllers
	ctrl := deliveryControllers.NewController(taskUC, userUC)

	// Router
	r := routers.SetupRouter(ctrl, jwtSvc)

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("server error:", err)
	}
}
