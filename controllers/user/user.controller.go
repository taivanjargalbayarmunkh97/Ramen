package user

import (
	"example.com/ramen/initializers"
	user "example.com/ramen/models/user"
	"example.com/ramen/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"reflect"
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
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Хэрэглэгч олдсонгүй id"})
	}

	result := initializers.DB.Where("id = ?", id).Preload("Photo").Preload("Photo1").Preload("Photo2").Preload("PRole").Preload("Role").First(&user)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Хэрэглэгч олдсонгүй"})
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
	var conn1 = initializers.DB

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}

	var where string
	var filter func(db *gorm.DB) *gorm.DB

	for i, v := range request.Filter {
		if v.FieldName == "role.id" {
			if where != "" {
				where += " and "
			}
			where = where + "ID in (" + "select entity_id from role_maps where deleted_at is" +
				" null and " +
				" role_id in ("
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
		conn1 = initializers.DB.
			Model(&user.User{}).Preload("Photo").Preload("Photo1").Preload("Photo2").Preload("PRole").Preload("Role").Scopes(
			filter).Where(where)
	} else {
		conn1 = initializers.DB.
			Model(&user.User{}).Preload("Photo").Preload("Photo1").Preload("Photo2").Preload("PRole").Preload("Role").Scopes(
			filter)
	}

	pagination := utils.Pagination{CurrentPageNo: request.PageNo, PerPage: request.PerPage, Sort: request.Sort}
	conn1.Debug().
		Scopes(utils.Paginate(users, &pagination, conn1)).
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
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "user_id оруулна уу"})
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}
	tx := initializers.DB.Begin()
	result := tx.Where("id = ?", id).First(&users)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Хэрэглэгч олдсонгүй"})
	}

	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
		}
		request.Password = string(hashedPassword)
	}
	var base64String string
	if request.Photo != nil && request.Photo.Base64 != "" {
		base64String = request.Photo.Base64
	}

	users.Name = request.Name
	users.Email = request.Email
	users.PhoneNumber = request.PhoneNumber
	users.Password = request.Password
	users.Provider = request.Provider
	users.Followers = request.Followers
	users.Location = request.Location
	users.EngagementRate = request.EngagementRate
	users.AverageLikes = request.AverageLikes
	users.Bio = request.Bio
	users.TotalPosts = request.TotalPosts
	users.AvgLikes = request.AvgLikes
	users.AvgComments = request.AvgComments
	users.AvgViews = request.AvgViews
	users.AvgReelPlays = request.AvgReelPlays
	users.GenderSplit = request.GenderSplit
	users.AudienceInterests = request.AudienceInterests
	users.PopularPosts = request.PopularPosts
	users.InfluencerIgName = request.InfluencerIgName
	users.CompanyAccount = request.CompanyAccount
	users.ManagerPhoneNumber = request.ManagerPhoneNumber

	result = tx.Save(&users)
	if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}

	if base64String != "" {

		err1 := utils.FileUpload(base64String, id, "", tx)
		if err1 != nil {
			tx.Rollback()
			return c.Status(http.StatusOK).JSON(utils.ResponseObj{ResponseCode: http.StatusBadRequest,
				ResponseMsg: err1.Error()})
		}

	}

	tx.Commit()
	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай шинэчиллээ", Data: users})

}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "ID"
// @Security ApiKeyAuth
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /users/{user_id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	var users user.User
	id := c.Params("user_id")

	if id == "" {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "user_id оруулна уу"})
	}

	tx := initializers.DB.Begin()
	result := tx.Where("id = ?", id).First(&users)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Хэрэглэгч олдсонгүй"})
	}

	result = tx.Delete(&users)
	if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: result.Error.Error()})
	}

	tx.Commit()
	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай устгагдлаа", Data: users})
}
