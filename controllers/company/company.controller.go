package company

import (
	"net/http"

	"example.com/ramen/initializers"
	"example.com/ramen/models/Company"
	"example.com/ramen/models/file"
	"example.com/ramen/utils"
	"github.com/gofiber/fiber/v2"
)

// ListCompany godoc
// @Summary Get company list
// @Description Get company list
// @Tags Company
// @Accept json
// @Produce json
// @Param company body utils.RequestObj true "RequestObj"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /company/list [post]
func ListCompany(c *fiber.Ctx) error {
	var company []Company.Company
	var request utils.RequestObj
	var conn = initializers.DB

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}

	conn = initializers.DB.
		Model(&Company.Company{}).Preload("Image").
		Scopes(utils.Filter(request.Filter))

	pagination := utils.Pagination{CurrentPageNo: request.PageNo, PerPage: request.PerPage, Sort: request.Sort}
	conn.Debug().
		Scopes(utils.Paginate(company, &pagination, conn)).
		Find(&company)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: company, Pagination: &pagination})

}

// GetCompany godoc
// @Summary Get company
// @Description Get company
// @Tags Company
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /company/{id} [get]
func GetCompany(ctx *fiber.Ctx) error {
	var company Company.Company
	id := ctx.Params("id")

	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Id оруулна уу"})
	}

	result := initializers.DB.Where("id = ?", id).Preload("Image").First(&company)
	if result.RowsAffected == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Олдсонгүй"})
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: company})

}

// CreateCompany godoc
// @Summary Create company
// @Description Create company
// @Tags Company
// @Accept json
// @Produce json
// @Param company body Company.CreateCompany true "CreateCompany"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /company [post]
func CreateCompany(c *fiber.Ctx) error {
	var payload Company.CreateCompany
	var company Company.Company
	tx := initializers.DB.Begin()
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}
	company.Name = payload.Name
	company.Description = payload.Description
	company.Website = payload.Website
	company.Email = payload.Email
	company.Phone = payload.Phone
	company.Address = payload.Address
	company.City = payload.City
	company.AreasOfActivity = payload.AreasOfActivity

	if err := tx.Create(&company).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}
	if payload.Image != "" {
		err := utils.FileUpload(payload.Image, company.ID, "company", tx)
		if err != nil {
			tx.Rollback()
			return c.Status(http.StatusOK).JSON(utils.ResponseObj{ResponseCode: http.StatusBadRequest,
				ResponseMsg: err.Error()})
		}
	}
	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай бүртгэлээ", Data: company})

}

// UpdateCompany godoc
// @Summary Update company
// @Description Update company
// @Tags Company
// @Accept json
// @Produce json
// @Param company body Company.UpdateCompany true "UpdateCompany"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /company [put]
func UpdateCompany(c *fiber.Ctx) error {
	var payload Company.UpdateCompany
	var company Company.Company
	tx := initializers.DB.Begin()
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}
	result := tx.Where("id = ?", payload.Id).First(&company)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Олдсонгүй"})

	}
	company.Name = payload.Name
	company.Description = payload.Description
	company.Website = payload.Website
	company.Email = payload.Email
	company.Phone = payload.Phone
	company.Address = payload.Address
	company.City = payload.City
	company.AreasOfActivity = payload.AreasOfActivity

	if err := tx.Save(&company).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}

	if payload.Image != "" {
		var file []file.File
		tx.Where("parent_id = ?", company.ID).Delete(&file)
		err := utils.FileUpload(payload.Image, company.ID, "company", tx)
		if err != nil {
			tx.Rollback()
			return c.Status(http.StatusOK).JSON(utils.ResponseObj{ResponseCode: http.StatusBadRequest,
				ResponseMsg: err.Error()})
		}
	}
	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: company})
}

// DeleteCompany godoc
// @Summary Delete company
// @Description Delete company
// @Tags Company
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /company/{id} [delete]
func DeleteCompany(c *fiber.Ctx) error {
	id := c.Params("id")
	var company Company.Company
	result := initializers.DB.Where("id = ?", id).Delete(&company)
	if result.Error != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай"})

}
