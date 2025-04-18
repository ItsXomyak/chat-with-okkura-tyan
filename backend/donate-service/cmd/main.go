package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"donate-service/internal/handler"
	"donate-service/internal/repository"
)

func main() {
	app := fiber.New()

	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}

	app.Use(cors.New(cors.Config{
	AllowOrigins: "*",
	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	AllowMethods: "GET, POST, PUT, DELETE, OPTIONS, PATCH",
}))


	handler.RegisterRoutes(app, db)

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
