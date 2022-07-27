//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package validator

import (
	"encoding/json"
	v10 "github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

const (
	InvalidPayload = "invalid_payload"
)

type (
	Violation struct {
		Namespace string      `json:"-"`
		Field     string      `json:"-"`
		FieldJSON string      `json:"field"`
		Tag       string      `json:"error"`
		Value     interface{} `json:"value,omitempty"`
	}

	ValidationError struct {
		OriginalMessage string
		Message         string
		Violations      []Violation
	}

	Validator interface {
		Validate(v interface{}) *ValidationError
	}

	validator struct {
		validate *v10.Validate
	}
)

func New() Validator {
	return &validator{
		validate: v10.New(),
	}
}

// Validate is used to validate a struct using rules defined in the 'validate' tag.
// This method return the struct *ValidationError that contains details of the rules violation.
// ValidationError is compatible with 'error' interface and can be returned as error.
func (v *validator) Validate(val interface{}) *ValidationError {
	err := v.validate.Struct(val)
	if err == nil {
		return nil
	}

	if violationEntries, ok := err.(v10.ValidationErrors); ok {
		violations := make([]Violation, len(violationEntries))
		for i, err := range violationEntries {
			violations[i] = Violation{
				Namespace: err.Namespace(),
				Field:     err.Field(),
				FieldJSON: getJSONTag(val, err.StructField()),
				Tag:       err.Tag(),
				Value:     err.Value(),
			}
		}
		return &ValidationError{
			OriginalMessage: err.Error(),
			Message:         InvalidPayload,
			Violations:      violations,
		}
	}
	return &ValidationError{
		OriginalMessage: err.Error(),
		Message:         err.Error(),
	}
}

// Error return error data as json string
func (v *ValidationError) Error() string {
	dt, _ := json.Marshal(v)
	return string(dt)
}

// getJSONTag is user to get a json field name using reflection.
func getJSONTag(val interface{}, fieldName string) string {
	field, ok := reflect.TypeOf(val).FieldByName(fieldName)
	if !ok {
		return ""
	}

	labelTaf, hasTag := field.Tag.Lookup("label")
	if hasTag {
		return labelTaf
	}

	jsonTag, hasTag := field.Tag.Lookup("json")
	if hasTag {
		return strings.Split(jsonTag, ",")[0]
	}

	return fieldName
}
