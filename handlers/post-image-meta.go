package handlers

import (
	"blur-rest/config"
	"blur-rest/store"

	"github.com/gin-gonic/gin"
)

type response struct {
	GUID string `json:"guid"`
	URL  string `json:"upload-url"`
}

// PostImageMetaHandler receives the JSON data posted by the user and stores it
func PostImageMetaHandler(c *gin.Context) {
	var err error

	var meta store.ImageMeta

	// Bind JSON
	if err = c.BindJSON(&meta); err != nil {
		c.JSON(400, config.ErrorWithStatus{
			Code:    config.Binding,
			Message: "Could not bind the JSON body",
		})
		return
	}

	// Validate meta and apply defaults
	if err = meta.Validate(); err != nil {
		c.JSON(400, config.ErrorWithStatus{
			Code:    config.JSONBody,
			Message: err.Error(),
		})
		return
	}
	meta.ApplyDefaults()

	c.JSON(200, response{})
}
