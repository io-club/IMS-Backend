package main

import (
	"fishnet/glb"
	"fishnet/initialize"
)

// init
func init() {
	initialize.InitAll()
}

func main() {
	glb.G.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
