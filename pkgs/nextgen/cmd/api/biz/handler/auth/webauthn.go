package webauthn

import (
	"context"
	"encoding/binary"
	"fish_net/cmd/api/glb"
	"fish_net/cmd/user/domain"
	userDomain "fish_net/cmd/user/domain"
	"fish_net/cmd/user/usecase"
	"fish_net/cmd/api/rpc"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/hertz-contrib/sessions"
)

var _userUsecase userDomain.UserUsecase
var _webAuthnCredentialUsecase domain.WebAuthnCredentialUsecase

func init() {
	_userUsecase = usecase.NewUserUsecase()
	_webAuthnCredentialUsecase = usecase.NewWebAuthnCredentialUsecase()
}

const REGISTER_SESSION_DATA_KEY = "register_session_data"
const LOGIN_SESSION_DATA_KEY = "login_session_data"

func ByteToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}

func Play(ctx context.Context, c *app.RequestContext) {

	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	if v == nil {
		count = 0
	} else {
		count = v.(int)
		count++
	}
	session.Set("count", count)
	session.Save()
	c.JSON(200, utils.H{"count": count})
}

func RegisterBegin(ctx context.Context, c *app.RequestContext) {
	userName, ok := c.Params.Get("username")
	if !ok {
		c.JSON(http.StatusBadRequest, utils.H{
			"msg": "muse give a username",
		})
		return
	}
	klog.Info("RegisterBegin username: " + userName)

	rpc.
	// 没有用户就新建
	users, count, err := _userUsecase.QueryUser(ctx, &userName, 1, 0)
	if err != nil || count == 0 {
		klog.Info("Creating user: " + userName)
		user := &domain.User{
			Username: userName,
			Nickname: userName,
			Avater:   "https://pics.com/avatar.png",
		}
		err = _userUsecase.CreateUser(ctx, []*domain.User{user})
		if err != nil {
			klog.Warn("USER CREATION ERROR: " + err.Error())
			c.JSON(http.StatusInternalServerError, utils.H{
				"failed": true,
				"msg":    "user not created",
			})
			return
		}
	} else {
		klog.Info("user exists")
	}

	users, count, err = _userUsecase.QueryUser(ctx, &userName, 1, 0)
	user := users[0]
	user.Credentials = _webAuthnCredentialUsecase.QueryCredential(int64(user.ID))

	// 检查用户已经注册的设备，不允许重复注册
	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
	}
	options, sessionData, err := glb.Auth.BeginRegistration(user, registerOptions)
	if err != nil {
		klog.Warn("SESSION ERROR: " + err.Error())
		c.JSON(http.StatusInternalServerError, utils.H{
			"failed": true,
			"msg":    "server is broken",
		})
		return
	}

	// store the sessionData values
	session := sessions.Default(c)
	session.Set(REGISTER_SESSION_DATA_KEY, *sessionData)
	session.Set("a", 1)
	err = session.Save()
	if err != nil {
		msg := "can not save session"
		klog.Warn(msg + err.Error())
		c.JSON(http.StatusInternalServerError, utils.H{
			"failed": true,
			"msg":    msg,
		})
		return
	}

	// return the options generated
	c.JSON(http.StatusOK, utils.H{
		"options": options,
		"user": map[string]any{
			"id": user.ID,
		},
	})
	// options.publicKey contain our registration options
}

