package v1

import (
	"github.com/gin-gonic/gin"
)

// SystemVersion contains information about the system version.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#show-version
type SystemVersion struct {
	ApplianceVersion string `json:"appliance_version"`
	Hostname         string `json:"hostname"`
	KernelVersion    string `json:"kernel_version"`
	SystemDate       string `json:"system_date"`
	ZevenetVersion   string `json:"zevenet_version"`
	RestApiVersion   string `json:"rest_api_version"`
}

// GetRoot retrieves basic system information like various version numbers.
// @Summary Get basic system information.
// @Description Retrieves basic system information like various version numbers.
// @ID get-root
// @Tags System
// @Accept  json
// @Produce  json
// @Success 200 {object} v1.SystemVersion
// @Failure 400 {object} v1.ApiError
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
