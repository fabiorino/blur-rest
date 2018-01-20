package handlers

import (
	"blur-rest/config"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

// UploadImageHandler handles the PUT request. The image uploaded by the user is stored in a temporary file.
// At the end of the computation, the image is returned as an octet-stream and the temporary file is removed
func UploadImageHandler(c *gin.Context) {
	guid := c.Param("guid")

	// Get meta from db
	meta, err := config.GlobalConfig.Store.Get(guid)
	if err != nil {
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.GUIDNotFoundError,
			Message: "Provided GUID does not exist",
		})
		return
	}

	// Create temp file for source image
	srcImage, err := ioutil.TempFile("", "source-image")
	if err != nil {
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.TempFileError,
			Message: "Could not create source image",
		})
		return
	}
	defer os.Remove(srcImage.Name())

	// Read uploaded data
	body := c.Request.Body
	content, err := ioutil.ReadAll(body)
	if err != nil {
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.ReadError,
			Message: "Could not read uploded raw data",
		})
		return
	}

	// Write data into the file
	if _, err := srcImage.Write(content); err != nil {
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.WriteError,
			Message: "Could not write uploded raw data",
		})
		return
	}

	// Close source image file
	if err := srcImage.Close(); err != nil {
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.CloseError,
			Message: "Could not close source image file",
		})
		return
	}
}
