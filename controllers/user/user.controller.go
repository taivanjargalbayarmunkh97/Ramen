package user

import (
	"example.com/ramen/initializers"
	user "example.com/ramen/models/user"
	"example.com/ramen/utils"
	"github.com/gofiber/fiber/v2"
)

// GetMe godoc
// @Summary Get user info
// @Description Get user info
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /users/me [get]
func GetMe(c *fiber.Ctx) error {
	var user user.User
	id := c.Locals("user")

	if id == nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Хэрэглэгч олдсонгүй id"})
	}

	result := initializers.DB.Where("id = ?", id).Preload("Photo").Preload("PRole").First(&user)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Хэрэглэгч олдсонгүй"})
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай",
		Data: user})
}

// GetUserList godoc
// @Summary Get user list
// @Description Get user list
// @Tags User
// @Accept json
// @Produce json
// @Param user body utils.RequestObj true "RequestObj"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /users/list [post]
func GetUserList(c *fiber.Ctx) error {
	var users []user.UserSimple
	var request utils.RequestObj
	var conn = initializers.DB

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}

	conn = initializers.DB.
		Model(&user.User{}).Preload("Photo").Preload("PRole").
		Scopes(utils.Filter(request.Filter))

	pagination := utils.Pagination{CurrentPageNo: request.PageNo, PerPage: request.PerPage, Sort: request.Sort}
	conn.Debug().
		Scopes(utils.Paginate(users, &pagination, conn)).
		Find(&users)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: users, Pagination: &pagination})

}

// UpdateUser godoc
// @Summary Update user
// @Description Update user
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "ID"
// @Param user body user.UserUpdate true "UserUpdate"
// @Security ApiKeyAuth
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /users/{user_id} [put]
func UpdateUser(c *fiber.Ctx) error {
	var users user.User
	var request user.UserUpdate
	id := c.Params("user_id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "user_id оруулна уу"})
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}

	result := initializers.DB.Where("id = ?", id).First(&users)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Хэрэглэгч олдсонгүй"})
	}

	result = initializers.DB.Model(&users).Updates(request)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай шинэчиллээ", Data: users})

}
