package auth

import (
	"example.com/ramen/models/role"
	"fmt"
	"net/http"
	"strings"
	"time"

	"example.com/ramen/initializers"
	_map "example.com/ramen/models/map"
	reference2 "example.com/ramen/models/reference"
	"example.com/ramen/models/user"
	"example.com/ramen/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// SignUpAdmin godoc
// @Summary Create a new admin
// @Description Create a new admin
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body user.SignUpInput true "User"
// @Success 201 {object} string
// @Failure 400 {object} string
// @Router /auth/signup/admin [post]
func SignUpAdmin(c *fiber.Ctx) error {
	var payload *user.SignUpInput
	tx := initializers.DB.Begin()
	if err := c.BodyParser(&payload); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := user.ValidateStruct(payload)
	if errors != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})

	}

	if payload.Password != payload.PasswordConfirm {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Passwords do not match"})

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}

	newUser := user.User{
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: string(hashedPassword),
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})

	}

	var role role.Role
	resultRef := tx.Where("id = ?", 1).First(&role)
	if resultRef.RowsAffected == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Role олдсонгүй"})

	}

	var roleMap _map.RoleMap
	roleMap.EntityId = newUser.ID.String()
	roleMap.Name = role.Name
	roleMap.RoleId = role.Id
	roleMap.EntityName = "Admin"

	err1 := tx.Create(&roleMap)
	if err1.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: err.Error()})

	}

	if payload.Photo.Base64 != "" {
		err := utils.FileUpload(payload.Photo.Base64, newUser.ID, "Influencer", tx)
		if err != nil {
			tx.Rollback()
			return c.Status(http.StatusOK).JSON(utils.ResponseObj{ResponseCode: http.StatusBadRequest,
				ResponseMsg: err.Error()})
		}
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай бүртгэлээ", Data: newUser})
}

// SignUpInfluencer godoc
// @Summary Sign up influencer
// @Description Sign up influencer
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body  user.SignUpInfluencer true "User"
// @Success 201 {object} string
// @Failure 400 {object} string
// @Router /auth/signup/influencer [post]
func SignUpInfluencer(c *fiber.Ctx) error {
	var payload *user.SignUpInfluencer
	tx := initializers.DB.Begin()
	if err := c.BodyParser(&payload); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{
			ResponseCode: fiber.StatusBadRequest,
			ResponseMsg:  err.Error(),
		})
	}

	errors := user.ValidateStruct(payload)
	if errors != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(
			utils.ResponseObj{
				ResponseCode: fiber.StatusBadRequest,
				ResponseMsg:  "Утга зөв эсэхийг шалгана уу",
			})

	}

	if payload.Password != payload.PasswordConfirm {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(
			utils.ResponseObj{
				ResponseCode: fiber.StatusBadRequest,
				ResponseMsg:  "Passwords do not match",
			})

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(
			utils.ResponseObj{
				ResponseCode: fiber.StatusBadRequest,
				ResponseMsg:  err.Error(),
			})
	}

	newUser := user.User{
		Name:              payload.Name,
		Email:             strings.ToLower(payload.Email),
		Password:          string(hashedPassword),
		InfluencerIgName:  payload.IgName,
		Followers:         payload.Followers,
		Location:          payload.Location,
		EngagementRate:    payload.EngagementRate,
		AverageLikes:      payload.AverageLikes,
		Bio:               payload.Bio,
		TotalPosts:        payload.TotalPosts,
		AvgLikes:          payload.AvgLikes,
		AvgComments:       payload.AvgComments,
		AvgViews:          payload.AvgViews,
		AvgReelPlays:      payload.AvgReelPlays,
		GenderSplit:       payload.GenderSplit,
		AudienceInterests: payload.AudienceInteres,
		PopularPosts:      payload.PopularPosts,
		PhoneNumber:       payload.PhoneNumber,
	}

	result := tx.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(
			utils.ResponseObj{
				ResponseCode: fiber.StatusBadRequest,
				ResponseMsg:  "User with that email already exists",
			})
	} else if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{
			ResponseCode: fiber.StatusBadRequest,
			ResponseMsg:  result.Error.Error(),
		})
	}

	var role role.Role
	resultRef := tx.Where("id = ?", 2).First(&role)
	if resultRef.RowsAffected == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Role олдсонгүй"})

	}

	var roleMap _map.RoleMap
	roleMap.EntityId = newUser.ID.String()
	roleMap.Name = role.Name
	roleMap.RoleId = role.Id
	roleMap.EntityName = "Company"

	err1 := tx.Create(&roleMap)
	if err1.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: err.Error()})

	}

	if payload.Photo.Base64 != "" {
		err := utils.FileUpload(payload.Photo.Base64, newUser.ID, "Influencer", tx)
		if err != nil {
			tx.Rollback()
			return c.Status(http.StatusOK).JSON(utils.ResponseObj{ResponseCode: http.StatusBadRequest,
				ResponseMsg: err.Error()})
		}
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай бүртгэлээ", Data: newUser})
}

