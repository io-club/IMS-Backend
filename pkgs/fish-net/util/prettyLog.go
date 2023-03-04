package util

import (
	"encoding/json"
	"fishnet/glb"
	"fmt"
)

func PrettyLog(objs ...interface{}) {
	for _, obj := range objs {
		empJSON, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			glb.LOG.Warn(err.Error())
		}
		glb.LOG.Info(fmt.Sprintln(string(empJSON)))
	}
}
