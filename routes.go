package main

import "github.com/gofiber/fiber/v2"

// Sets up the API endpoints
func setupAPIRoutes(app *fiber.App) {
	app.Get("/faq", getFAQ)
	app.Post("/faq", addFAQ)
	app.Delete("/faq/:id", deleteFAQ)
	app.Patch("/faq/:id", patchFAQ)
}
