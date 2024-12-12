package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Vansh3140/GOVault/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Creating a new Fiber app instance
	app := fiber.New()

	// Defining routes and grouping API routes under /api/govault
	app.Route("/api/govault", func(api fiber.Router) {
		// Route to get all resources in a collection
		api.Get("/:collection", routes.GetAll)
		// Route to get a specific resource in a collection
		api.Get("/:collection/:resource", routes.GetOne)
		// Route to create a new resource in a collection
		api.Post("/:collection", routes.CreateOne)
		// Route to update a specific resource in a collection
		api.Put("/:collection/:resource", routes.UpdateOne)
		// Route to delete all resources in a collection
		api.Delete("/:collection", routes.DeleteAll)
		// Route to delete a specific resource in a collection
		api.Delete("/:collection/:resource", routes.DeleteOne)
	})

	// Channel to listen for OS signals like SIGINT, SIGTERM, etc.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)

	// Running the Fiber app in a separate goroutine to prevent blocking
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Blocking the main goroutine when OS signal is received
	<-stop
	log.Println("Received shutdown signal, shutting down...")

	// Gracefully shutdown the Fiber app
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server shutdown successfully")
}
