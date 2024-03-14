package utils

import (
	base642 "encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"example.com/ramen/models/file"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func FileUpload(base64Image string, ParentId string, Category string, tx *gorm.DB) error {
	var c *fiber.Ctx
	var fileName string
	decoded, err := base642.StdEncoding.DecodeString(base64Image)
	size := strconv.FormatInt(int64(len(base64Image)*3/4), 10)
	if err != nil {
		return err
	}
	// 5mb limit
	if len(decoded) > 5000000 {
		return c.Status(http.StatusOK).JSON(ResponseObj{ResponseCode: http.StatusBadRequest, ResponseMsg: "File size is too large"})
	}
	// generate uui for file name
	uid := uuid.New()

	fileName = uid.String() + ".png"

	if err = ioutil.WriteFile("./uploads/"+fileName, decoded, 0644); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(ResponseObj{ResponseCode: http.
			StatusBadRequest, ResponseMsg: "File uploads failed"})
	}
	var file file.File

	if Category == "" {
		err := tx.Where("parent_id = ? or company_parent_id = ? or influencer_parent_id = ?",
			ParentId, ParentId, ParentId).Find(&file)
		if err.Error != nil {
			return err.Error
		}
		file.FileName = fileName
		file.Size = size
		file.FilePath = fmt.Sprintf("/uploads/%s", fileName)
		if err := tx.Save(&file).Error; err != nil {
			return err
		}
		return nil
	}

	if Category == "Company" {
		file.CompanyParentId = ParentId
		file.Category = Category
		file.FileName = fileName
		file.Size = size
		file.FilePath = fmt.Sprintf("/uploads/%s", fileName)
	} else if Category == "Agency" {
		file.AgencyParentId = ParentId
		file.Category = Category
		file.FileName = fileName
		file.Size = size
		file.FilePath = fmt.Sprintf("/uploads/%s", fileName)
	} else if Category == "Channel" {
		file.ChannelParentId = ParentId
		file.Category = Category
		file.FileName = fileName
		file.Size = size
		file.FilePath = fmt.Sprintf("/uploads/%s", fileName)

	} else if Category == "Reference" {
		file.ReferenceParentId = ParentId
		file.Category = Category
		file.FileName = fileName
		file.Size = size
		file.FilePath = fmt.Sprintf("/uploads/%s", fileName)
	} else if Category == "Resource" {
		file.ResourceParentId = ParentId
		file.Category = Category
		file.FileName = fileName
		file.Size = size
		file.FilePath = fmt.Sprintf("/uploads/%s", fileName)
	} else if Category == "Campaigns" {
		file.CampaignsParentId = ParentId
		file.Category = Category
		file.FileName = fileName
		file.Size = size
		file.FilePath = fmt.Sprintf("/uploads/%s", fileName)
	} else {
		file.InfluencerParentId = ParentId
		file.Category = Category
		file.FileName = fileName
		file.Size = size
		file.FilePath = fmt.Sprintf("/uploads/%s", fileName)
	}

	if err := tx.Create(&file).Error; err != nil {
		return err
	}
	return nil
}

type Base64Struct struct {
	Base64 string `json:"base64"`
}
