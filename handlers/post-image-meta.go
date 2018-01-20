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
			Code:    config.BindingError,
			Message: "Could not bind the JSON body",
		})
		return
	}

	// Validate meta and apply defaults
	if err = meta.Validate(); err != nil {
		c.JSON(400, config.ErrorWithStatus{
			Code:    config.JSONBodyError,
			Message: err.Error(),
		})
		return
	}
	meta.ApplyDefaults()

	// Store meta
	var guid string
	if guid, err = config.GlobalConfig.Store.Insert(meta); err != nil {
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.StoreError,
			Message: "Could not store meta",
		})
		return
	}

	r := response{
		GUID: guid,
		URL:  "http://" + config.GlobalConfig.Fqdn + ":" + config.GlobalConfig.Port + "/blur/" + guid,
	}

	c.JSON(200, r)
}
