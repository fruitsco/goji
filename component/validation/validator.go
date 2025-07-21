package validation

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var (
	// initMutex protects concurrent initialization of validators
	initMutex sync.Mutex
	// initializedValidators tracks which validator instances have been initialized
	initializedValidators = make(map[*validator.Validate]bool)
)

func New() *validator.Validate {
	validate := validator.New()

	Init(validate)

	return validate
}

func Init(v *validator.Validate) {
	initMutex.Lock()
	defer initMutex.Unlock()

	// Check if this validator instance has already been initialized
	if initializedValidators[v] {
		return
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
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

		// If the field has a header tag, use that.
		if headerTagName := fld.Tag.Get("header"); headerTagName != "" {
			name := strings.SplitN(headerTagName, ",", 2)[0]

			if name == "-" {
				return ""
			}

			return fmt.Sprintf("header.%s", name)
		}

		// TODO: handle other tag types
		return ""
	})

	v.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(decimal.Decimal); ok {
			return valuer.InexactFloat64()
		}
		return nil
	}, decimal.Decimal{})

	// Mark this validator instance as initialized
	initializedValidators[v] = true
}
