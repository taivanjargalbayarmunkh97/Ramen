package role

import (
	"example.com/ramen/models/role"
	"example.com/ramen/models/user"
	"example.com/ramen/utils"
	"github.com/gofiber/fiber/v2"
)

//
//// GetMe godoc
//// @Summary Get user info
//// @Description Get user info
//// @Tags User
//// @Accept json
//// @Produce json
//// @Security ApiKeyAuth
//// @Success 200 {object} UserResponse
//// @Failure 400 {object} ErrorResponse
//// @Router /users/me [get]
//func GetMe(c *fiber.Ctx) error {
//	User := c.Locals("user").(user.UserResponse)
//
//	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай", Data: User})
//}

func CreateRole(c *fiber.Ctx) error {
	var payload *role.RoleCreateInput
	var role role.Role

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := user.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})
	}

	role.Name = payload.Name
	role.Description = payload.Description
	role.Field1 = payload.Field1
	role.Field2 = payload.Field2
	role.Field3 = payload.Field3

	return nil

}
