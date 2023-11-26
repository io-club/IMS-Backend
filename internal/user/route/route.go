package route

import (
	"ims-server/internal/user/api"
	ioginx "ims-server/pkg/ginx"
	"net/http"
)

var Routes = []ioginx.Route{
	// User,
	// TODO: 删除私有创建示例
	{Func: api.CreateUser(), FuncName: "CreateUser", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api.GetUserByID(), FuncName: "GetUserByID", Permission: nil, Methods: []string{http.MethodGet, http.MethodPost}},
	{Func: api.MGetUserByID(), FuncName: "MGetUserByID", Permission: nil, Methods: []string{http.MethodGet, http.MethodPost}},
	{Func: api.UpdateUserByID(), FuncName: "UpdateUserByID", Permission: nil, Methods: []string{http.MethodPost}},
	{Func: api.DeleteUserByID(), FuncName: "DeleteUserByID", Permission: nil, Methods: []string{http.MethodPost}},
}
