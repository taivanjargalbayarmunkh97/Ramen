package main

import (
	"example.com/ramen/controllers/agency"
	"example.com/ramen/controllers/auth"
	"example.com/ramen/controllers/file"
	"example.com/ramen/controllers/role"
	"example.com/ramen/controllers/user"
	_ "example.com/ramen/docs"
	initializers2 "example.com/ramen/initializers"
	"example.com/ramen/middleware"
	"fmt"
	fiberSwagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
)

func init() {
	config, err := initializers2.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers2.ConnectDB(&config)
}

// @title Ramen API
// @version 1.0
// @description This is a sample API with Fiber and Swagger
// @host http://103.168.56.249:8080
// @BasePath /
func main() {
	app := apiRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}

func apiRoutes() *fiber.App {
	app := fiber.New()

	// Swagger documentation route
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	// Custom Swagger configuration
	app.Get("/swagger/*", fiberSwagger.New(fiberSwagger.Config{
		URL:          "http://example.com/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
		OAuth: &fiberSwagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))

	app.Use(cors.New())

	// API routes
	//app.Mount("/api", app)

	// Authentication routes
	app.Route("/auth", func(router fiber.Router) {
		router.Post("/signup/admin", auth.SignUpAdmin)
		router.Post("/signup/influencer", auth.SignUpInfluencer)
		router.Post("/signup/company", auth.SignUpCompany)
		router.Post("/login", auth.SignInUser)
		router.Get("/logout", middleware.DeserializeUser, auth.LogoutUser)
	})

	// Үндсэн хэрэглэгчийн мэдээлэл
	app.Route("/users", func(router fiber.Router) {
		router.Get("/me", middleware.DeserializeUser, user.GetMe)
		router.Post("/list", middleware.DeserializeUser, user.GetUserList)
		router.Put("/:user_id", middleware.DeserializeUser, user.UpdateUser)
	})

	// File
	app.Route("/file", func(router fiber.Router) {
		router.Get("/:name", middleware.DeserializeUser, file.GetFile)
	})

	// Хэрэглэгчийн эрх
	app.Route("/role", func(router fiber.Router) {
		router.Post("/list", middleware.DeserializeUser, role.GetRoleList)
		router.Post("/create", middleware.DeserializeUser, role.CreateRole)
		router.Put("/:id", middleware.DeserializeUser, role.UpdateRole)
		router.Delete("/:id", middleware.DeserializeUser, role.DeleteRole)
	})

	// Агент
	app.Route("/agent", func(router fiber.Router) {
		router.Post("/list", middleware.DeserializeUser, agency.GetAgentList)
		router.Get("/:id", middleware.DeserializeUser, agency.GetAgent)
		router.Post("/create", middleware.DeserializeUser, agency.CreateAgency)
		router.Put("/:id", middleware.DeserializeUser, agency.UpdateAgent)
		router.Delete("/:id", middleware.DeserializeUser, agency.DeleteUser)
	})

	app.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "JWT Authentication with Golang, Fiber, and GORM",
		})
	})

	app.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exist on this server", path),
		})
	})

	return app
}
