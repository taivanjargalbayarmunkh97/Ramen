package news

import "github.com/gofiber/fiber/v2"

func CreateNews(c *fiber.Ctx) error {
	return c.SendString("Create News")
}
