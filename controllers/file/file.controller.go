package file

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

// GetFile HandlerGetFileByName godoc
//
//	@Summary		Get file by name
//	@Description	Get file by name
//	@Tags			File
//	@Accept			json
//	@Produce		json
//	@Param			name path string true "name"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/file/{name} [get]
func GetFile(c *fiber.Ctx) error {
	filename := c.Params("name")
	filePath := "./uploads/" + filename

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	}

	return c.SendFile(filePath)
}
