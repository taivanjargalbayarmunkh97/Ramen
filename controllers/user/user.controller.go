package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wpcodevo/golang-fiber-jwt/models/user"
	"github.com/wpcodevo/golang-fiber-jwt/utils"
)

// GetMe godoc
// @Summary Get user info
// @Description Get user info
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Router /users/me [get]
func GetMe(c *fiber.Ctx) error {
	User := c.Locals("user").(user.UserResponse)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай", Data: User})
}
