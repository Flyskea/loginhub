package validator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/zh"
)

func RegisterZHTrans(v *validator.Validate, trans ut.Translator) (err error) {
	err = v.RegisterTranslation("emailVerify", trans, func(ut ut.Translator) error {
		return ut.Add("emailVerify", "{0} 不是有效的邮箱地址", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T("emailVerify", fe.Field())
		if err != nil {
			return fe.(error).Error()
		}
		return t
	})
	if err != nil {
		return err
	}
	return zh.RegisterDefaultTranslations(v, trans)
}

func RegisterENTrans(v *validator.Validate, trans ut.Translator) (err error) {
	err = v.RegisterTranslation("emailVerify", trans, func(ut ut.Translator) error {
		return ut.Add("emailVerify", "{0} is not a valid email address", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T("emailVerify", fe.Field())
		if err != nil {
			return fe.(error).Error()
		}
		return t
	})
	if err != nil {
		return err
	}
	return en.RegisterDefaultTranslations(v, trans)
}
