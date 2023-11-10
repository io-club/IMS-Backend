package util

import (
	"IMS-Backend/pkgs/fish-net/glb"
	"encoding/json"
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

func SPrettyLog(objs ...interface{}) string {
	res := ""
	for _, obj := range objs {
		empJSON, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			glb.LOG.Warn(err.Error())
		}
		res += fmt.Sprintln(string(empJSON))
		res += "\n"
	}
	return res
}
