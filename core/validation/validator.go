package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// If the field has a protobuf tag, use that.
		if protoTagName := fld.Tag.Get("protobuf"); protoTagName != "" {
			options := strings.Split(protoTagName, ",")

			for _, option := range options {
				if strings.HasPrefix(option, "name=") {
					return strings.TrimPrefix(option, "name=")
				}
			}

			return ""
		}

		// If the field has a json tag, use that.
		if jsonTagName := fld.Tag.Get("json"); jsonTagName != "" {
			name := strings.SplitN(jsonTagName, ",", 2)[0]

			if name == "-" {
				return ""
			}

			return name
		}

		// TODO: handle other tag types
		return ""
	})

	return validate
}
