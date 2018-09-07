package helpers

import (
	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
)

func GetArray(input *gabs.Container, path ...string) ([]*gabs.Container, error) {
	c, err := input.ArrayCount(path...)
	if err != nil {
		return nil, err
	}

	ret := make([]*gabs.Container, 0)

	for i := 0; i < c; i++ {
		e, err := input.ArrayElement(i, path...)
		if err != nil {
			return nil, err
		}

		ret = append(ret, e)
	}

	return ret, nil
}

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
