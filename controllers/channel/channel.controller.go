package channel

import (
	"example.com/ramen/initializers"
	"example.com/ramen/models/channel"
	"example.com/ramen/models/file"
	_map "example.com/ramen/models/map"
	reference2 "example.com/ramen/models/reference"
	"example.com/ramen/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
	"reflect"
)

// ListChannel godoc
// @Summary Get channel list
// @Description Get channel list
// @Tags Channel
// @Accept json
// @Produce json
// @Param channel body utils.RequestObj true "RequestObj"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /channel/list [post]
func ListChannel(c *fiber.Ctx) error {
	var channels []channel.Channel
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

			where = where + "ID in (" + "select channel_type_entity_id from maps where deleted_at is" +
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
			Model(&channel.Channel{}).Preload("Type").Preload("Image").Scopes(
			filter).Where(where)
	} else {
		conn = initializers.DB.
			Model(&channel.Channel{}).Preload("Type").Preload("Image").Scopes(
			filter)
	}

	pagination := utils.Pagination{CurrentPageNo: request.PageNo, PerPage: request.PerPage, Sort: request.Sort}
	conn.Debug().
		Scopes(utils.Paginate(channels, &pagination, conn)).
		Find(&channels)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: channels, Pagination: &pagination})

}

// CreateChannel godoc
// @Summary Create channel
// @Description Create channel
// @Tags Channel
// @Accept json
// @Produce json
// @Param channel body channel.CreateChannel true "CreateChannel"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /channel [post]
func CreateChannel(c *fiber.Ctx) error {
	var payload channel.CreateChannel
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})

	}
	tx := initializers.DB.Begin()

	var channel channel.Channel
	channel.Name = payload.Name
	channel.Description = payload.Description
	channel.Email = payload.Email
	channel.Phone = payload.Phone
	channel.Website = payload.Website
	channel.Address = payload.Address
	channel.TvDailyAvgViews = payload.TvDailyAvgViews
	channel.TvUnivisionNumber = payload.TvUnivisionNumber
	channel.FmDailyAvg1 = payload.FmDailyAvg1
	channel.FmDailyAvg2 = payload.FmDailyAvg2
	channel.FmSecondEval1 = payload.FmSecondEval1
	channel.FmSecondEval2 = payload.FmSecondEval2
	channel.FmCpe1 = payload.FmCpe1
	channel.FmCpe2 = payload.FmCpe2

	result := tx.Create(&channel)
	if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}

	if payload.Image.Base64 != "" {
		if payload.Image.Base64 != "" {
			err := utils.FileUpload(payload.Image.Base64, channel.ID.String(), "Channel", tx)
			if err != nil {
				tx.Rollback()
				return c.Status(http.StatusOK).JSON(utils.ResponseObj{ResponseCode: http.StatusBadRequest,
					ResponseMsg: err.Error()})
			}
		}
	}

	if payload.Type != "" {
		var referrers reference2.Reference
		result := tx.Where("id = ?", payload.Type).First(&referrers)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: "Төрөл олдсонгүй"})
		}
		var map1 _map.Map
		map1.ChannelTypeEntityId = channel.ID.String()
		map1.ReferenceId = referrers.ID
		map1.Name = referrers.Name
		map1.EntityName = "channelType"

		result = tx.Create(&map1)
		if result.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
		}

	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай бүртгэлээ", Data: channel})

}

// GetChannel godoc
// @Summary Get channel
// @Description Get channel
// @Tags Channel
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /channel/{id} [get]
func GetChannel(c *fiber.Ctx) error {
	id := c.Params("id")
	var channel channel.Channel
	result := initializers.DB.Where("id = ?", id).Preload("Type").Preload("Image").First(&channel)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Мэдээлэл олдсонгүй"})
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: channel})

}

// UpdateChannel godoc
// @Summary Update channel
// @Description Update channel
// @Tags Channel
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param channel body channel.UpdateChannel true "UpdateChannel"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /channel/{id} [put]
func UpdateChannel(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload channel.UpdateChannel
	var channel channel.Channel

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})
	}
	tx := initializers.DB.Begin()

	result := tx.Where("id = ?", id).First(&channel)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Мэдээлэл олдсонгүй"})
	}

	channel.Name = payload.Name
	channel.Description = payload.Description
	channel.Email = payload.Email
	channel.Phone = payload.Phone
	channel.Website = payload.Website
	channel.Address = payload.Address
	channel.TvDailyAvgViews = payload.TvDailyAvgViews
	channel.TvUnivisionNumber = payload.TvUnivisionNumber
	channel.FmDailyAvg1 = payload.FmDailyAvg1
	channel.FmDailyAvg2 = payload.FmDailyAvg2
	channel.FmSecondEval1 = payload.FmSecondEval1
	channel.FmSecondEval2 = payload.FmSecondEval2
	channel.FmCpe1 = payload.FmCpe1
	channel.FmCpe2 = payload.FmCpe2

	result = tx.Save(&channel)
	if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}

	if payload.Image.Base64 != "" {
		var file []file.File
		tx.Where("channel_parent_id = ?", channel.ID).Delete(&file)
		err := utils.FileUpload(payload.Image.Base64, channel.ID.String(), "Channel", tx)
		if err != nil {
			tx.Rollback()
			return c.Status(http.StatusOK).JSON(utils.ResponseObj{ResponseCode: http.StatusBadRequest,
				ResponseMsg: err.Error()})
		}
	}

	if payload.Type != "" {
		var referrers reference2.Reference
		result := tx.Where("id = ?", payload.Type).First(&referrers)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: "Төрөл олдсонгүй"})
		}
		var map1 _map.Map
		result1 := tx.Where("channel_type_entity_id = ?", channel.ID).Find(&map1)
		if result1.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: "Төрөл олдсонгүй"})
		}
		map1.ChannelTypeEntityId = channel.ID.String()
		map1.ReferenceId = referrers.ID
		map1.Name = referrers.Name
		map1.EntityName = "channelType"

		result = tx.Save(&map1)
		if result.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
		}
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай бүртгэлээ", Data: channel})

}

// DeleteChannel godoc
// @Summary Delete channel
// @Description Delete channel
// @Tags Channel
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Security ApiKeyAuth
// @Success 200 {object} utils.ResponseObj
// @Failure 400 {object} utils.ResponseObj
// @Router /channel/{id} [delete]
func DeleteChannel(c *fiber.Ctx) error {
	id := c.Params("id")
	var channel channel.Channel
	result := initializers.DB.Where("id = ?", id).First(&channel)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Мэдээлэл олдсонгүй"})
	}

	initializers.DB.Delete(&channel, id)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай устгагдлаа"})

}
