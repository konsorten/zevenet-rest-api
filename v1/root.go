package v1

import (
	"github.com/gin-gonic/gin"
)

// GetRoot retrieves basic system information like various version numbers.
// @Summary Get basic system information.
// @Description Retrieves basic system information like various version numbers.
// @ID get-root
// @Tags System
// @Accept  json
// @Produce  json
// @Success 200 {object} models.SystemVersion
// @Failure 400 {object} models.ApiError
// @Security ApiKeyAuth
// @Router /v1/ [get]
func (controller *ApiController) GetRoot(c *gin.Context) {
	res, err := controller.callZAPI("GET", "/system/version", nil)
	if err != nil {
		c.Error(err)
		return
	}

	resultBody := res.Content.S("params")

	resultBody.Set(mainVersionSimple, "rest_api_version")

	c.Data(res.HTTPStatus, JsonMimeType, resultBody.Bytes())
}
