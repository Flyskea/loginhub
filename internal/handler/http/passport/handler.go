package passport

import (
	"context"

	"github.com/Flyskea/gotools/errors"
	"github.com/gin-gonic/gin"

	oauth2apiv1 "loginhub/api/v1/oauth2"
	apiv1 "loginhub/api/v1/passport"
	"loginhub/internal/base/handler"
	"loginhub/internal/base/iface"
	"loginhub/internal/base/middleware/authn"
	"loginhub/internal/base/reason"
	"loginhub/internal/domain/passport/entity"
	"loginhub/internal/service/passport"
)

type PassportHandler struct {
	ps       *passport.PassportService
	searcher iface.IP2RegionSearcher
}

func NewPassportHandler(
	ps *passport.PassportService,
	searcher iface.IP2RegionSearcher,
) (*PassportHandler, error) {
	return &PassportHandler{
		ps:       ps,
		searcher: searcher,
	}, nil
}

func (h *PassportHandler) createDevice(c *gin.Context) (*entity.Device, error) {
	device := NewDeviceEntity(c.GetHeader(UserAgentKey), c.ClientIP())
	err := h.fillLocation(c.Request.Context(), device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (h *PassportHandler) fillLocation(ctx context.Context, device *entity.Device) error {
	region, err := h.searcher.RegionBasicByIPV4(ctx, device.IP)
	if err != nil {
		return err
	}
	device.Location = region.String()
	return nil
}

// @Summary      发送电子邮箱验证码
// @Description  根据类型发送电子邮箱验证码
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Param        body  body      passport.EmailSendRequest  true  "发送邮箱验证码请求"
// @Success      200   {object}  handler.RespBody{}         "请求成功"
// @Failure      400   {object}  handler.RespBody{}         "参数有误"
// @Failure      500   {object}  handler.RespBody{}         "服务器内部错误"
// @Router       /passport/mail/send [POST]
func (h *PassportHandler) EmailCaptchaSend(c *gin.Context) {
	req := &apiv1.EmailSendRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		handler.HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	ctx := c.Request.Context()
	err := h.ps.EmailCaptcha(ctx, &passport.EmailSend{
		Email:       req.Email,
		CaptchaType: passport.CaptchaType(req.Type),
	})
	handler.HandleResponse(c, err, nil)
}

// @Summary      用户注册
// @Description  根据类型用户注册
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Param        body  body      passport.RegisterRequest                          true  "用户注册请求"
// @Success      200   {object}  handler.RespBody{data=passport.RegisterResponse}  "请求成功"
// @Failure      400   {object}  handler.RespBody{}                                "参数有误"
// @Failure      500   {object}  handler.RespBody{}                                "服务器内部错误"
// @Router       /passport/register [POST]
func (h *PassportHandler) Register(c *gin.Context) {
	info := &apiv1.RegisterRequest{}
	if err := c.ShouldBindJSON(info); err != nil {
		handler.HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	device, err := h.createDevice(c)
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}
	ctx := c.Request.Context()
	resp, err := h.ps.Register(ctx, &passport.RegisterInfo{
		RegisterType:    passport.RegisterType(info.RegisterType),
		UserName:        info.UserName,
		Email:           info.Email,
		Password:        info.Password,
		PasswordConfirm: info.PasswordConfirm,
		Code:            info.Captcha,
		Device:          device,
	})
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}
	c.SetCookie(authn.AuthNHTTPCookieName, resp.AccessToken, resp.AccessTokenTTL, "/", "", false, true)
	handler.HandleResponse(c, nil, apiv1.RegisterResponse{
		RefreshToken: resp.RefreshToken,
		User:         UserEntityToVO(resp.User),
	})
}

// @Summary      用户登录
// @Description  根据类型用户登录
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Param        body  body      passport.LoginRequest                          true  "用户登录请求"
// @Success      200   {object}  handler.RespBody{data=passport.LoginResponse}  "请求成功"
// @Failure      400   {object}  handler.RespBody{}                             "参数有误"
// @Failure      500   {object}  handler.RespBody{}                             "服务器内部错误"
// @Router       /passport/login [POST]
func (h *PassportHandler) Login(c *gin.Context) {
	req := &apiv1.LoginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		handler.HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	device, err := h.createDevice(c)
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}
	ctx := c.Request.Context()
	la := &passport.LoginAction{
		Type:   req.LoginType.ToService(),
		Device: device,
	}
	switch req.LoginType {
	case apiv1.LocalPasswordLoginType:
		la.LocalPasswordLoginAction = &passport.LocalPasswordLoginAction{
			Account:  req.Account,
			Password: req.Password,
		}
	case apiv1.OAuth2LoginType:
		oauthSessionID, err := c.Cookie(oauth2apiv1.OAuth2SessionName)
		if err != nil {
			handler.HandleResponse(c, errors.BadRequest(reason.OAuth2StateEmpty), nil)
			return
		}
		la.OAuth2LoginAction = &passport.OAuth2LoginAction{
			Provider:        req.Provider,
			RequestState:    req.State,
			Code:            req.Code,
			OAuth2SessionID: oauthSessionID,
		}
	}
	resp, err := h.ps.Login(ctx, la)
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}
	c.SetCookie(authn.AuthNHTTPCookieName, resp.AccessToken, resp.AccessTokenTTL, "/", "", false, true)
	handler.HandleResponse(c, nil, apiv1.LoginResponse{
		RefreshToken: resp.RefreshToken,
		User:         UserEntityToVO(resp.User),
	})
}

