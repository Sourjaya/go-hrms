package main

import (
	"log"

	"github.com/Sourjaya/go-hrms/router"
	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()
	router.SetupRoutes(app)
	log.Println("Server starting at port 8080")
	app.Listen(8080)
}
