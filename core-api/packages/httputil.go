package pkg

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validationMessages = map[string]string{
	"required": "field is required",
	"email":    "must be a valid email address",
	"gt":       "must be greater than 0",
}

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func ValidationErrors(err error) (map[string]string, bool) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil, false
	}
	fields := make(map[string]string, len(errs))
	for _, e := range errs {
		msg, exists := validationMessages[e.Tag()]
		if !exists {
			msg = "invalid value"
		}
		fields[e.Field()] = msg
	}
	return fields, true
}
