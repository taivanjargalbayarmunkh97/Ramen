package main

import (
	"fmt"
	fiberSwagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/wpcodevo/golang-fiber-jwt/controllers/agency"
	"github.com/wpcodevo/golang-fiber-jwt/controllers/auth"
	"github.com/wpcodevo/golang-fiber-jwt/controllers/user"
	"github.com/wpcodevo/golang-fiber-jwt/initializers"
	"github.com/wpcodevo/golang-fiber-jwt/middleware"
	"log"
	"os"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(&config)
}

// @title Ramen API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host https://tranquil-river-84673-f93f313aee4e.herokuapp.com
// @BasePath /api
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

	app.Get("/swagger/*", fiberSwagger.HandlerDefault)
	app.Get("/swagger/*", fiberSwagger.New(fiberSwagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &fiberSwagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))

}
