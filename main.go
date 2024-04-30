package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	InitDB()
	LoadTranslations()

	if err := createBackup(); err != nil {
		fmt.Println("Error creating backup:", err)
		return
	}

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "fiber",
		EnableIPValidation: true,
		EnableTrustedProxyCheck: true,
		AppName: "BestHome API v1.0.0",
	})
		port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST, DELETE, PATCH",
		AllowHeaders:     "Content-Type , x-custom-lang",
		MaxAge: 		 3600,
	}))

	// Set up middleware
	app.Use(logRequests)

	// Define API endpoints
	setupAPIRoutes(app)

	// Print a message to the console
	log.Printf("Server is running on port %s\n", port)

	app.Listen(":" + port)

}
