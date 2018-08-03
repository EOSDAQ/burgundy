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
	"go.uber.org/zap"
)

type (
	// BurgundyStatus for common response status
	BurgundyStatus struct {
		TRID       string      `json:"trID"`
		ResultCode string      `json:"resultCode"`
		ResultMsg  string      `json:"resultMsg"`
		ResultData interface{} `json:"resultData"`
	}
)

var mlog *zap.SugaredLogger

func init() {
	mlog, _ = util.InitLog("controller", "console")
}

// InitHandler ...
func InitHandler(burgundy conf.ViperConfig, e *echo.Echo, db *gorm.DB) (err error) {

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
	newUserHTTPHandler(user, userSvc)

	return nil
}

// HTTPUserHandler ...
type HTTPUserHandler struct {
	UserService service.UserService
}

func newUserHTTPHandler(eg *echo.Group, us service.UserService) {
	handler := &HTTPUserHandler{
		UserService: us,
	}

	// /api/v1/acct/user
	eg.POST("", handler.CreateUser)
	eg.PUT("/:accountName", handler.UpdateUser)
	eg.GET("/:accountName", handler.GetUser)
	eg.DELETE("/:accountName", handler.DeleteUser)
}
