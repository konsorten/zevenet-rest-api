package helpers

import (
	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
)

func AddLinks(ctx *gin.Context, input *gabs.Container, links map[string]string) error {
	_, err := input.Set(ResolveRelativePath(ctx, "."), "links", "self")
	if err != nil {
		return err
	}

	for k, l := range links {
		_, err = input.Set(ResolveRelativePath(ctx, l), "links", k)
		if err != nil {
			return err
		}
	}

	return nil
}
