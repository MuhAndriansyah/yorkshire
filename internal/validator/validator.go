package validator

import (
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var (
	trans ut.Translator
)

func NewValidator() *validator.Validate {
	v := validator.New()

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	t, _ := uni.GetTranslator("en")
	trans = t
	_ = en_translations.RegisterDefaultTranslations(v, trans)

	return v
}

func FormatErrors(err error) []FieldError {
	var fieldErrors []FieldError

	if err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrs {
				fieldErrors = append(fieldErrors, FieldError{
					Field:   err.Field(),
					Message: strings.ToLower(err.Translate(trans)),
				})
			}
		}
	}

	return fieldErrors
}

func ErrorMap(err error) map[string]string {
	out := make(map[string]string)
	for _, fe := range FormatErrors(err) {
		if _, ok := out[fe.Field]; ok {
			continue
		}
		out[fe.Field] = fe.Message
	}
	return out
}
