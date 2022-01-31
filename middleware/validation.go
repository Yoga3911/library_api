package middleware

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Failed string
	Tag    string
	Value  interface{}
}

func StructValidator(models interface{}) []*Validator {
	var errors []*Validator
	
	err := validator.New().Struct(models)
	if err != nil {
		for _, i := range err.(validator.ValidationErrors) {
			var e Validator
			e.Failed = i.StructNamespace()
			e.Tag = i.Tag()
			e.Value = i.Param()
			errors = append(errors, &e)
		}
	}

	return errors
}

func InputChecker(data ...string) error {
	regex, err := regexp.Compile(`([a-zA-Z1-9@. ]+)`)
	if err != nil {
		return fmt.Errorf("cant compile regex")
	}

	for _, item := range data {
		result := regex.FindAllString(item, -1)
		if item != result[0] {
			return fmt.Errorf("forbidden input")
		}
	}

	return nil
}