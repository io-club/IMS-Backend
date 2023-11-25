package ioginx

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	ioconfig "ims-server/pkg/config"
	ioconst "ims-server/pkg/consts"
	"ims-server/pkg/db"
	iologger "ims-server/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

// DefaultRoute Default route that includes common interfaces like ping
var DefaultRoute = []Route{
	{
		Func: func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		},
		private:    false,
		Methods:    []string{http.MethodGet},
		Permission: nil,
		funcName:   "ping",
		reqName:    "",
		respName:   "",
	},
}

type IOServer struct {
	router *gin.Engine
	routes []Route
}

// NewIOServer creates a new instance of the IOServer struct.
//
// This function allows for an optional custom *gin.Engine to be passed in, which can be used to implement a Gin router with specific requirements.
// However, if it is not provided, gin.Default() will be used by default.
//
// The function reads the value of the `Mode` configuration property and sets the Gin mode
// accordingly. Possible values for `Mode` are "debug", "release", and "test".
//
// The function then configures the middleware.
//
// The function returns a pointer to the newly created `IOServer` instance.
func NewIOServer(router *gin.Engine) *IOServer {
	if router == nil {
		router = gin.Default()
	}

	switch strings.ToLower(ioconfig.V.GetString("mode")) {
	case ioconst.ModeDebug.String():
		gin.SetMode(gin.DebugMode)
	case ioconst.ModeRelease.String():
		gin.SetMode(gin.ReleaseMode)
	case ioconst.ModeTest.String():
		gin.SetMode(gin.TestMode)
	}
	// TODO: Configure middleware

	return &IOServer{
		router: router,
		// Include default ping method
		routes: []Route{},
	}
}

func (s *IOServer) GetRouter(f func(*gorm.DB) error) *IOServer {
	err := f(db.NewDB())
	if err != nil {
		iologger.Panicf("db connect fail,err:%v", err)
	}
	iologger.Info("db connect success")
	return s
}

func (s *IOServer) SetRoutes(routes ...Route) *IOServer {
	for _, route := range routes {
		err := ParseRoute(&route)
		if err != nil {
			iologger.Fatalf("Service registration failed: %s", err)
		}
		s.routes = append(s.routes, route)
	}
	return s
}

func (s *IOServer) ServiceRegister() {
	r := s.router
	// Register routes
	routes := []Route{}
	routes = append(routes, s.routes...)

	// Check for duplicate routes
	if err := CheckRoutes(routes); err != nil {
		iologger.Fatalf("route verification failed: %s", err)
	}

	// Add default route
	routes = append(routes, DefaultRoute...)
	// Register routes (ignore case)
	for _, route := range routes {
		route := route
		url := "/" + strings.ToLower(route.funcName)
		fmt.Printf("Register route: %s\n", url)
		fn := func(c *gin.Context) {}
		if route.funcName != "ping" {
			fn = func(c *gin.Context) {
				ctx := context.Background()
				// Parse request parameters
				req, err := BindRequest(c, route)
				if err != nil {
					c.JSON(200, NewErr(c, err))
					return
				}

				// TODO: Permission check
				// if ioconfig.GetAuthConfig().CheckAccountType {
				// if route.Permission != nil && !route.Permission.IsEmpty() {
				//    if _, ok := route.Permission[cliams.UserType]; !ok {
				//       iolog.Info().Any("Permission", route.Permission).Any("cliams", cliams).Msg("Permission check failed")
				//       c.JSON(200, NewErr(c, ioerror.ErrNotPermitted))
				//       return
				//    }
				// }
				//}

				// Call the handler function
				rets, err := ToHandle(&ctx, req, route.Func)
				if err != nil {
					iologger.Warn("ToHandle Failed")
					c.JSON(200, NewErr(c, err))
					return
				}

				// Process the response
				res, err := BindResponse(rets)
				iologger.Info("response: %+v", res)
				if err != nil {
					c.JSON(200, NewErr(c, err))
					return
				}
				c.JSON(http.StatusOK, NewOk(c, res))
			}
		} else {
			fn = (route.Func).(func(c *gin.Context))
		}

		for _, method := range route.Methods {
			switch method {
			case http.MethodGet:
				r.GET(url, fn)
			case http.MethodPost:
				r.POST(url, fn)
			default:
				panic("Method not supported")
			}
		}
	}
}

func (s *IOServer) Run(addr string, serviceName string) {
	s.ServiceRegister()

	server := &http.Server{
		Addr:           addr,             // Server listening address
		Handler:        s.router,         // Handler for handling HTTP requests
		ReadTimeout:    10 * time.Second, // Timeout for reading the request body
		WriteTimeout:   10 * time.Second, // Timeout for writing the response body
		MaxHeaderBytes: 1 << 20,          // Max size of received request headers
	}
	// TODO: Add service registration

	// Start a goroutine to listen and wait for an interrupt signal to gracefully shut down the server
	go func() {
		iologger.Debug("Listening and serving HTTP on %s, service name: %s", addr, serviceName)
		err := server.ListenAndServe()
		if err != nil {
			iologger.Panicf("Failed to start server: %s", err)
		}
	}()
	// Wait for the interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	iologger.Info("Shutdown Server ...")
	// After receiving the interrupt signal, set a 5-second timeout to handle remaining requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		iologger.Fatalf("Failed to shutdown server: %v", err)
	}
}
