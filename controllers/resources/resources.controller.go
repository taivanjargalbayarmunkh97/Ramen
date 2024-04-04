package resources

import (
	"example.com/ramen/initializers"
	"example.com/ramen/models/file"
	_map "example.com/ramen/models/map"
	reference2 "example.com/ramen/models/reference"
	"example.com/ramen/models/resources"
	"example.com/ramen/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"reflect"
)

// ListReference godoc
// @Summary List resources
// @Description List resources
// @Tags Resources
// @Accept  json
// @Produce  json
// @Param body body utils.RequestObj true "body"
// @Success 200 {object} utils.ResponseObj
// @Router /resources/list [post]
func ListReference(c *fiber.Ctx) error {
	var resourcesList []resources.Resources
	var request utils.RequestObj
	var conn = initializers.DB

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}

	var where string
	var filter func(db *gorm.DB) *gorm.DB

	for i, v := range request.Filter {
		if v.FieldName == "type.id" {
			if where != "" {
				where += " and "
			}

			where = where + "ID in (" + "select resources_entity_id from maps where deleted_at is" +
				" null and " +
				" reference_id in ("
			for i, val := range v.Values {
				if i > 0 {
					where += ", "
				}
				if reflect.TypeOf(val) == reflect.TypeOf(int(0)) {
					where += fmt.Sprintf("%d", val)
				} else {
					where += fmt.Sprintf("'%s'", val)
				}
			}
			where += "))"
			filter = utils.Filter(request.Filter[i+1:], request.GlobOperation)
		}
	}
	filter = utils.Filter(request.Filter, request.GlobOperation)

	if where != "" {
		conn = initializers.DB.
			Model(&resources.Resources{}).Preload("Type").Preload("Image").Scopes(
			filter).Where(where)
	} else {
		conn = initializers.DB.
			Model(&resources.Resources{}).Preload("Type").Preload("Image").Scopes(
			filter)
	}

	pagination := utils.Pagination{CurrentPageNo: request.PageNo, PerPage: request.PerPage, Sort: request.Sort}
	conn.Debug().
		Scopes(utils.Paginate(resourcesList, &pagination, conn)).
		Find(&resourcesList)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: resourcesList, Pagination: &pagination})

}

// CreateReference godoc
// @Summary Create resources
// @Description Create resources
// @Tags Resources
// @Accept  json
// @Produce  json
// @Param resources body resources.CreateResources true "CreateReference"
// @Success 200 {object} utils.ResponseObj
// @Router /resources [post]
func CreateReference(c *fiber.Ctx) error {
	var resource resources.Resources
	var payload resources.CreateResources

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})
	}

	resource.Name = payload.Name
	resource.Description = payload.Description
	resource.Body = payload.Body
	resource.YoutubeLink = payload.YoutubeLink
	resource.FacebookLink = payload.FacebookLink
	resource.TwitterLink = payload.TwitterLink
	resource.InstagramLink = payload.InstagramLink
	resource.LinkedinLink = payload.LinkedinLink
	resource.PinterestLink = payload.PinterestLink

	tx := initializers.DB.Begin()

	if err := tx.Create(&resource).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	if payload.Image.Base64 != "" {
		if err := utils.FileUpload(payload.Image.Base64, resource.ID.String(), "Resource", tx); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
		}
	}

	if len(payload.Type) > 0 {
		var reference []reference2.Reference
		var map1 []_map.Map

		result := tx.Where("id in ?", payload.Type).Find(&reference)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Төрөл олдсонгүй"})
		}
		for _, v := range reference {
			map1 = append(map1, _map.Map{
				Name:              v.Name,
				ReferenceId:       v.ID,
				ResourcesEntityId: resource.ID.String(),
				EntityName:        "resources",
			})

		}

		if err := tx.Save(&map1).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
		}

	}

	tx.Commit()
	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай бүртгэгдлээ"})

}

// UpdateReference godoc
// @Summary Update resources
// @Description Update resources
// @Tags Resources
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param resources body resources.UpdateResources true "UpdateReference"
// @Success 200 {object} utils.ResponseObj
// @Router /resources/{id} [put]
func UpdateReference(c *fiber.Ctx) error {
	var resource resources.Resources
	var payload resources.UpdateResources

	var id = c.Params("id")

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})
	}

	tx := initializers.DB.Begin()

	if result := tx.Where("id = ?", id).First(&resource).RowsAffected; result == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Мэдээлэл олдсонгүй"})
	}

	resource.Name = payload.Name
	resource.Description = payload.Description
	resource.Body = payload.Body
	resource.YoutubeLink = payload.YoutubeLink
	resource.FacebookLink = payload.FacebookLink
	resource.TwitterLink = payload.TwitterLink
	resource.InstagramLink = payload.InstagramLink
	resource.LinkedinLink = payload.LinkedinLink
	resource.PinterestLink = payload.PinterestLink

	if err := tx.Save(&resource).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	if payload.Image.Base64 != "" {
		var file []file.File
		tx.Where("resource_parent_id = ?", resource.ID).Delete(&file)
		if err := utils.FileUpload(payload.Image.Base64, resource.ID.String(), "Resource", tx); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
		}
	}

	if len(payload.Type) > 0 {
		var reference []reference2.Reference
		var map1 []_map.Map

		result := tx.Where("id in ?", payload.Type).Find(&reference)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Төрөл олдсонгүй"})
		}
		if err := tx.Where("resources_entity_id = ?", resource.ID.String()).Delete(&_map.Map{}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
		}

		for _, v := range reference {
			map1 = append(map1, _map.Map{
				Name:              v.Name,
				ReferenceId:       v.ID,
				ResourcesEntityId: resource.ID.String(),
				EntityName:        "resources",
			})
		}

		if err := tx.Save(&map1).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
		}
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай шинэчлэгдлээ"})

}

// GetReference godoc
// @Summary Get resources
// @Description Get resources
// @Tags Resources
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} utils.ResponseObj
// @Router /resources/{id} [get]
func GetReference(c *fiber.Ctx) error {
	var id = c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Id is empty"})
	}
	var resource resources.Resources

	if result := initializers.DB.Where("id = ?", id).Preload("Type").Preload("Image").First(&resource).
		RowsAffected; result == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Мэдээлэл олдсонгүй"})
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай", Data: resource})

}

// DeleteReference godoc
// @Summary Delete resources
// @Description Delete resources
// @Tags Resources
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} utils.ResponseObj
// @Router /resources/{id} [delete]
func DeleteReference(c *fiber.Ctx) error {
	var id = c.Params("id")
	var resource resources.Resources

	tx := initializers.DB.Begin()

	if result := tx.Where("id = ?", id).First(&resource).RowsAffected; result == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Мэдээлэл олдсонгүй"})
	}

	if err := tx.Where("resources_entity_id = ?", resource.ID.String()).Delete(&_map.Map{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	if err := tx.Where("resource_parent_id = ?", resource.ID).Delete(&file.File{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	if err := tx.Delete(&resource).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай устгагдлаа"})
}
