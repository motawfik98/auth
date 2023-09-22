package utils

import (
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"regexp"
	"strings"
)

type CustomValidator struct {
	Validator  *validator.Validate
	Translator *ut.Translator
}

func (cv *CustomValidator) TranslateErrors() {
	englishLocale := en.New()
	universalTranslator := ut.New(englishLocale, englishLocale)
	trans, _ := universalTranslator.GetTranslator("en")
	cv.Translator = &trans
	cv.Validator.RegisterValidation("strongPassword", strongPasswordCheck)
	enTranslations.RegisterDefaultTranslations(cv.Validator, *cv.Translator)

	cv.Validator.RegisterTranslation("strongPassword", trans, func(ut ut.Translator) error {
		return ut.Add("strongPassword", "Password is not strong enough", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("strongPassword", fe.Field())
		return t
	})

}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)
	if err == nil {
		return nil
	}
	var errMessages validator.ValidationErrors
	errors.As(err, &errMessages)
	var sb strings.Builder
	for _, e := range errMessages {
		sb.WriteString(e.Translate(*cv.Translator) + "\n")
	}
	return errors.New(sb.String())
}

func strongPasswordCheck(fieldLevel validator.FieldLevel) bool {
	fieldValue := fieldLevel.Field().String()
	if len(fieldValue) < 8 {
		return false
	}
	containsNumber, numberErr := regexp.MatchString("\\d", fieldValue)
	containsSpecial, specialErr := regexp.MatchString("[`!@#%^&*()_+={};\\':\"\\\\|,.<>\\[\\]?~-]", fieldValue)
	containsUpper := strings.ToLower(fieldValue) != fieldValue
	containsLower := strings.ToUpper(fieldValue) != fieldValue
	if numberErr != nil || specialErr != nil {
		return false
	}
	return containsUpper && containsLower && containsSpecial && containsNumber
}
