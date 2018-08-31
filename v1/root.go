package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// ShowBottle godoc
// @Summary Show a bottle
// @Description get string by ID
// @ID get-string-by-int
// @Tags bottles
// @Accept  json
// @Produce  json
// @Param  id path int true "Bottle ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /v1/ [get]
func (controller *ApiController) GetRoot(c *gin.Context) {
	c.JSON(200, gin.H{
		"request": fmt.Sprintf("%+v", c.Request.URL),
	})
}

func (controller *ApiController) Register() {
	v1 := controller.handler.Group("/v1")

	v1.GET("/", controller.GetRoot)
}
