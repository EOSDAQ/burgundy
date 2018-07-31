package controller

import (
	"burgundy/models"
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// CreateUser ..
func (h *HTTPUserHandler) CreateUser(c echo.Context) (err error) {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)

	user := new(models.User)
	if err = c.Bind(user); err != nil {
		log.Printf("[CreateUser] bind error trID[%s] req[%v] err[%s]", trID, *user, err)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	log.Printf("[CreateUser] trID[%s] req[%v]", trID, user)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	saveUser, err := h.UserService.Store(ctx, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1000",
			ResultMsg:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, BurgundyStatus{
		TRID:       trID,
		ResultCode: "0000",
		ResultMsg:  "Request OK",
		ResultData: saveUser.String(),
	})
}

// GetUser ..
func (h *HTTPUserHandler) GetUser(c echo.Context) (err error) {

	trID := c.Response().Header().Get(echo.HeaderXRequestID)
	accName := c.Param("accountName")

	log.Printf("[GetUser] accountName[%s] trID[%s]", accName, trID)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	user, err := h.UserService.GetByID(ctx, accName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1000",
			ResultMsg:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, BurgundyStatus{
		TRID:       trID,
		ResultCode: "0000",
		ResultMsg:  "Request OK",
		ResultData: user.String(),
	})
}

// DeleteUser ..
func (h *HTTPUserHandler) DeleteUser(c echo.Context) (err error) {

	trID := c.Response().Header().Get(echo.HeaderXRequestID)
	accName := c.Param("accountName")

	log.Printf("[DeleteUser] accountName[%s] trID[%s]", accName, trID)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := h.UserService.Delete(ctx, accName)
	if !result || err != nil {
		return c.JSON(http.StatusInternalServerError, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1000",
			ResultMsg:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, BurgundyStatus{
		TRID:       trID,
		ResultCode: "0000",
		ResultMsg:  "Request OK",
		ResultData: accName,
	})
}
