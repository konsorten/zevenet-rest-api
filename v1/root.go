package v1

import (
	"github.com/gin-gonic/gin"
)

func (controller *ApiController) GetRoot(c *gin.Context) {
	c.Redirect(301, "./system/version/")
}
