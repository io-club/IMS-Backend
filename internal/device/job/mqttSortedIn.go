// 接收与发布
package job

import (
	"ims-server/pkg/mqtt"
)

func GetInDate() {
	client, err := iomqtt.NewClient()
	if err != nil {
		return
	}

	if err := client.Sub("", nil); err != nil {
		return
	}
}
