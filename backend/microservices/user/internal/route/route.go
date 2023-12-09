package route

import (
	api2 "ims-server/microservices/user/internal/api"
	ioconsts "ims-server/pkg/consts"
	ioginx "ims-server/pkg/ginx"
	"ims-server/pkg/util"
	"net/http"
)

var Routes = []ioginx.Route{
	// User
	{Func: api2.GetUserByID(), FuncName: "GetUserByID", Methods: []string{http.MethodGet, http.MethodPost},
		Permission: util.NewSet(ioconsts.UserTypeAdmin, ioconsts.UserTypeInsider, ioconsts.UserTypeOutsiders),
	},
	{Func: api2.MGetUserByIDs(), FuncName: "MGetUserByIDs", Methods: []string{http.MethodGet, http.MethodPost},
		Permission: util.NewSet(ioconsts.UserTypeAdmin, ioconsts.UserTypeInsider, ioconsts.UserTypeOutsiders),
	},
	{Func: api2.GetUsers(), FuncName: "GetUsers", Methods: []string{http.MethodGet, http.MethodPost}, Permission: nil},
	{Func: api2.UpdateUserByID(), FuncName: "UpdateUserByID", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api2.UploadAvatar(), FuncName: "UploadAvatar", Methods: []string{http.MethodPost},
		Permission: util.NewSet(ioconsts.UserTypeAdmin, ioconsts.UserTypeInsider, ioconsts.UserTypeOutsiders),
	},
	{Func: api2.DeleteUserByID(), FuncName: "DeleteUserByID", Methods: []string{http.MethodPost}, Permission: nil},

	// Register
	{Func: api2.Register(), FuncName: "Register", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api2.SendVerification(), FuncName: "SendVerification", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api2.NameLogin(), FuncName: "NameLogin", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api2.EmailLogin(), FuncName: "EmailLogin", Methods: []string{http.MethodPost}, Permission: nil},
	{Func: api2.RetrievePassword(), FuncName: "RetrievePassword", Methods: []string{http.MethodPost}, Permission: nil},
}
