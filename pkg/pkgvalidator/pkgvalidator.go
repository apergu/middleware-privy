package pkgvalidator

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var uni *ut.UniversalTranslator
var enTrans locales.Translator
var trans ut.Translator
var v *validator.Validate

func InitValidator() {

	enTrans = en.New()
	uni = ut.New(enTrans, enTrans)

	trans, _ = uni.GetTranslator("en")

	v = validator.New()

	en_translations.RegisterDefaultTranslations(v, trans)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see
}

func Validate(c interface{}) []map[string]interface{} {
	var validationErrors []map[string]interface{}

	err := v.Struct(c)
	if err != nil {
		// Validation failed, print the error messages
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors,
				map[string]interface{}{
					"field":       err.Field(),
					"description": err.Translate(trans),
				})
		}
	}

	return validationErrors
}
