package helpers

import (
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

func ResolveRelativePath(ctx *gin.Context, relPath ...string) string {
	elem := []string{
		strings.TrimSuffix(ctx.Request.URL.Path, "/"),
	}

	elem = append(elem, relPath...)

	return path.Join(elem...)
}
