package validator

import (
	"fmt"
	"loginhub/internal/base/constant"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/language"
)

type TranslatorRegister struct {
	Lang         string
	Lo           locales.Translator
	RegisterFunc func(v *validator.Validate, trans ut.Translator) (err error)
}

var allLanguageTranslators = []*TranslatorRegister{
	{Lang: language.Chinese.String(), Lo: zh.New(), RegisterFunc: RegisterZHTrans},
	{Lang: language.SimplifiedChinese.String(), Lo: zh.New()},
	{Lang: language.English.String(), Lo: en.New(), RegisterFunc: RegisterENTrans},
	{Lang: language.AmericanEnglish.String(), Lo: en.New()},
	{Lang: language.BritishEnglish.String(), Lo: en.New()},
}

var GlobalValidatorTranslator = make(map[string]ut.Translator, 0)

func InitValidator() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	_ = v.RegisterValidation("emailVerify", ValidateEmail)
	for _, t := range allLanguageTranslators {
		tran, val := getTran(t.Lo), v
		if t.RegisterFunc != nil {
			if err := t.RegisterFunc(val, tran); err != nil {
				panic(err)
			}
		}
		GlobalValidatorTranslator[t.Lang] = tran
	}
	v.RegisterTagNameFunc(func(fld reflect.StructField) (res string) {
		if jsonTag := fld.Tag.Get("json"); len(jsonTag) > 0 {
			if jsonTag == "-" {
				return ""
			}
			return jsonTag
		}
		if formTag := fld.Tag.Get("form"); len(formTag) > 0 {
			return formTag
		}
		return fld.Name
	})
}

func getTran(lo locales.Translator) ut.Translator {
	tran, ok := ut.New(lo, lo).GetTranslator(lo.Locale())
	if !ok {
		panic(fmt.Sprintf("not found translator %s", lo.Locale()))
	}
	return tran
}

func GetTranslator(lang string) ut.Translator {
	tran, ok := GlobalValidatorTranslator[lang]
	if !ok {
		return GlobalValidatorTranslator[constant.DefaultLanguage]
	}
	return tran
}
