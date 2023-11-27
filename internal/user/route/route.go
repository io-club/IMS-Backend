package route

import (
	"ims-server/internal/user/api"
	ioginx "ims-server/pkg/ginx"
	"net/http"
)

var Routes = []ioginx.Route{
	// User,
	// TODO: 删除私有创建示例
	{Func: api.CreateUser(), FuncName: "CreateUser", Methods: []string{http.MethodPost}, Permission: nil, Private: true},
	{Func: api.GetUserByID(), FuncName: "GetUserByID", Methods: []string{http.MethodGet, http.MethodPost}, Permission: nil},
	{Func: api.MGetUserByIDs(), FuncName: "MGetUserByIDs", Methods: []string{http.MethodGet, http.MethodPost}, Permission: nil},
	{Func: api.UpdateUserByID(), FuncName: "UpdateUserByID", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api.DeleteUserByID(), FuncName: "DeleteUserByID", Methods: []string{http.MethodPost}, Permission: nil},
}
