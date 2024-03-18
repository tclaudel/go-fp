package fp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOk(t *testing.T) {
	tests := []struct {
		name          string
		value         any
		expectedValue any
	}{
		{
			name:          "string",
			value:         "test",
			expectedValue: "test",
		},
		{
			name:          "int",
			value:         1,
			expectedValue: 1,
		},
		{
			name:          "float",
			value:         1.0,
			expectedValue: 1.0,
		},
		{
			name:          "bool",
			value:         true,
			expectedValue: true,
		},
		{
			name:          "struct",
			value:         struct{ Name string }{Name: "test"},
			expectedValue: struct{ Name string }{Name: "test"},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			v := Ok(test.value)

			assert.True(t, v.IsOk())
			assert.False(t, v.IsErr())

			value, err := v.Unwrap()
			assert.Equal(t, test.expectedValue, value)
			assert.Nil(t, err)

			assert.Nil(t, v.UnwrapErr())
		})
	}
}

func TestErr(t *testing.T) {
	result := Err[string](errors.New("test"))

	assert.False(t, result.IsOk())
	assert.True(t, result.IsErr())

	value, err := result.Unwrap()
	assert.Equal(t, "", value)
	assert.Equal(t, errors.New("test"), err)

	assert.Equal(t, errors.New("test"), result.UnwrapErr())
}

func TestUnwrapOr(t *testing.T) {
	tests := []struct {
		name          string
		result        Result[string]
		defaultValue  string
		expectedValue any
	}{
		{
			name:          "string",
			result:        Ok("test"),
			defaultValue:  "default",
			expectedValue: "test",
		},
		{
			name:          "string",
			result:        Err[string](errors.New("test")),
			defaultValue:  "default",
			expectedValue: "default",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			value := test.result.UnwrapOr(test.defaultValue)
			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func TestUnwrapOrElse(t *testing.T) {
	tests := []struct {
		name          string
		result        Result[string]
		defaultValue  string
		expectedValue any
	}{
		{
			name:          "string",
			result:        Ok("test"),
			defaultValue:  "default",
			expectedValue: "test",
		},
		{
			name:          "string",
			result:        Err[string](errors.New("test")),
			defaultValue:  "default",
			expectedValue: "default",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			value := test.result.UnwrapOrElse(func() string {
				return test.defaultValue
			})
			assert.Equal(t, test.expectedValue, value)
		})
	}
}
