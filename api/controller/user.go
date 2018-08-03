package controller

import (
	"burgundy/models"
	"context"
	"net/http"

	"github.com/labstack/echo"
)

// CreateUser ..
func (h *HTTPUserHandler) CreateUser(c echo.Context) (err error) {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)

	user := &models.User{
		ConfirmEmail: false,
		ConfirmOTP:   false,
	}
	if err = c.Bind(user); err != nil {
		mlog.Infow("CreateUser bind error ", "trID", trID, "req", *user, "err", err)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	mlog.Infow("CreateUser ", "trID", trID, "req", user)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	saveUser, err := h.UserService.Store(ctx, user)
	if err != nil {
		mlog.Infow("CreateUser error ", "trID", trID, "req", *user, "err", err)
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

	mlog.Infow("GetUser ", "trID", trID, "account", accName)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	user, err := h.UserService.GetByID(ctx, accName)
	if err != nil {
		mlog.Infow("GetUser error ", "trID", trID, "account", accName, "err", err)
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
		ResultData: user,
	})
}

// UpdateUser ..
func (h *HTTPUserHandler) UpdateUser(c echo.Context) (err error) {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)

	user := &models.User{}
	if err = c.Bind(user); err != nil {
		mlog.Infow("UpdateUser bind error ", "trID", trID, "req", *user, "err", err)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	mlog.Infow("UpdateUser ", "trID", trID, "req", user)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	saveUser, err := h.UserService.Update(ctx, user)
	if err != nil {
		mlog.Infow("UpdateUser error ", "trID", trID, "req", *user, "err", err)
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
		ResultData: saveUser,
	})
}

// DeleteUser ..
func (h *HTTPUserHandler) DeleteUser(c echo.Context) (err error) {

	trID := c.Response().Header().Get(echo.HeaderXRequestID)
	accName := c.Param("accountName")

	mlog.Infow("DeleteUser ", "trID", trID, "account", accName)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := h.UserService.Delete(ctx, accName)
	if !result || err != nil {
		mlog.Infow("DeleteUser error ", "trID", trID, "account", accName, "err", err)
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
