package controller

import (
	"burgundy/models"
	"context"
	"net/http"

	"github.com/juju/errors"
	"github.com/labstack/echo"
)

// CreateUser ..
func (h *HTTPUserHandler) CreateUser(c echo.Context) (err error) {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)

	user := &models.User{
		EmailConfirm: false,
		OTPConfirm:   false,
	}
	if err = c.Bind(user); err != nil {
		mlog.Errorw("CreateUser bind error ", "trID", trID, "req", *user, "err", err)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	mlog.Debugw("CreateUser ", "trID", trID, "req", user)

	if !user.Validate() {
		mlog.Errorw("CreateUser Invalid data", "trID", trID, "req", *user)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	saveUser, err := h.UserService.Store(ctx, user)
	if err != nil {
		mlog.Errorw("CreateUser error ", "trID", trID, "req", *user, "err", err)
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

	mlog.Debugw("GetUser ", "trID", trID, "account", accName)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	user, err := h.UserService.GetByID(ctx, accName)
	if errors.IsUserNotFound(err) {
		return response(c, http.StatusNotFound, trID, "0404", err.Error())
	} else if err != nil {
		return response(c, http.StatusInternalServerError, trID, "1000", err.Error())
	}

	return response(c, http.StatusOK, trID, "0000", "Request OK", user)
}

// DeleteUser ..
func (h *HTTPUserHandler) DeleteUser(c echo.Context) (err error) {

	trID := c.Response().Header().Get(echo.HeaderXRequestID)
	accName := c.Param("accountName")

	mlog.Debugw("DeleteUser ", "trID", trID, "account", accName)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := h.UserService.Delete(ctx, accName)
	if !result || err != nil {
		mlog.Errorw("DeleteUser error ", "trID", trID, "account", accName, "err", err)
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

// Login ...
func (h *HTTPUserHandler) Login(c echo.Context) (err error) {
	type LoginRequest struct {
		AccountName string `json:"accountName"`
		AccountHash string `json:"accountHash"`
	}

	trID := c.Response().Header().Get(echo.HeaderXRequestID)

	req := &LoginRequest{}

	if err = c.Bind(req); err != nil {
		mlog.Errorw("Login bind error ", "trID", trID, "req", *req, "err", err)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	if req.AccountName == "" || req.AccountHash == "" {
		mlog.Errorw("Login error ", "trID", trID, "accName", req.AccountName, "hash", req.AccountHash)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	mlog.Debugw("Login ", "trID", trID, "accName", req.AccountName, "accHash", req.AccountHash)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	user, err := h.UserService.Login(ctx, req.AccountName, req.AccountHash)
	if err != nil {
		mlog.Errorw("Login error ", "trID", trID, "accName", req.AccountName, "err", err)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
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

type EmailRequest struct {
	Email     string `json:"email"`
	EmailHash string `json:"emailHash"`
}

// ConfirmEmail ..
func (h *HTTPUserHandler) ConfirmEmail(c echo.Context) (err error) {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)

	accName := c.Param("accountName")
	req := &EmailRequest{}

	if err = c.Bind(req); err != nil {
		mlog.Errorw("ConfirmEmail bind error ", "trID", trID, "req", *req, "err", err)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	if accName == "" || req.Email == "" || req.EmailHash == "" {
		mlog.Errorw("ConfirmEmail error ", "trID", trID, "accName", accName, "email", req.Email, "emailHash", req.EmailHash)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	mlog.Debugw("ConfirmEmail ", "trID", trID, "accName", accName, "email", req.Email, "emailHash", req.EmailHash)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	saveUser, err := h.UserService.ConfirmEmail(ctx, accName, req.Email, req.EmailHash)
	if err != nil {
		mlog.Errorw("ConfirmEmail error ", "trID", trID, "accName", accName, "err", err)
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

// RevokeEmail ..
func (h *HTTPUserHandler) RevokeEmail(c echo.Context) (err error) {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)

	accName := c.Param("accountName")
	req := &EmailRequest{}

	if err = c.Bind(req); err != nil {
		mlog.Errorw("ConfirmEmail bind error ", "trID", trID, "req", *req, "err", err)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	if accName == "" || req.EmailHash == "" {
		mlog.Errorw("RevokeEmail error ", "trID", trID, "accName", accName, "emailHash", req.EmailHash)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	mlog.Debugw("RevokeEmail ", "trID", trID, "accName", accName, "emailHash", req.EmailHash)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	revokeUser, err := h.UserService.RevokeEmail(ctx, accName, req.Email, req.EmailHash)
	if err != nil {
		mlog.Errorw("RevokeEmail error ", "trID", trID, "accName", accName, "err", err)
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
		ResultData: revokeUser,
	})
}

func (h *HTTPUserHandler) NewOTP(c echo.Context) (err error) {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)
	accName := c.Param("accountName")

	mlog.Debugw("NewOTP ", "trID", trID, "account", accName)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	key, err := h.UserService.GenerateOTPKey(ctx, accName)
	if key == "" || err != nil {
		mlog.Errorw("NewOTP error ", "trID", trID, "account", accName, "err", err)
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
		ResultData: struct {
			AccountName string `json:"accountName"`
			OTPKey      string `json:"otpKey"`
		}{
			AccountName: accName,
			OTPKey:      key,
		},
	})
}

func (h *HTTPUserHandler) RevokeOTP(c echo.Context) (err error) {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)
	accName := c.Param("accountName")

	mlog.Debugw("RevokeOTP ", "trID", trID, "account", accName)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = h.UserService.RevokeOTP(ctx, accName)
	if err != nil {
		mlog.Errorw("RevokeOTP error ", "trID", trID, "account", accName, "err", err)
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
	})
}

func (h *HTTPUserHandler) ValidateOTP(c echo.Context) (err error) {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)
	accName := c.Param("accountName")
	code := c.FormValue("code")

	mlog.Debugw("ValidateOTP ", "trID", trID, "account", accName)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	ok, err := h.UserService.ValidateOTP(ctx, accName, code)
	if !ok {
		mlog.Errorw("ValidateOTP error ", "trID", trID, "account", accName, "err", err)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, BurgundyStatus{
		TRID:       trID,
		ResultCode: "0000",
		ResultMsg:  "Request OK",
	})
}
