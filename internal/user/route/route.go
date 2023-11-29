package route

import (
	"ims-server/internal/user/api"
	ioginx "ims-server/pkg/ginx"
	"net/http"
)

var Routes = []ioginx.Route{
	// User
	{Func: api.GetUserByID(), FuncName: "GetUserByID", Methods: []string{http.MethodGet, http.MethodPost}, Permission: nil},
	{Func: api.MGetUserByIDs(), FuncName: "MGetUserByIDs", Methods: []string{http.MethodGet, http.MethodPost}, Permission: nil},
	{Func: api.UpdateUserByID(), FuncName: "UpdateUserByID", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api.DeleteUserByID(), FuncName: "DeleteUserByID", Methods: []string{http.MethodPost}, Permission: nil},

	// Register
	{Func: api.Register(), FuncName: "Register", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api.SendVerification(), FuncName: "SendVerification", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api.NameLogin(), FuncName: "NameLogin", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api.EmailLogin(), FuncName: "EmailLogin", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api.RetrievePassword(), FuncName: "RetrievePassword", Methods: []string{http.MethodPost}, Permission: nil},
}
