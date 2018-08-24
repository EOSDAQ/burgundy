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
		EmailConfirm: false,
		OTPConfirm:   false,
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

	if !user.Validate() {
		mlog.Infow("CreateUser Invalid data", "trID", trID, "req", *user)
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
		mlog.Infow("ConfirmEmail bind error ", "trID", trID, "req", *req, "err", err)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	if accName == "" || req.Email == "" || req.EmailHash == "" {
		mlog.Infow("ConfirmEmail error ", "trID", trID, "accName", accName, "email", req.Email, "emailHash", req.EmailHash)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	mlog.Infow("ConfirmEmail ", "trID", trID, "accName", accName, "email", req.Email, "emailHash", req.EmailHash)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	saveUser, err := h.UserService.ConfirmEmail(ctx, accName, req.Email, req.EmailHash)
	if err != nil {
		mlog.Infow("ConfirmEmail error ", "trID", trID, "accName", accName, "err", err)
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
		mlog.Infow("ConfirmEmail bind error ", "trID", trID, "req", *req, "err", err)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	if accName == "" || req.EmailHash == "" {
		mlog.Infow("RevokeEmail error ", "trID", trID, "accName", accName, "emailHash", req.EmailHash)
		return c.JSON(http.StatusBadRequest, BurgundyStatus{
			TRID:       trID,
			ResultCode: "1101",
			ResultMsg:  "Invalid Parameter",
		})
	}

	mlog.Infow("RevokeEmail ", "trID", trID, "accName", accName, "emailHash", req.EmailHash)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	revokeUser, err := h.UserService.RevokeEmail(ctx, accName, req.Email, req.EmailHash)
	if err != nil {
		mlog.Infow("RevokeEmail error ", "trID", trID, "accName", accName, "err", err)
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

	mlog.Infow("NewOTP ", "trID", trID, "account", accName)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	key, err := h.UserService.GenerateOTPKey(ctx, accName)
	if key == "" || err != nil {
		mlog.Infow("NewOTP error ", "trID", trID, "account", accName, "err", err)
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

	mlog.Infow("RevokeOTP ", "trID", trID, "account", accName)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = h.UserService.RevokeOTP(ctx, accName)
	if err != nil {
		mlog.Infow("RevokeOTP error ", "trID", trID, "account", accName, "err", err)
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

	mlog.Infow("ValidateOTP ", "trID", trID, "account", accName)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	ok, err := h.UserService.ValidateOTP(ctx, accName, code)
	if !ok {
		mlog.Infow("ValidateOTP error ", "trID", trID, "account", accName, "err", err)
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
