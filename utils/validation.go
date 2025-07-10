package utils

import (
	"strings"

	validator "github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	if err.Error() == "EOF" {
		return "please check your request body"
	}

	var errors []string
	var errorMessage string
	switch v := err.(type) {
	case validator.ValidationErrors:
		for _, e := range v {
			errors = append(errors, e.Field()+" "+e.Tag())
		}
		errorMessage = strings.Join(errors, ", ")
	default:
		parts := strings.SplitN(err.Error(), "Error:", 2)
		if len(parts) == 2 {
			errors = append(errors, strings.TrimSpace(parts[1]))
		} else {
			errors = append(errors, err.Error())
		}
		errorMessage = strings.Join(errors, ", ")
	}

	return errorMessage
}
