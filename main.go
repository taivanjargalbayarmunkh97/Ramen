package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/wpcodevo/golang-fiber-jwt/controllers/agency"
	"github.com/wpcodevo/golang-fiber-jwt/controllers/auth"
	"github.com/wpcodevo/golang-fiber-jwt/controllers/user"
	"github.com/wpcodevo/golang-fiber-jwt/initializers"
	"github.com/wpcodevo/golang-fiber-jwt/middleware"
	"log"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(&config)
}

func main() {
	app := fiber.New()
	micro := fiber.New()
	app.Mount("/api", micro)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST",
		AllowCredentials: true,
	}))

	// Нэвтрэх, бүртгүүлэх, гарах
	micro.Route("/auth", func(router fiber.Router) {
		router.Post("/register", auth.SignUpUser)
		router.Post("/login", auth.SignInUser)
		router.Get("/logout", middleware.DeserializeUser, auth.LogoutUser)
	})

	// Үндсэн хэрэглэгчийн мэдээлэл
	micro.Route("/users", func(router fiber.Router) {
		micro.Get("/me", middleware.DeserializeUser, user.GetMe)
	})

	// Агент
	micro.Route("/agent", func(router fiber.Router) {
		router.Post("/list", middleware.DeserializeUser, agency.GetAgentList)
		router.Get("/:id", middleware.DeserializeUser, agency.GetAgent)
		router.Post("/create", middleware.DeserializeUser, agency.CreateAgency)
		router.Put("/:id", middleware.DeserializeUser, agency.UpdateAgent)
		router.Delete("/:id", middleware.DeserializeUser, agency.DeleteUser)
	})

	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "JWT Authentication with Golang, Fiber, and GORM",
		})
	})

	micro.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exists on this server", path),
		})
	})

	port := "54018"
	log.Fatal(app.Listen(":" + port))

}
