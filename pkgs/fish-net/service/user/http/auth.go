package http

import (
	"IMS-Backend/pkgs/fish-net/domain"
	"IMS-Backend/pkgs/fish-net/glb"
	"IMS-Backend/pkgs/fish-net/util"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"go.uber.org/zap"
)

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
	users, err := _userUsecase.QueryUser(nil, &userName, nil, 1, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"failed": true,
			"msg":    "server is broken",
		})
		return
	}
	if len(users) == 0 {
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
	}
	users, err = _userUsecase.QueryUser(nil, &userName, nil, 1, 0)
	user := users[0]
	user.Credentials = _webAuthnCredentialUsecase.QueryCredential(int64(user.ID))
	util.PrettyLog("user.CredentialExcludeList(): ", user.CredentialExcludeList())
	// 检查用户已经注册的设备，不允许重复注册
	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
		credCreationOpts.AuthenticatorSelection.ResidentKey = protocol.ResidentKeyRequirementPreferred
	}
	options, sessionData, err := glb.Auth.BeginRegistration(user, registerOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"failed": true,
			"msg":    "server is broken",
		})
		return
	}

	// store the sessionData values
	session := sessions.Default(c)
	session.Set(REGISTER_SESSION_DATA_KEY, *sessionData)
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
	options, sessionData, err := glb.Auth.BeginDiscoverableLogin(
		webauthn.WithUserVerification(protocol.VerificationPreferred),
	)
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
	})
}

func LoginFinish(c *gin.Context) {
	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(c.Request.Body)
	if err != nil {
		glb.LOG.Info("parsedResponse error: " + err.Error())
	}
	glb.LOG.Info("parsedResponse: " + util.SPrettyLog(parsedResponse))

	if err != nil {
		glb.LOG.Info("parsedResponse failed: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":      "parse credential request response failed",
			"verified": false,
		})
		return
	}

	// Get the session data stored from the function above
	session := sessions.Default(c)
	sessionData := session.Get(LOGIN_SESSION_DATA_KEY)
	if sessionData == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":      "get sessionData failed",
			"verified": false,
		})
		return
	}

	// Validate the response
	_, err = glb.Auth.ValidateDiscoverableLogin(func(rawID, userHandle []byte) (user webauthn.User, err error) {
		// Get the user's credentials from the database
		webAuthn, err := _webAuthnCredentialUsecase.QueryByCredentialID(rawID)
		glb.LOG.Info("aaa: " + util.SPrettyLog(webAuthn))
		if err != nil {
			return nil, err
		}
		users, err := _userUsecase.MGetUsers([]int64{webAuthn.UserID})
		if err != nil {
			return nil, err
		}
		if len(users) == 0 {
			return nil, errors.New("user not exists")
		}
		domainUser := users[0]
		domainUser.Credentials = _webAuthnCredentialUsecase.QueryCredential(int64(domainUser.ID))
		return domainUser, nil
	}, sessionData.(webauthn.SessionData), parsedResponse)
	// _, err = glb.Auth.ValidateDiscoverableLogin(user, sessionData.(webauthn.SessionData), parsedResponse)

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
