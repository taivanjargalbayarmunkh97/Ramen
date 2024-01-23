package utils

import (
	base642 "encoding/base64"
	"example.com/ramen/initializers"
	"example.com/ramen/models/file"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func FileUpload(base64Image string, ParentId *uuid.UUID, Category string) error {
	var c *fiber.Ctx
	var fileName string
	decoded, err := base642.StdEncoding.DecodeString(base64Image)
	size := strconv.FormatInt(int64(len(base64Image)*3/4), 10)
	if err != nil {
		return err
	}
	// generate uui for file name
	uid := uuid.New()

	fileName = uid.String() + ".png"

	if err = ioutil.WriteFile("./uploads/"+fileName, decoded, 0644); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(ResponseObj{ResponseCode: http.
			StatusBadRequest, ResponseMsg: "File uploads failed"})
	}

	var file file.File

	file.ParentId = ParentId.String()
	file.Category = Category
	file.FileName = fileName
	file.Size = size
	file.FilePath = fmt.Sprintf("/uploads/%s", fileName)

	if err := initializers.DB.Create(&file).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(ResponseObj{ResponseCode: http.
			StatusBadRequest, ResponseMsg: "File uploads failed" + err.Error()})
	}
	return nil
}
