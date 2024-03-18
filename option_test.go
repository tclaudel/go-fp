package fp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSome(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "string",
			value: "test",
		},
		{
			name:  "int",
			value: 1,
		},
		{
			name:  "float",
			value: 1.0,
		},
		{
			name:  "bool",
			value: true,
		},
		{
			name: "struct",
			value: struct {
				Name string
			}{
				Name: "test",
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			opt := Some(test.value)

			assert.True(t, opt.IsSome())
			assert.False(t, opt.IsNone())

			v, ok := opt.Unwrap()
			assert.True(t, ok)
			assert.Equal(t, test.value, v)
		})
	}
}

func TestNone(t *testing.T) {
	opt := None[any]()

	assert.False(t, opt.IsSome())
	assert.True(t, opt.IsNone())

	v, ok := opt.Unwrap()
	assert.False(t, ok)
	assert.Nil(t, v)
}

func TestOption_UnwrapOr(t *testing.T) {
	test := []struct {
		name          string
		option        Option[string]
		or            string
		expectedValue any
	}{
		{
			name: "some",
			option: func() Option[string] {
				return Some("test")
			}(),
			or:            "default",
			expectedValue: "test",
		},
		{
			name: "none",
			option: func() Option[string] {
				return None[string]()
			}(),
			or:            "default",
			expectedValue: "default",
		},
	}

	for _, test := range test {
		test := test

		t.Run(test.name, func(t *testing.T) {
			v := test.option.UnwrapOr("default")

			assert.Equal(t, test.expectedValue, v)
		})
	}
}

func TestOption_UnwrapOrElse(t *testing.T) {
	tests := []struct {
		name          string
		option        Option[string]
		fn            func() string
		expectedValue any
	}{
		{
			name: "some",
			option: func() Option[string] {
				return Some("test")
			}(),
			fn: func() string {
				return "default"
			},
			expectedValue: "test",
		},
		{
			name: "none",
			option: func() Option[string] {
				return None[string]()
			}(),
			fn: func() string {
				return "default"
			},
			expectedValue: "default",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			v := test.option.UnwrapOrElse(test.fn)

			assert.Equal(t, test.expectedValue, v)
		})
	}
}
