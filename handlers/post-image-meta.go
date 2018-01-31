package handlers

import (
	"blur-rest/config"
	"blur-rest/store"
	"fmt"

	"github.com/gin-gonic/gin"
)

// The JSON response message delivered to the user in case of success
type response struct {
	GUID string `json:"guid"`
	URL  string `json:"upload-url"`
}

// PostImageMetaHandler receives the JSON data posted by the user and stores it.
// Returns a JSON response to the user
func PostImageMetaHandler(c *gin.Context) {
	var err error

	var meta store.ImageMeta

	// Bind JSON
	if err = c.BindJSON(&meta); err != nil {
		errMsg := fmt.Sprintf("Could not bind the JSON body: %s", err.Error())
		c.JSON(400, config.ErrorWithStatus{
			Code:    config.BindingError,
			Message: errMsg,
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
		errMsg := fmt.Sprintf("Could not store meta: %s", err.Error())
		c.JSON(500, config.ErrorWithStatus{
			Code:    config.StoreError,
			Message: errMsg,
		})
		return
	}

	r := response{
		GUID: guid,
		URL:  "http://" + config.GlobalConfig.Fqdn + ":" + config.GlobalConfig.Port + "/blur/" + guid,
	}

	c.JSON(200, r)
}
