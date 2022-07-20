package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_Validate(t *testing.T) {
	type Example struct {
		Name     string `validate:"required,max=10" json:"name"`
		Email    string `validate:"required,email" json:"email"`
		Username string `validate:"required,min=4,max=8" json:"username"`
	}

	cases := map[string]struct {
		inputExample  Example
		expectedError *ValidationError
	}{
		"should return success": {
			inputExample: Example{
				Name:     "John Doe",
				Email:    "john.doe@tc.com.br",
				Username: "john.doe",
			},
			expectedError: nil,
		},
		"should return error on name: required": {
			inputExample: Example{
				Name:     "",
				Email:    "john.doe@tc.com.br",
				Username: "john.doe",
			},
			expectedError: &ValidationError{
				OriginalMessage: "Key: 'Example.Name' Error:Field validation for 'Name' failed on the 'required' tag",
				Message:         InvalidPayload,
				Violations: []Violation{
					{Namespace: "Example.Name", Field: "Name", FieldJSON: "name", Tag: "required", Value: ""},
				},
			},
		},
		"should return error on email: invalid email": {
			inputExample: Example{
				Name:     "John Doe",
				Email:    "invalid_email",
				Username: "john.doe",
			},
			expectedError: &ValidationError{
				OriginalMessage: "Key: 'Example.Email' Error:Field validation for 'Email' failed on the 'email' tag",
				Message:         InvalidPayload,
				Violations: []Violation{
					{Namespace: "Example.Email", Field: "Email", FieldJSON: "email", Tag: "email", Value: "invalid_email"},
				},
			},
		},
		"should return error on username: value is less than minimal": {
			inputExample: Example{
				Name:     "John Doe",
				Email:    "john.doe@tc.com.bt",
				Username: "jd",
			},
			expectedError: &ValidationError{
				OriginalMessage: "Key: 'Example.Username' Error:Field validation for 'Username' failed on the 'min' tag",
				Message:         InvalidPayload,
				Violations: []Violation{
					{Namespace: "Example.Username", Field: "Username", FieldJSON: "username", Tag: "min", Value: "jd"},
				},
			},
		},
		"should return error on all fields": {
			inputExample: Example{
				Name:     "John Doe 12345678901",
				Email:    "john.doe",
				Username: "12345678901",
			},
			expectedError: &ValidationError{
				OriginalMessage: "Key: 'Example.Name' Error:Field validation for 'Name' failed on the 'max' tag\nKey: 'Example.Email' Error:Field validation for 'Email' failed on the 'email' tag\nKey: 'Example.Username' Error:Field validation for 'Username' failed on the 'max' tag",
				Message:         InvalidPayload,
				Violations: []Violation{
					{Namespace: "Example.Name", Field: "Name", FieldJSON: "name", Tag: "max", Value: "John Doe 12345678901"},
					{Namespace: "Example.Email", Field: "Email", FieldJSON: "email", Tag: "email", Value: "john.doe"},
					{Namespace: "Example.Username", Field: "Username", FieldJSON: "username", Tag: "max", Value: "12345678901"},
				},
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			validator := New()
			assert.Equal(
				t,
				cs.expectedError,
				validator.Validate(cs.inputExample),
			)
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	var err = &ValidationError{
		OriginalMessage: "invalid payload",
		Message:         "invalid payload",
		Violations: []Violation{{
			Namespace: "UserModel",
			Field:     "Username",
			FieldJSON: "username",
			Tag:       "max",
			Value:     "username test",
		}},
	}

	expected := `{"OriginalMessage":"invalid payload","Message":"invalid payload","Violations":[{"field":"username","error":"max","value":"username test"}]}`
	assert.Equal(t, expected, err.Error())
}
