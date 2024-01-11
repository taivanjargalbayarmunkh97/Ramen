package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wpcodevo/golang-fiber-jwt/models/user"
	"github.com/wpcodevo/golang-fiber-jwt/utils"
)

func GetMe(c *fiber.Ctx) error {
	User := c.Locals("user").(user.UserResponse)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай", Data: User})
}
