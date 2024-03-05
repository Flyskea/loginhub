package oauth2

import (
	"net/http"

	"github.com/Flyskea/gotools/errors"
	"github.com/gin-gonic/gin"

	apiv1 "loginhub/api/v1/oauth2"
	"loginhub/internal/base/handler"
	"loginhub/internal/base/pager"
	"loginhub/internal/base/reason"
	"loginhub/internal/domain/oauth2/entity"
	"loginhub/internal/service/oauth2"
)

type OAuth2Handler struct {
	oaps *oauth2.OAuth2Service
}

func NewOAuth2ProviderHandler(
	oaps *oauth2.OAuth2Service,
) *OAuth2Handler {
	return &OAuth2Handler{
		oaps: oaps,
	}
}

// @Summary      创建 OAuth2 服务提供商
// @Description  创建 OAuth2 服务提供商
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @Param        body  body      oauth2provider.CreateProviderRequest                          true  "创建 OAuth2 服务提供商请求"
// @Success      200   {object}  handler.RespBody{data=oauth2provider.CreateProviderResponse}  "请求成功"
// @Failure      400   {object}  handler.RespBody{}                                            "参数有误"
// @Failure      500   {object}  handler.RespBody{}                                            "服务器内部错误"
// @Router       /oauth2/provider [post]
func (o *OAuth2Handler) CreateProvider(c *gin.Context) {
	var req apiv1.CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	providerInfo := req.ToProviderInfo()
	err := o.oaps.CreateProviderInfo(c.Request.Context(), providerInfo)
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}
	handler.HandleResponse(c, nil, apiv1.CreateProviderResponse{
		ProviderInfo: *providerInfo,
	})
}

// @Summary      获取 OAuth2 服务提供商
// @Description  获取 OAuth2 服务提供商
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @Param        body  query     oauth2provider.ListProviderRequest                                               true  "获取 OAuth2 服务提供商请求"
// @Success      200   {object}  handler.RespBody{data=oauth2provider.ListProviderResponse[entity.ProviderInfo]}  "请求成功"
// @Failure      400   {object}  handler.RespBody{}                                                               "参数有误"
// @Failure      500   {object}  handler.RespBody{}                                                               "服务器内部错误"
// @Router       /oauth2/providers [get]
func (o *OAuth2Handler) ListProviders(c *gin.Context) {
	req := &apiv1.ListProviderRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		handler.HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	providers, total, err := o.oaps.ListProviderInfos(c.Request.Context(), &req.PageCond)
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}
	handler.HandleResponse(c, nil, apiv1.ListProviderResponse{
		PageResp: pager.PageResp[*entity.ProviderInfo]{
			Count: total,
			List:  providers,
		},
	})
}

// @Summary      更新 OAuth2 服务提供商
// @Description  更新 OAuth2 服务提供商
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @Param        body  body      oauth2provider.UpdateProviderRequest                          true  "更新 OAuth2 服务提供商请求"
// @Success      200   {object}  handler.RespBody{data=oauth2provider.UpdateProviderResponse}  "请求成功"
// @Failure      400   {object}  handler.RespBody{}                                            "参数有误"
// @Failure      500   {object}  handler.RespBody{}                                            "服务器内部错误"
// @Router       /oauth2/provider/{id}/update [post]
func (o *OAuth2Handler) UpdateProvider(c *gin.Context) {
	var id apiv1.ProviderURI
	if err := c.ShouldBindUri(&id); err != nil {
		handler.HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	var req apiv1.UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	providerInfo := req.ToProviderInfo()
	providerInfo.ID = id.ID
	err := o.oaps.UpdateProviderInfo(c.Request.Context(), providerInfo)
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}
	handler.HandleResponse(c, nil, apiv1.UpdateProviderResponse{
		ProviderInfo: *providerInfo,
	})
}

// @Summary      删除 OAuth2 服务提供商
// @Description  删除 OAuth2 服务提供商
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @Param        id   path      int                 true  "删除 OAuth2 服务提供商请求"
// @Success      200  {object}  handler.RespBody{}  "请求成功"
// @Failure      400  {object}  handler.RespBody{}  "参数有误"
// @Failure      500  {object}  handler.RespBody{}  "服务器内部错误"
// @Router       /oauth2/provider/{id}/delete
func (o *OAuth2Handler) DeleteProvider(c *gin.Context) {
	var req apiv1.DeleteProviderRequest
	if err := c.ShouldBindUri(&req); err != nil {
		handler.HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	err := o.oaps.DeleteProviderInfo(c.Request.Context(), req.ID)
	handler.HandleResponse(c, err, nil)
}

// @Summary      获取oauth2重定向url
// @Description  获取oauth2重定向url
// @Tags         Passport
// @Accept       json
// @Produce      json
// @Param        provider  path      string                                                       true  "第三方平台"
// @Success      200       {object}  handler.RespBody{data=passport.GetOauthRedirectURLResponse}  "请求成功"
// @Failure      400       {object}  handler.RespBody{}                                           "参数有误"
// @Failure      500       {object}  handler.RespBody{}                                           "服务器内部错误"
// @Router       /oauth2/redirect/{provider} [GET]
func (h *OAuth2Handler) OAuth2RequestURL(c *gin.Context) {
	req := &apiv1.GetOauthRequestURLRequest{}
	if err := c.ShouldBindUri(req); err != nil {
		handler.HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	requestURL, tmpSession, err := h.oaps.OAuth2RequestURL(c.Request.Context(), req.Provider)
	if err != nil {
		handler.HandleResponse(c, err, nil)
		return
	}

	c.SetCookie(apiv1.OAuth2SessionName, tmpSession, 0, "/", "", false, true)
	c.Redirect(http.StatusFound, requestURL)
}

func (h *OAuth2Handler) OAuth2SupportedProviders(c *gin.Context) {
	handler.HandleResponse(c, nil, apiv1.GetSupportedProvidersResponse{
		Providers: entity.GetSupportedProviders(),
	})
}
