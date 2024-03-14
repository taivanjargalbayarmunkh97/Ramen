package agency

import (
	"example.com/ramen/initializers"
	"example.com/ramen/models/Agency"
	"example.com/ramen/models/file"
	_map "example.com/ramen/models/map"
	reference2 "example.com/ramen/models/reference"
	"example.com/ramen/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"reflect"
)

// CreateAgency godoc
// @Summary Create a new agency
// @Description Create a new agency
// @Tags Agency
// @Accept json
// @Produce json
// @Param agency body Agency.CreateAgency true "Agency"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /agent/create [post]
func CreateAgency(c *fiber.Ctx) error {
	var payload *Agency.CreateAgency
	var agency Agency.Agency

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})

	}

	agency.Name = payload.Name
	agency.Address = payload.Address
	agency.Phone = payload.Phone
	agency.Email = payload.Email
	agency.Website = payload.Website
	agency.Description = payload.Description
	agency.Body = payload.Body
	agency.City = payload.City
	tx := initializers.DB.Begin()
	result := tx.Create(&agency)
	if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}

	if payload.Image.Base64 != "" {
		err := utils.FileUpload(payload.Image.Base64, agency.ID.String(), "Agency", tx)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: err.Error()})
		}
	}

	if len(payload.Brands) > 0 {
		var referrers []reference2.Reference
		var map1 []_map.Map
		result := tx.Where("id in ?", payload.Brands).Find(&referrers)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: "Бренд олдсонгүй"})
		}

		for _, v := range referrers {
			map1 = append(map1, _map.Map{AgencyBrandsEntityId: agency.ID.String(), EntityName: "agencyBrand",
				ReferenceId: v.ID, Name: v.Name})
		}

		result = tx.Create(&map1)
		if result.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
		}

	}

	if len(payload.Type) > 0 {
		var referrers []reference2.Reference
		result := tx.Debug().Where("id in ?", payload.Type).Find(&referrers)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: "Төрөл олдсонгүй"})
		}
		var map1 []_map.Map

		for _, ref := range referrers {
			map1 = append(map1, _map.Map{AgencyEntityId: agency.ID.String(), Name: ref.Name,
				ReferenceId: ref.ID, EntityName: "agent"})
		}

		result = tx.Save(&map1)
		if result.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
		}

	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай бүртгэлээ", Data: agency})

}

// UpdateAgent godoc
// @Summary Update agent
// @Description Update agent
// @Tags Agency
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param agency body Agency.UpdateAgency true "Agency"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /agent/{id} [put]
func UpdateAgent(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload *Agency.UpdateAgency
	var agent Agency.Agency
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})

	}

	tx := initializers.DB.Begin()

	result := tx.Model(&agent).Where("id = ?", id).Updates(Agency.Agency{Name: payload.Name,
		Address: payload.Address, Phone: payload.Phone, Email: payload.Email, Website: payload.Website,
		Description: payload.Description, Body: payload.Body, City: payload.City})
	if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}

	if len(payload.Type) > 0 {
		var reference []reference2.Reference
		result := tx.Where("id in ?", payload.Type).Find(&reference)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: " төрөл олдсонгүй"})
		}
		if len(reference) > 0 {
			var mapC []_map.Map
			tx.Where("agency_entity_id = ?", agent.ID).Delete(&_map.Map{})
			for _, ref := range reference {
				mapC = append(mapC, _map.Map{AgencyEntityId: agent.ID.String(),
					Name: ref.Name, ReferenceId: ref.ID,
					EntityName: "agent"})
			}
			if err := tx.Save(&mapC).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
					ResponseMsg: err.Error()})
			}
		}
	}

	if payload.Image.Base64 != "" {
		var file []file.File
		tx.Where("agency_parent_id = ?", id).Delete(&file)
		err := utils.FileUpload(payload.Image.Base64, id, "Agency", tx)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: err.Error()})
		}

	}

	if len(payload.Brands) > 0 {
		tx.Where("agency_brands_entity_id = ?", id).Delete(&_map.Map{})
		var referrers []reference2.Reference
		var map1 []_map.Map
		result := tx.Where("id in ?", payload.Brands).Find(&referrers)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: "Бренд олдсонгүй"})
		}

		for _, v := range referrers {
			map1 = append(map1, _map.Map{AgencyBrandsEntityId: id, EntityName: "agencyBrand",
				ReferenceId: v.ID, Name: v.Name})
		}

		result = tx.Save(&map1)
		if result.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})

		}
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: agent})

}

// GetAgentList godoc
// @Summary Get agent list
// @Description Get agent list
// @Tags Agency
// @Accept json
// @Produce json
// @Param agency body utils.RequestObj true "Agency"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /agent/list [post]
func GetAgentList(c *fiber.Ctx) error {
	var agents []Agency.Agency
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

			where = where + "ID in (" + "select agency_entity_id from maps where deleted_at is" +
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
			Model(&Agency.Agency{}).Preload("Type").Preload("Brands").Preload("Image").Scopes(
			filter).Where(where)
	} else {
		conn = initializers.DB.
			Model(&Agency.Agency{}).Preload("Type").Preload("Brands").Preload("Image").Scopes(
			filter)
	}

	pagination := utils.Pagination{CurrentPageNo: request.PageNo, PerPage: request.PerPage, Sort: request.Sort}
	conn.Debug().
		Scopes(utils.Paginate(agents, &pagination, conn)).
		Find(&agents)
	// list deer brandiin zurag
	//if len(agents) > 0 {
	//	for i, v := range agents {
	//		for ii, val := range v.Brands {
	//			var file file.File
	//			result := initializers.DB.Where("reference_parent_id = ?", val.ReferenceId).First(&file)
	//			if result.RowsAffected == 0 {
	//				continue
	//			}
	//			agents[i].Brands[ii].Image = file.FileName
	//		}
	//	}
	//
	//}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: agents, Pagination: &pagination})

}

// GetAgent godoc
// @Summary Get agent
// @Description Get agent
// @Tags Agency
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /agent/{id} [get]
func GetAgent(c *fiber.Ctx) error {
	id := c.Params("id")
	var agent Agency.Agency
	result := initializers.DB.Preload("Type").Preload("Brands").Preload("Image").First(&agent, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}

	for ii, val := range agent.Brands {
		var file file.File
		result := initializers.DB.Where("reference_parent_id = ?", val.ReferenceId).First(&file)
		if result.RowsAffected == 0 {
			continue
		}
		agent.Brands[ii].Image = file.FileName
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: agent})

}

// DeleteAgent godoc
// @Summary Delete user
// @Description Delete user
// @Tags Agency
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /agent/{id} [delete]
func DeleteAgent(c *fiber.Ctx) error {
	id := c.Params("id")
	var agent Agency.Agency
	result := initializers.DB.Delete(&agent, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: agent})

}
