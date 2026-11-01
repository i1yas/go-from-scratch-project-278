package httpapi

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"hexleturlshort/internal/links"
)

// RegisterValidators registers custom validators
func RegisterValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("link_url", func(fl validator.FieldLevel) bool {
			url := fl.Field().String()
			_, err := links.NewURL(url)

			return err == nil
		})
		if err != nil {
			return err
		}

		err = v.RegisterValidation("shortcode", func(fl validator.FieldLevel) bool {
			code := fl.Field().String()
			_, err := links.NewShortCode(code)

			return err == nil
		})
		if err != nil {
			return err
		}
	}

	// NOTE: need to register this to use json names in validation errors
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}

			return name
		})
	}

	return nil
}