func RegisterFinish(ctx context.Context, c *app.RequestContext) {

	idString, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, utils.H{
			"msg": "must be registering",
		})
		return
	}
	idUint, err := strconv.ParseInt(idString, 10, 64)
	users, err := _userUsecase.MGetUsers(ctx, []int64{idUint})
	if err != nil || len(users) == 0 {
		klog.Warn("user not exists" + err.Error())
		c.JSON(http.StatusBadRequest, utils.H{
			"msg": "user not exists " + idString,
		})
		return
	}
	user := users[0]

	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	session := sessions.Default(c)

	sessionData := session.Get(REGISTER_SESSION_DATA_KEY)
	if sessionData == nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"msg": "get sessionData failed",
		})
		return
	}
	klog.Info("RegisterBegin sessionData.UserID: " + sessionData.(webauthn.SessionData).Challenge)

	response, err := protocol.ParseCredentialCreationResponseBody(c.Request.BodyStream())
	klog.Info("\nRegisterBegin request origin: " + response.Response.CollectedClientData.Origin + "\n")
	klog.Info("\nRegisterBegin allowed origin: " + strings.Join(glb.Auth.Config.RPOrigins, "|") + "\n")

	if err != nil {
		klog.Info("FinishRegistration error")
		c.JSON(http.StatusBadRequest, utils.H{
			"msg":      "get sessionData failed",
			"verified": false,
		})
		return
	}

	// parsedResponse, err := protocol.ParseCredentialRequestResponseBody(c.Request.Body)
	credential, err := glb.Auth.CreateCredential(user, sessionData.(webauthn.SessionData), response)

	// Handle validation or input errors
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"msg":      "get credential failed: " + err.Error(),
			"verified": false,
		})
		return
	}

	klog.Info("FinishRegistration success!!!")

	// If login was successful, handle next steps
	if _, err := _webAuthnCredentialUsecase.CreateCredential(int64(user.ID), credential); err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"msg":      "Save WebAuthnCredential Failed",
			"verified": true,
		})
		return
	}

	klog.Info("RegisterBegin sessionData.UserID: " + string(user.WebAuthnID()))
	c.JSON(http.StatusOK, utils.H{
		"msg":      "Register Success",
		"verified": true,
	})
}

func LoginBegin(ctx context.Context, c *app.RequestContext) {
	userName, ok := c.Params.Get("username")
	if !ok {
		c.JSON(http.StatusBadRequest, utils.H{
			"failed": true,
			"msg":    "muse give a username",
		})
		return
	}
	klog.Info("LoginBegin username: " + userName)

	users, count, err := _userUsecase.QueryUser(ctx, &userName, 10, 0)
	if err != nil || count == 0 {
		c.JSON(http.StatusBadRequest, utils.H{
			"failed": true,
			"msg":    "user not exists",
		})
		return
	}
	user := users[0]
	user.Credentials = _webAuthnCredentialUsecase.QueryCredential(int64(user.ID))

	// generate PublicKeyCredentialRequestOptions, session data
	options, sessionData, err := glb.Auth.BeginLogin(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"failed": true,
			"msg":    "login begin failed",
		})
		return
	}

	// store session data as marshaled JSON
	session := sessions.Default(c)
	session.Set(LOGIN_SESSION_DATA_KEY, *sessionData)
	if err = session.Save(); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"failed": true,
			"msg":    "can't set session data",
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"failed":  false,
		"options": options,
		"user": map[string]any{
			"id": user.ID,
		},
	})
}

func LoginFinish(ctx context.Context, c *app.RequestContext) {

	idString, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, utils.H{
			"msg":      "must be registering",
			"verified": false,
		})
		return
	}
	idUint, err := strconv.ParseInt(idString, 10, 64)
	users, err := _userUsecase.MGetUsers(ctx, []int64{idUint})
	if err != nil || len(users) == 0 {
		c.JSON(http.StatusBadRequest, utils.H{
			"msg":      "user not exists: " + idString,
			"verified": false,
		})
		return
	}
	user := users[0]
	user.Credentials = _webAuthnCredentialUsecase.QueryCredential(int64(user.ID))

	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	session := sessions.Default(c)

	sessionData := session.Get(LOGIN_SESSION_DATA_KEY)
	if sessionData == nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"msg":      "get sessionData failed",
			"verified": false,
		})
		return
	}
	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(c.Request.BodyStream())

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"msg":      "bad request body data, need CredentialRequestResponse",
			"verified": false,
		})
	}
	_, err = glb.Auth.ValidateLogin(user, sessionData.(webauthn.SessionData), parsedResponse)

	// Handle validation or input errors
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"msg":      "validate login failed: " + err.Error(),
			"verified": false,
		})
		klog.Info(err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"msg":      "Login Success",
		"verified": true,
	})
}