// @Summary      用户注销
// @Description  用户注销
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Success      200  {object}  handler.RespBody{}  "请求成功"
// @Failure      400  {object}  handler.RespBody{}  "参数有误"
// @Failure      500  {object}  handler.RespBody{}  "服务器内部错误"
// @Router       /passport/logout [POST]
func (h *PassportHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	info := authn.UserInfo(ctx)
	err := h.ps.Logout(ctx, &info)
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}
	c.SetCookie(authn.AuthNHTTPCookieName, "", -1, "/", "", false, true)
	handler.HandleResponse(c, nil, nil)
}

// @Summary      用户刷新会话
// @Description  刷新会话，返回新的refresh_token并设置cookie
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Param        body  body      passport.SessionRefreshRequest                          true  "刷新会话请求"
// @Success      200   {object}  handler.RespBody{data=passport.SessionRefreshResponse}  "请求成功"
// @Failure      400   {object}  handler.RespBody{}                                      "参数有误"
// @Failure      500   {object}  handler.RespBody{}                                      "服务器内部错误"
// @Router       /passport/session/refresh [POST]
func (h *PassportHandler) RefreshCookie(c *gin.Context) {
	req := &apiv1.SessionRefreshRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		handler.HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	ctx := c.Request.Context()
	resp, err := h.ps.Refresh(ctx, req.RefreshToken)
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}
	c.SetCookie(authn.AuthNHTTPCookieName, resp.AccessToken, resp.AccessTokenTTL, "/", "", false, true)
	handler.HandleResponse(c, nil, apiv1.SessionRefreshResponse{
		RefreshToken: resp.RefreshToken,
	})
}

// @Summary      用户登录设备
// @Description  获取用户登录的设备列表
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Success      200  {object}  handler.RespBody{data=[]passport.LoginDevice}  "请求成功"
// @Failure      400  {object}  handler.RespBody{}                             "参数有误"
// @Failure      500  {object}  handler.RespBody{}                             "服务器内部错误"
// @Router       /passport/device [GET]
func (h *PassportHandler) LoginDevices(c *gin.Context) {
	ctx := c.Request.Context()
	info := authn.UserInfo(ctx)
	resp, err := h.ps.LoginDevices(ctx, info.UserID)
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}
	handler.HandleResponse(c, nil, DevicesEntityToVO(resp))
}

// @Summary      用户注销设备
// @Description  用户注销设备
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Param        id   path      int                 true  "设备ID"
// @Success      200  {object}  handler.RespBody{}  "请求成功"
// @Failure      400  {object}  handler.RespBody{}  "参数有误"
// @Failure      500  {object}  handler.RespBody{}  "服务器内部错误"
// @Router       /passport/device/{id}/kick [POST]
func (h *PassportHandler) KickDevice(c *gin.Context) {
	req := &apiv1.KickDeviceRequest{}
	if err := c.ShouldBindUri(req); err != nil {
		handler.HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	ctx := c.Request.Context()
	info := authn.UserInfo(ctx)
	err := h.ps.DeleteDevice(ctx, &passport.DeleteDevice{
		DeviceID: req.DeviceID,
		UserInfo: info,
	})
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}
	handler.HandleResponse(c, nil, nil)
}
