package lang

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

	"loginhub/internal/base/constant"
)

var (
	langMapping = map[string]bool{
		language.Chinese.String():           true,
		language.SimplifiedChinese.String(): true,
		language.English.String():           true,
	}

	defaultLang = language.English.String()
)

func ExtractAndSetAcceptLanguage(c *gin.Context) {
	var lang string
	tag, _, err := language.ParseAcceptLanguage(c.GetHeader(constant.AcceptLanguageHeader))
	if err != nil {
		lang = defaultLang
	}
	if len(tag) > 0 {
		lang = tag[0].String()
	}
	if _, ok := langMapping[lang]; !ok {
		lang = defaultLang
	}
	ctx := context.WithValue(c.Request.Context(), constant.LangCtxKey, lang)
	c.Request = c.Request.WithContext(ctx)
}

func LangFromCtx(ctx context.Context) string {
	v := ctx.Value(constant.LangCtxKey)
	if v != nil {
		lang, ok := v.(string)
		if ok {
			return lang
		}
	}
	return constant.DefaultLanguage
}
