package reference

import (
	"example.com/ramen/initializers"
	"example.com/ramen/models/reference"
	"example.com/ramen/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// CreateReference godoc
// @Summary Create a new reference
// @Description Create a new reference
// @Tags Reference
// @Accept json
// @Produce json
// @Param reference body CreateReference true "Reference"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /reference [post]
func CreateReference(c *fiber.Ctx) error {
	var payload *reference.CreateReference
	var reference reference.Reference

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})

	}

	reference.Name = payload.Name
	reference.Description = payload.Description
	reference.Field1 = payload.Field1
	reference.Field2 = payload.Field2
	reference.Field3 = payload.Field3
	reference.Code = payload.Code
	tx := initializers.DB.Begin()
	result := tx.Create(&reference)
	if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}

	if payload.Image.Base64 != "" {
		err := utils.FileUpload(payload.Image.Base64, strconv.Itoa(int(reference.ID)), "Reference", tx)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: err.Error()})
		}

	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай бүртгэлээ", Data: reference})
}

// ListReference godoc
// @Summary Get reference list
// @Description Get reference list
// @Tags Reference
// @Accept json
// @Produce json
// @Param reference body utils.RequestObj true "RequestObj"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /reference/list [post]
func ListReference(c *fiber.Ctx) error {
	var references []reference.Reference
	var request utils.RequestObj
	var conn = initializers.DB

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}

	conn = initializers.DB.
		Model(&reference.Reference{}).
		Scopes(utils.Filter(request.Filter, request.GlobOperation))

	pagination := utils.Pagination{CurrentPageNo: request.PageNo, PerPage: request.PerPage, Sort: request.Sort}
	conn.Debug().
		Scopes(utils.Paginate(references, &pagination, conn)).
		Find(&references)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: references, Pagination: &pagination})

}

// GetReference godoc
// @Summary Get reference
// @Description Get reference
// @Tags Reference
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /reference/{id} [get]
func GetReference(ctx *fiber.Ctx) error {
	var reference reference.Reference
	id := ctx.Params("id")

	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Олдсонгүй"})
	}

	result := initializers.DB.Where("id = ?", id).First(&reference)
	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Алдаа гарлаа"})
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай", Data: reference})

}

// UpdateReference godoc
// @Summary Update reference
// @Description Update reference
// @Tags Reference
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param reference body UpdateReference true "UpdateReference"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /reference/{id} [put]
func UpdateReference(c *fiber.Ctx) error {
	var referencedb reference.Reference
	var request reference.UpdateReference
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "reference_id оруулна уу"})
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}

	result := initializers.DB.Where("id = ?", id).First(&referencedb)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Reference олдсонгүй"})
	}

	referencedb.Name = request.Name
	referencedb.Description = request.Description
	referencedb.Field1 = request.Field1
	referencedb.Field2 = request.Field2
	referencedb.Field3 = request.Field3
	referencedb.Code = request.Code
	result = initializers.DB.Save(&referencedb)
	if result.Error != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай шинэчлэлээ", Data: referencedb})

}

// DeleteReference godoc
// @Summary Delete reference
// @Description Delete reference
// @Tags Reference
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /reference/{id} [delete]
func DeleteReference(ctx *fiber.Ctx) error {
	var reference reference.Reference
	id := ctx.Params("id")

	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Олдсонгүй"})
	}

	result := initializers.DB.Where("id = ?", id).Delete(&reference)
	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Алдаа гарлаа"})
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай устгалаа"})

}