// SignUpCompany godoc
// @Summary Sign up company
// @Description Sign up company
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body user.SignUpCompany true "User"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /auth/signup/company [post]
func SignUpCompany(c *fiber.Ctx) error {
	var payload *user.SignUpCompany
	tx := initializers.DB.Begin()
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := user.ValidateStruct(payload)
	if errors != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})

	}

	if payload.Password != payload.PasswordConfirm {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: "Passwords do not match"})

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	newUser := user.User{
		Name:               payload.Name,
		Email:              strings.ToLower(payload.Email),
		Password:           string(hashedPassword),
		CompanyAccount:     payload.CompanyAccount,
		Location:           payload.Location,
		PhoneNumber:        payload.PhoneNumber,
		ManagerPhoneNumber: payload.ManagerPhoneNumber,
	}

	result := tx.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{
			ResponseCode: fiber.StatusBadRequest,
			ResponseMsg:  "User with that email already exists",
		})
	} else if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadGateway).JSON(
			utils.ResponseObj{
				ResponseCode: fiber.StatusBadGateway,
				ResponseMsg:  "Something bad happened",
			})
	}

	var role role.Role
	resultRef := tx.Where("id = ?", 3).First(&role)
	if resultRef.RowsAffected == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Role олдсонгүй"})

	}

	var roleMap _map.RoleMap
	roleMap.EntityId = newUser.ID.String()
	roleMap.Name = role.Name
	roleMap.RoleId = role.Id
	roleMap.EntityName = "Company"

	err1 := tx.Create(&roleMap)
	if err1.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Алдаа гарлаа", Data: err.Error()})

	}

	if payload.ProleId != "" {
		var reference reference2.Reference
		result := initializers.DB.Where("id = ?", payload.ProleId).First(&reference)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Reference олдсонгүй"})
		}
		var mapDb _map.Map
		mapDb.EntityId = newUser.ID.String()
		mapDb.Name = reference.Name
		mapDb.ReferenceId = reference.ID
		mapDb.EntityName = "Company"

		err := tx.Create(&mapDb)
		if err.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: err.Error.Error()})
		}

	}

	if payload.Photo.Base64 != "" {
		err := utils.FileUpload(payload.Photo.Base64, newUser.ID, "Company", tx)
		if err != nil {
			tx.Rollback()
			return c.Status(http.StatusOK).JSON(utils.ResponseObj{ResponseCode: http.StatusBadRequest,
				ResponseMsg: err.Error()})
		}
	}

	tx.Commit()
	return c.Status(fiber.StatusOK).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
		ResponseMsg: "Амжилттай бүртгэлээ", Data: newUser})

}

// SignInUser godoc
// @Summary Sign in user
// @Description Sign in user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body user.SignInInput true "User"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /auth/login [post]
func SignInUser(c *fiber.Ctx) error {
	var payload *user.SignInInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := user.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusOK).JSON(errors)

	}

	var user user.User
	result := initializers.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
	}

	config, _ := initializers.LoadConfig(".")

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = user.ID
	claims["exp"] = now.Add(config.JwtExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(config.JwtSecret))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   config.JwtMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "token": tokenString})
}

func LogoutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
