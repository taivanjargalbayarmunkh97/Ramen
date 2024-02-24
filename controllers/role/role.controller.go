package role

import (
	"example.com/ramen/initializers"
	"example.com/ramen/models/role"
	"example.com/ramen/models/user"
	"example.com/ramen/utils"
	"github.com/gofiber/fiber/v2"
)

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role
// @Tags Role
// @Accept json
// @Produce json
// @Param role body role.RoleCreateInput true "Role"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /role/create [post]
func CreateRole(c *fiber.Ctx) error {
	var payload *role.RoleCreateInput
	var role role.Role

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := user.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})
	}

	role.Name = payload.Name
	role.Description = payload.Description
	role.Field1 = payload.Field1
	role.Field2 = payload.Field2
	role.Field3 = payload.Field3

	initializers.DB.Create(&role)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай бүртгэлээ", Data: role})

}

// GetRoleList godoc
// @Summary Get role list
// @Description Get role list
// @Tags Role
// @Accept json
// @Produce json
// @Param role body utils.RequestObj true "Role"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /role/list [post]
func GetRoleList(c *fiber.Ctx) error {
	var request utils.RequestObj
	var conn = initializers.DB
	var roles []role.Role

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}

	conn = initializers.DB.
		Model(&role.Role{}).
		Scopes(utils.Filter(request.Filter))

	pagination := utils.Pagination{CurrentPageNo: request.PageNo, PerPage: request.PerPage, Sort: request.Sort}
	conn.Debug().
		Scopes(utils.Paginate(roles, &pagination, conn)).
		Find(&roles)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: roles, Pagination: &pagination})

}

// UpdateRole godoc
// @Summary Update role
// @Description Update role
// @Tags Role
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param role body role.RoleUpdateInput true "Role"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /role/{id} [put]
func UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload *role.RoleUpdateInput
	var role role.Role

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := user.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})
	}

	result := initializers.DB.First(&role, id)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Хэрэглэгчийн эрх олдсонгүй"})
	}

	role.Name = payload.Name
	role.Description = payload.Description
	role.Field1 = payload.Field1
	role.Field2 = payload.Field2
	role.Field3 = payload.Field3

	initializers.DB.Save(&role)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай бүртгэлээ", Data: role})
}

// DeleteRole godoc
// @Summary Delete role
// @Description Delete role
// @Tags Role
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /role/{id} [delete]
func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var role role.Role

	result := initializers.DB.First(&role, id)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Хэрэглэгчийн эрх олдсонгүй"})
	}

	initializers.DB.Delete(&role)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай устгагдлаа"})

}
