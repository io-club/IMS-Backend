package http

import (
	"encoding/binary"
	"fishnet/domain"
	"fishnet/glb"
	"fishnet/user/usecase"
	"fishnet/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"go.uber.org/zap"
)

var _userUsecase domain.UserUsecase
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

func Play(c *gin.Context) {
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
	c.JSON(200, gin.H{"count": count})
}

func RegisterBegin(c *gin.Context) {
	userName, ok := c.Params.Get("username")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "muse give a username",
		})
		return
	}
	glb.LOG.Info("RegisterBegin username: " + userName)

	// 没有用户就新建
	users, count, err := _userUsecase.QueryUser(&userName, 1, 0)
	if err != nil || count == 0 {
		glb.LOG.Info("Creating user: " + userName)
		user := &domain.User{
			Username: userName,
			Nickname: userName,
			Icon:     "https://pics.com/avatar.png",
		}
		err = _userUsecase.CreateUser([]*domain.User{user})
		if err != nil {
			glb.LOG.Warn("USER CREATION ERROR: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"failed": true,
				"msg":    "user not created",
			})
			return
		}
	} else {
		glb.LOG.Info("user exists")
	}

	users, count, err = _userUsecase.QueryUser(&userName, 1, 0)
	user := users[0]
	user.Credentials = _webAuthnCredentialUsecase.QueryCredential(int64(user.ID))
	util.PrettyLog("user.CredentialExcludeList(): ", user.CredentialExcludeList())
	// 检查用户已经注册的设备，不允许重复注册
	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
	}
	options, sessionData, err := glb.Auth.BeginRegistration(user, registerOptions)
	if err != nil {
		glb.LOG.Warn("SESSION ERROR: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
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
		glb.LOG.Warn(msg + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"failed": true,
			"msg":    msg,
		})
		return
	}

	// return the options generated
	c.JSON(http.StatusOK, gin.H{
		"options": options,
		"user": map[string]any{
			"id": user.ID,
		},
	})
	// options.publicKey contain our registration options
}

func RegisterFinish(c *gin.Context) {

	idString, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "must be registering",
		})
		return
	}
	idUint, err := strconv.ParseInt(idString, 10, 64)
	glb.LOG.Info("field to query", zap.Int64("idUint", idUint))
	users, err := _userUsecase.MGetUsers([]int64{idUint})
	if err != nil || len(users) == 0 {
		glb.LOG.Warn("user not exists" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
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
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "get sessionData failed",
		})
		return
	}
	glb.LOG.Info("RegisterBegin sessionData.UserID: " + sessionData.(webauthn.SessionData).Challenge)

	response, err := protocol.ParseCredentialCreationResponseBody(c.Request.Body)
	glb.LOG.Info("\nRegisterBegin request origin: " + response.Response.CollectedClientData.Origin + "\n")
	glb.LOG.Info("\nRegisterBegin allowed origin: " + strings.Join(glb.Auth.Config.RPOrigins, "|") + "\n")

	if err != nil {
		glb.LOG.Info("FinishRegistration error")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":      "get sessionData failed",
			"verified": false,
		})
		return
	}

	// parsedResponse, err := protocol.ParseCredentialRequestResponseBody(c.Request.Body)
	credential, err := glb.Auth.CreateCredential(user, sessionData.(webauthn.SessionData), response)

	// Handle validation or input errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":      "get credential failed: " + err.Error(),
			"verified": false,
		})
		return
	}

	glb.LOG.Info("FinishRegistration success!!!")

	// If login was successful, handle next steps
	if _, err := _webAuthnCredentialUsecase.CreateCredential(int64(user.ID), credential); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":      "Save WebAuthnCredential Failed",
			"verified": true,
		})
		return
	}

	glb.LOG.Info("RegisterBegin sessionData.UserID: " + string(user.WebAuthnID()))
	c.JSON(http.StatusOK, gin.H{
		"msg":      "Register Success",
		"verified": true,
	})
}

func LoginBegin(c *gin.Context) {
	userName, ok := c.Params.Get("username")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"failed": true,
			"msg":    "muse give a username",
		})
		return
	}
	glb.LOG.Info("LoginBegin username: " + userName)

	users, count, err := _userUsecase.QueryUser(&userName, 1, 0)
	if err != nil || count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"failed": true,
			"msg":    "login begin failed",
		})
		return
	}

	// store session data as marshaled JSON
	session := sessions.Default(c)
	session.Set(LOGIN_SESSION_DATA_KEY, *sessionData)
	if err = session.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"failed": true,
			"msg":    "can't set session data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"failed":  false,
		"options": options,
		"user": map[string]any{
			"id": user.ID,
		},
	})
}

func LoginFinish(c *gin.Context) {

	idString, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":      "must be registering",
			"verified": false,
		})
		return
	}
	idUint, err := strconv.ParseInt(idString, 10, 64)
	users, err := _userUsecase.MGetUsers([]int64{idUint})
	if err != nil || len(users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
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
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":      "get sessionData failed",
			"verified": false,
		})
		return
	}
	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":      "bad request body data, need CredentialRequestResponse",
			"verified": false,
		})
	}
	_, err = glb.Auth.ValidateLogin(user, sessionData.(webauthn.SessionData), parsedResponse)

	// Handle validation or input errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":      "validate login failed: " + err.Error(),
			"verified": false,
		})
		glb.LOG.Info(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":      "Login Success",
		"verified": true,
	})
}
