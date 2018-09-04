package controller

import (
	"time"

	mw "burgundy/api/middleware"
	"burgundy/conf"
	_Repo "burgundy/repository"
	"burgundy/service"
	"burgundy/util"

	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

type (
	// BurgundyStatus for common response status
	BurgundyStatus struct {
		TRID       string      `json:"trID"`
		ResultCode string      `json:"resultCode"`
		ResultMsg  string      `json:"resultMsg"`
		ResultData interface{} `json:"resultData,omitempty"`
	}
)

var mlog *zap.SugaredLogger

func init() {
	mlog, _ = util.InitLog("controller", "console")
}

// InitHandler ...
func InitHandler(burgundy *conf.ViperConfig, e *echo.Echo, db *gorm.DB) (err error) {

	mlog, _ = util.InitLog("controller", burgundy.GetString("logmode"))
	timeout := time.Duration(burgundy.GetInt("timeout")) * time.Second

	// Default Group
	api := e.Group("/api")
	api.File("/swagger.json", "swagger.json")
	ver := api.Group("/v1")
	sys := ver.Group("/acct")
	sys.Use(mw.TransID())
	user := sys.Group("/user")

	userRepo := _Repo.NewGormUserRepository(db)
	userSvc, err := service.NewUserService(burgundy, userRepo, timeout)
	if err != nil {
		return errors.Annotatef(err, "InitHandler")
	}
	newUserHTTPHandler(user, userSvc, burgundy.GetString("jwt_access_key"))

	return nil
}

// HTTPUserHandler ...
type HTTPUserHandler struct {
	UserService service.UserService
}

func newUserHTTPHandler(eg *echo.Group, us service.UserService, jwtkey string) {
	handler := &HTTPUserHandler{
		UserService: us,
	}

	// POST /api/v1/acct/user
	eg.POST("", handler.CreateUser)
	// POST /api/v1/acct/user/signin
	eg.POST("/signin", handler.Login)

	// METHOD /api/v1/acct/user/:accountName
	r := eg.Group("/:accountName")

	if jwtkey != "" {
		r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
			SigningKey:  []byte(jwtkey),
			TokenLookup: "header:Authorization",
		}))
	}

	r.GET("", handler.GetUser)
	r.DELETE("", handler.DeleteUser)

	r.POST("/confirmEmail", handler.ConfirmEmail)
	r.DELETE("/revokeEmail", handler.RevokeEmail)

	r.POST("/newOTP", handler.NewOTP)
	r.DELETE("/revokeOTP", handler.RevokeOTP)
	r.POST("/validateOTP", handler.ValidateOTP)
}

func response(c echo.Context, code int, trID, errCode, errMsg string, result ...interface{}) error {
	res := BurgundyStatus{
		TRID:       trID,
		ResultCode: errCode,
		ResultMsg:  errMsg,
	}

	if result != nil {
		res.ResultData = result[0]
	}
	return c.JSON(code, res)
}
