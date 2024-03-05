package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"loginhub/internal/base/middleware/lang"
	"loginhub/internal/base/reason"
	myvalidator "loginhub/internal/base/validator"

	"github.com/Flyskea/gotools/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func HandleRequestError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if ve, ok := err.(validator.ValidationErrors); ok {
		if len(ve) > 0 {
			HandleResponse(c, errors.BadRequest(reason.RequestFormatError).WithMsg(
				ve[0].Translate(myvalidator.GetTranslator(lang.LangFromCtx(c.Request.Context()))),
			), nil)
			return true
		}
	}
	HandleResponse(c, errors.BadRequest(reason.RequestFormatError), nil)
	return true
}

func HandleResponse(c *gin.Context, err error, data interface{}) {
	// no error
	if err == nil {
		c.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, reason.Success, data))
		return
	}

	ctx := c.Request.Context()
	var myErr *errors.Error
	// unknown error
	if !errors.As(err, &myErr) {
		slog.ErrorContext(ctx, "unknown error",
			slog.String("error", err.Error()),
			slog.String("stack", errors.LogStack(2, 0)))
		c.JSON(http.StatusInternalServerError, NewRespBody(
			http.StatusInternalServerError, reason.UnknownError))
		return
	}

	// log internal server error
	if errors.IsInternalServer(myErr) {
		slog.ErrorContext(c, "internal server error", slog.String("error", fmt.Sprintf("%+v", myErr)))
	}

	respBody := NewRespBodyFromError(myErr)
	if data != nil {
		respBody.Data = data
	}
	c.JSON(int(myErr.Code), respBody)
}
