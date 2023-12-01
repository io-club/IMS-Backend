package iodb

import (
	"fmt"
	"ims-server/pkg/util"
	"testing"
)

func TestPageRequest(t *testing.T) {
	pr := PageRequest{
		Page:         0,
		Size:         10,
		Filter:       "( NOT page eq 1 AND page <= 10 AND name = 'test' ) OR ( NOT id = 1 AND name = 'test2' )",
		Search:       "",
		Sort:         "",
		filterFields: util.NewSet("id", "name", "page"),
	}
	pr.Build()
	fmt.Printf("expr: %+v\n", pr.expr)
}
