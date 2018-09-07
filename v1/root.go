package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/konsorten/zevenet-rest-api/helpers"
)

func (controller *ApiController) GetRoot(c *gin.Context) {
	c.Redirect(301, helpers.ResolveRelativePath(c, "./system/version"))
}
