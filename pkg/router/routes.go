package router

import (
	"github.com/Sourjaya/go-hrms/pkg/controllers"
	"github.com/gofiber/fiber"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/api/v1/employee", controllers.GetAllEmployees)
	app.Get("/api/v1/employee/:id", controllers.GetEmployeeByID)
	app.Post("/api/v1/employee", controllers.NewEmployee)
	// app.Put("api/v1/employee", controllers.UpdateEmployee)
	app.Delete("/api/v1/employee/:id", controllers.DeleteEmployee)
}
