package handlers

import (
	"blur-rest/blur"
	"blur-rest/config"
	"fmt"
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
		c.JSON(404, config.ErrorWithStatus{
			Code:    config.GUIDNotFoundError,
			Message: "Provided GUID does not exist",
		})
		return
	}

	// Create temp file for source image
	srcImage, err := ioutil.TempFile("", "source-image")
	if err != nil {
		errMsg := fmt.Sprintf("Could not create source image: %s", err.Error())
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.TempFileError,
			Message: errMsg,
		})
		return
	}
	defer os.Remove(srcImage.Name())

	// Read uploaded data
	body := c.Request.Body
	content, err := ioutil.ReadAll(body)
	if err != nil {
		errMsg := fmt.Sprintf("Could not read uploded raw data: %s", err.Error())
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.ReadError,
			Message: errMsg,
		})
		return
	}

	// Write data into the file
	if _, err := srcImage.Write(content); err != nil {
		errMsg := fmt.Sprintf("Could not write uploded raw data: %s", err.Error())
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.WriteError,
			Message: errMsg,
		})
		return
	}

	// Open source image file
	srcImage, err = os.Open(srcImage.Name())
	if err != nil {
		errMsg := fmt.Sprintf("Could not open source image file: %s", err.Error())
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.CloseError,
			Message: errMsg,
		})
		return
	}

	// Create destionation image file
	destImage, err := ioutil.TempFile("", "destination-image")
	if err != nil {
		errMsg := fmt.Sprintf("Could not create destination image: %s", err.Error())
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.TempFileError,
			Message: errMsg,
		})
		return
	}
	defer os.Remove(destImage.Name())

	// Blur
	err = blur.Blur(srcImage, destImage, meta.Blur)
	if err != nil {
		errMsg := fmt.Sprintf("Could not blur the image: %s", err.Error())
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.BlurError,
			Message: errMsg,
		})
		return
	}

	// Close source image file
	if err := srcImage.Close(); err != nil {
		errMsg := fmt.Sprintf("Could not close source image file: %s", err.Error())
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.CloseError,
			Message: errMsg,
		})
		return
	}

	// Close destination image file
	if err := destImage.Close(); err != nil {
		errMsg := fmt.Sprintf("Could not close destination image file: %s", err.Error())
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.CloseError,
			Message: errMsg,
		})
		return
	}

	// Return blurred image
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Type", "application/octet-stream")
	c.File(destImage.Name())

	// Delete entry from DB
	go config.GlobalConfig.Store.Delete(guid)
}
