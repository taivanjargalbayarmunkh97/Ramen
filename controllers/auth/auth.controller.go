package auth

import (
	"example.com/ramen/initializers"
	_map "example.com/ramen/models/map"
	reference2 "example.com/ramen/models/reference"
	"example.com/ramen/models/user"
	"example.com/ramen/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

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
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest, ResponseMsg: err.Error()})
	}

	errors := user.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Утга зөв эсэхийг шалгана уу", Data: errors})

	}

	if payload.Password != payload.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: "Passwords do not match"})

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})
	}

	newUser := user.User{
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: string(hashedPassword),
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
			ResponseMsg: err.Error()})

	}
	if payload.Photo.Base64 != "" {
		err := utils.FileUpload(payload.Photo.Base64, newUser.ID, "Influencer", tx)
		if err != nil {
			tx.Rollback()
			return c.Status(http.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: http.StatusBadRequest,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := user.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})

	}

	if payload.Password != payload.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Passwords do not match"})

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	newUser := user.User{
		Name:              payload.Name,
		Email:             strings.ToLower(payload.Email),
		Password:          string(hashedPassword),
		Role:              payload.RoleId,
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
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "User with that email already exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	if payload.Photo != "" {
		err := utils.FileUpload(payload.Photo, newUser.ID, "Influencer", tx)
		if err != nil {
			tx.Rollback()
			return c.Status(http.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: http.StatusBadRequest,
				ResponseMsg: err.Error()})
		}
	}

	tx.Commit()

	return c.Status(fiber.StatusCreated).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := user.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})

	}

	if payload.Password != payload.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Passwords do not match"})

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	newUser := user.User{
		Name:               payload.Name,
		Email:              strings.ToLower(payload.Email),
		Password:           string(hashedPassword),
		Role:               payload.RoleId,
		CompanyAccount:     payload.CompanyAccount,
		Location:           payload.Location,
		PhoneNumber:        payload.PhoneNumber,
		ManagerPhoneNumber: payload.ManagerPhoneNumber,
	}

	result := tx.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		tx.Rollback()
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "User with that email already exists"})
	} else if result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	if payload.ProleId != "" {
		var reference reference2.Reference
		result := initializers.DB.Where("id = ?", payload.ProleId).First(&reference)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
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
			return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseObj{ResponseCode: fiber.StatusBadRequest,
				ResponseMsg: "Алдаа гарлаа", Data: err.Error.Error()})
		}

	}
	tx.Commit()
	return c.Status(fiber.StatusCreated).JSON(utils.ResponseObj{ResponseCode: fiber.StatusOK,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := user.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	var user user.User
	result := initializers.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
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
