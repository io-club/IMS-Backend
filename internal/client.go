package main

import (
	"fmt"
	ioconfig "ims-server/pkg/config"
	ioginx "ims-server/pkg/ginx"
)

func main() {
	config := ioconfig.GetServiceConf().Nms
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	ioginx.NewIOServer(nil).Run(addr, config.Name)

}
