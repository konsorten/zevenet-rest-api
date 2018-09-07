package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/konsorten/zevenet-rest-api/helpers"
	"github.com/konsorten/zevenet-rest-api/models"
)

// GetAllFarms retrieves a list of all farms.
// @Summary Get list of all farms.
// @Description Retrieves a list of all farms. The farms can be up or down.
// @ID get-farms
// @Tags Farms
// @Accept  json
// @Produce  json
// @Success 200 {array} models.FarmInfo
// @Failure 400 {object} models.ApiError
// @Security ApiKeyAuth
// @Router /v1/farms [get]
func (controller *ApiController) GetAllFarms(c *gin.Context) {
	resultBody, err := callZAPICached(controller.zapi, "/farms")
	if err != nil {
		if apiErr, ok := err.(models.ApiError); ok {
			c.AbortWithStatusJSON(apiErr.StatusCode, apiErr)
			return
		} else {
			c.AbortWithStatusJSON(400, models.NewApiError(400, "Unexpected error: %v", err))
			return
		}
	}

	farms, _ := resultBody.Children()
	for _, farm := range farms {
		farm.Delete("status") // remove uncached value
		farm.Delete("vport")  // remove profile-specific value

		helpers.AddLinks(c, farm, map[string]string{
			"detail": fmt.Sprintf("./%v/%v", farm.S("profile").Data(), farm.S("farmname").Data()),
		})
	}

	c.Data(200, JsonMimeType, resultBody.Bytes())
}
