package news

import (
	"example.com/ramen/initializers"
	"example.com/ramen/models/file"
	_map "example.com/ramen/models/map"
	"example.com/ramen/models/news"
	reference2 "example.com/ramen/models/reference"
	"example.com/ramen/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"reflect"
)

// CreateNews godoc
// @Summary Create news
// @Description Create news
// @Tags News
// @Accept  json
// @Produce  json
// @Param resources body news.CreateNews true "CreateNews"
// @Success 200 {object} utils.ResponseObj
// @Router /news [post]
func CreateNews(c *fiber.Ctx) error {
	var payload news.CreateNews
	var news news.News

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})
	}

	news.Title = payload.Title
	news.Body = payload.Body
	news.Description = payload.Description

	tx := initializers.DB.Begin()

	if err := tx.Create(&news).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	if payload.Image.Base64 != "" {
		if err := utils.FileUpload(payload.Image.Base64, news.ID.String(), "News", tx); err != nil {
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
				Name:         v.Name,
				ReferenceId:  v.ID,
				NewsEntityId: news.ID.String(),
				EntityName:   "news",
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

func UpdateNews(c *fiber.Ctx) error {
	var payload news.UpdateNews
	var news news.News

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

	if result := tx.Where("id = ?", id).First(&news).RowsAffected; result == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Мэдээлэл олдсонгүй"})
	}

	news.Title = payload.Title
	news.Description = payload.Description
	news.Body = payload.Body

	if err := tx.Save(&news).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	if payload.Image.Base64 != "" {
		var file []file.File
		tx.Where("news_parent_id = ?", news.ID).Delete(&file)
		if err := utils.FileUpload(payload.Image.Base64, news.ID.String(), "News", tx); err != nil {
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
		if err := tx.Where("news_entity_id = ?", news.ID.String()).Delete(&_map.Map{}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
		}

		for _, v := range reference {
			map1 = append(map1, _map.Map{
				Name:         v.Name,
				ReferenceId:  v.ID,
				NewsEntityId: news.ID.String(),
				EntityName:   "news",
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

func ListNews(c *fiber.Ctx) error {
	var newsList []news.News
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

			where = where + "ID in (" + "select news_entity_id from maps where deleted_at is" +
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
			Model(&news.News{}).Preload("Type").Preload("Image").Scopes(
			filter).Where(where)
	} else {
		conn = initializers.DB.
			Model(&news.News{}).Preload("Type").Preload("Image").Scopes(
			filter)
	}

	pagination := utils.Pagination{CurrentPageNo: request.PageNo, PerPage: request.PerPage, Sort: request.Sort}
	conn.Debug().
		Scopes(utils.Paginate(newsList, &pagination, conn)).
		Find(&newsList)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай", Data: newsList, Pagination: &pagination})

}

func GetNews(c *fiber.Ctx) error {
	var id = c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Id is empty"})
	}
	var news news.News

	if result := initializers.DB.Where("id = ?", id).Preload("Type").Preload("Image").First(&news).
		RowsAffected; result == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Мэдээлэл олдсонгүй"})
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай", Data: news})

}

func DeleteNews(c *fiber.Ctx) error {
	var id = c.Params("id")
	var news news.News

	tx := initializers.DB.Begin()

	if result := tx.Where("id = ?", id).First(&news).RowsAffected; result == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Мэдээлэл олдсонгүй"})
	}

	if err := tx.Where("news_entity_id = ?", news.ID.String()).Delete(&_map.Map{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	if err := tx.Where("news_parent_id = ?", news.ID).Delete(&file.File{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	if err := tx.Delete(&news).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK, ResponseMsg: "Амжилттай устгагдлаа"})

}
