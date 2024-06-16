package pkgvalidator

import (
	"regexp"

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

// Custom validation function
func validatePath(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	regex := `^(([\w-]+/)+\d+)$`
	re := regexp.MustCompile(regex)
	return re.MatchString(path)
}

func InitValidator() {

	enTrans = en.New()
	uni = ut.New(enTrans, enTrans)

	trans, _ = uni.GetTranslator("en")

	v = validator.New()
	v.RegisterValidation("formatTopUpID", validatePath)

	en_translations.RegisterDefaultTranslations(v, trans)

	v.RegisterTranslation("formatTopUpID", trans, func(ut ut.Translator) error {
		return ut.Add("formatTopUpID", "{0} must be in a valid format, example:EID/MID/CID/001", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("formatTopUpID", fe.Field())
		return t
	})
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