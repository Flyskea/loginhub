package passport

import (
	"github.com/gin-gonic/gin"

	apiv1 "loginhub/api/v1/passport"
	"loginhub/internal/base/handler"
)

// @Summary      获取oauth2重定向url
// @Description  获取oauth2重定向url
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Param        provider  path      string                                                       true  "第三方平台"
// @Success      200       {object}  handler.RespBody{data=passport.GetOauthRedirectURLResponse}  "请求成功"
// @Failure      400       {object}  handler.RespBody{}                                           "参数有误"
// @Failure      500       {object}  handler.RespBody{}                                           "服务器内部错误"
// @Router       /passport/oauth2/redirect/{provider} [GET]
func (h *PassportHandler) OauthRedirectURL(c *gin.Context) {
	req := &apiv1.GetOauthRedirectURLRequest{}
	if err := c.ShouldBindUri(req); handler.HandleRequestError(c, err) {
		return
	}
	handler.HandleResponse(c, nil, nil)
}

// @Summary      创建oauth2第三方平台
// @Description  创建oauth2第三方平台
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Param        body  body      passport.CreateOauth2ProviderRequest  true  "创建oauth2第三方平台请求"
// @Success      200   {object}  handler.RespBody{}                    "请求成功"
// @Failure      400   {object}  handler.RespBody{}                    "参数有误"
// @Failure      500   {object}  handler.RespBody{}                    "服务器内部错误"
// @Router       /passport/oauth2/provider [POST]
func (h *PassportHandler) CreateOauth2Provider(c *gin.Context) {
	req := &apiv1.CreateOauth2ProviderRequest{}
	if err := c.ShouldBindJSON(req); handler.HandleRequestError(c, err) {
		return
	}
	handler.HandleResponse(c, nil, nil)
}

// @Summary      列出oauth2第三方平台
// @Description  列出oauth2第三方平台
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Success      200  {object}  handler.RespBody{}  "请求成功"
// @Failure      400  {object}  handler.RespBody{}  "参数有误"
// @Failure      500  {object}  handler.RespBody{}  "服务器内部错误"
// @Router       /passport/oauth2/provider [GET]
func (h *PassportHandler) ListOauth2Provider(c *gin.Context) {
	handler.HandleResponse(c, nil, nil)
}

// @Summary      更新oauth2第三方平台
// @Description  更新oauth2第三方平台
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Param        provider  path      string                                true  "第三方平台"
// @Param        body      body      passport.UpdateOauth2ProviderRequest  true  "更新oauth2第三方平台请求"
// @Success      200       {object}  handler.RespBody{}                    "请求成功"
// @Failure      400       {object}  handler.RespBody{}                    "参数有误"
// @Failure      500       {object}  handler.RespBody{}                    "服务器内部错误"
// @Router       /passport/oauth2/provider/{provider}/update [POST]
func (h *PassportHandler) UpdateOauth2Provider(c *gin.Context) {
	providerURI := &apiv1.ProviderURI{}
	if err := c.ShouldBindUri(providerURI); handler.HandleRequestError(c, err) {
		return
	}
	req := &apiv1.UpdateOauth2ProviderRequest{}
	if err := c.ShouldBindJSON(req); handler.HandleRequestError(c, err) {
		return
	}
	handler.HandleResponse(c, nil, nil)
}

// @Summary      删除oauth2第三方平台
// @Description  删除oauth2第三方平台
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Param        provider  path      string              true  "第三方平台"
// @Success      200       {object}  handler.RespBody{}  "请求成功"
// @Failure      400       {object}  handler.RespBody{}  "参数有误"
// @Failure      500       {object}  handler.RespBody{}  "服务器内部错误"
// @Router       /passport/oauth2/provider/{provider}/delete [POST]
func (h *PassportHandler) DeleteOauth2Provider(c *gin.Context) {
	req := &apiv1.DeleteOauth2ProviderRequest{}
	if err := c.ShouldBindUri(req); handler.HandleRequestError(c, err) {
		return
	}
	handler.HandleResponse(c, nil, nil)
}
