package route

import (
	"ims-server/microservices/device/internal/api"
	ioginx "ims-server/pkg/ginx"
	"net/http"
)

var Routes = []ioginx.Route{
	// TODO: 注册接口@姚礼兴
	// TODO：考虑权限问题
	{Func: api.GetDataByID(), FuncName: "GetDataByID", Methods: []string{http.MethodGet, http.MethodPost}, Permission: nil},
}
