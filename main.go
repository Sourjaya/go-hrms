package main

import (
	"log"

	"github.com/Sourjaya/go-hrms/pkg/router"
	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()
	router.SetupRoutes(app)
	log.Println("Server starting at port 8080")
	log.Fatal(app.Listen(3000))
}
