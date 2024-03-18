package fp

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult_Map(t *testing.T) {
	r := Ok(1).Map(func(i int) int {
		return i + 1
	})

	value, err := r.Unwrap()
	assert.Equal(t, 2, value)
	assert.Nil(t, err)
}

func TestResult_MapErr(t *testing.T) {
	r := Err[int](assert.AnError).Map(func(i int) int {
		return i + 1
	})

	value, err := r.Unwrap()
	assert.Equal(t, 0, value)
	assert.Equal(t, assert.AnError, err)
}

var errDivideByZero = errors.New("divide by zero")

func divide(t *testing.T, divider int) func(i int) Result[int] {
	t.Helper()

	return func(i int) Result[int] {
		if divider == 0 {
			return Err[int](errDivideByZero)
		}
		return Ok(i / divider)
	}
}

func TestResult_Bind(t *testing.T) {
	tests := []struct {
		name          string
		result        Result[int]
		fn            func(i int) Result[int]
		expectedValue int
		expectedErr   error
	}{
		{
			name:          "ok",
			result:        Ok(4),
			fn:            divide(t, 2),
			expectedValue: 2,
		},
		{
			name:        "err",
			result:      Err[int](assert.AnError),
			fn:          divide(t, 2),
			expectedErr: assert.AnError,
		},
		{
			name:        "divide by zero",
			result:      Ok(1),
			fn:          divide(t, 0),
			expectedErr: errDivideByZero,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			r := tt.result.Bind(tt.fn)

			value, err := r.Unwrap()
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestResult_Match(t *testing.T) {
	tests := []struct {
		name          string
		result        Result[int]
		expectedValue int
		expectedErr   error
	}{
		{
			name:          "ok",
			result:        Ok(1),
			expectedValue: 1,
		},
		{
			name:        "err",
			result:      Err[int](assert.AnError),
			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			tt.result.Match(func(i int) {
				assert.Equal(t, tt.expectedValue, i)
			}, func(err error) {
				assert.Equal(t, tt.expectedErr, err)
			})
		})
	}
}

func TestResult_BindFuncs(t *testing.T) {
	tests := []struct {
		name          string
		result        Result[int]
		fns           []func(int) Result[int]
		expectedValue int
		expectedErr   error
	}{
		{
			name:   "ok",
			result: Ok(8),
			fns: []func(int) Result[int]{
				divide(t, 2),
				divide(t, 2),
			},
			expectedValue: 2,
			expectedErr:   nil,
		},
		{
			name:   "propagate err",
			result: Ok(8),
			fns: []func(int) Result[int]{
				divide(t, 2),
				divide(t, 0),
			},
			expectedValue: 0,
			expectedErr:   errDivideByZero,
		},
		{
			name:   "err",
			result: Err[int](assert.AnError),
			fns: []func(int) Result[int]{
				divide(t, 2),
			},
			expectedValue: 0,
			expectedErr:   assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			r := tt.result.BindFuncs(tt.fns...)

			value, err := r.Unwrap()
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestResult_Filter(t *testing.T) {
	tests := []struct {
		name          string
		result        Result[int]
		predicate     func(int) bool
		expectedValue int
		expectedErr   error
	}{
		{
			name:          "ok",
			result:        Ok(1),
			predicate:     func(i int) bool { return i > 0 },
			expectedValue: 1,
			expectedErr:   nil,
		},
		{
			name:          "predicate is false",
			result:        Ok(0),
			predicate:     func(i int) bool { return i > 0 },
			expectedValue: 0,
			expectedErr:   ErrPredicate,
		},
		{
			name:        "err",
			result:      Err[int](assert.AnError),
			predicate:   func(i int) bool { return i > 0 },
			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			r := tt.result.Filter(tt.predicate)

			value, err := r.Unwrap()
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestResult_FilterOr(t *testing.T) {
	tests := []struct {
		name          string
		result        Result[int]
		predicate     func(int) bool
		err           error
		expectedValue int
		expectedErr   error
	}{
		{
			name:          "ok",
			result:        Ok(1),
			predicate:     func(i int) bool { return i > 0 },
			err:           assert.AnError,
			expectedValue: 1,
			expectedErr:   nil,
		},
		{
			name:          "predicate is false",
			result:        Ok(0),
			predicate:     func(i int) bool { return i > 0 },
			err:           assert.AnError,
			expectedValue: 0,
			expectedErr:   assert.AnError,
		},
		{
			name:        "err",
			result:      Err[int](assert.AnError),
			predicate:   func(i int) bool { return i > 0 },
			err:         assert.AnError,
			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			r := tt.result.FilterOr(0, tt.predicate, tt.err)

			value, err := r.Unwrap()
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func Add(a, b int) int {
	return a + b
}

func TestFold(t *testing.T) {
	tests := []struct {
		name          string
		results       []Result[int]
		initialValue  int
		fn            func(int, int) int
		expectedValue int
		expectedErr   error
	}{
		{
			name:          "ok",
			results:       []Result[int]{Ok(1), Ok(2), Ok(3)},
			initialValue:  0,
			fn:            Add,
			expectedValue: 6,
			expectedErr:   nil,
		},
		{
			name:          "propagate err",
			results:       []Result[int]{Ok(1), Err[int](assert.AnError), Ok(3)},
			initialValue:  0,
			fn:            Add,
			expectedValue: 0,
			expectedErr:   assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			result := Ok(tt.initialValue)

			r := result.Fold(tt.results, tt.fn)

			value, err := r.Unwrap()
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestSequence(t *testing.T) {
	tests := []struct {
		name          string
		results       []Result[int]
		expectedValue []int
		expectedErr   error
	}{
		{
			name:          "ok",
			results:       []Result[int]{Ok(1), Ok(2), Ok(3)},
			expectedValue: []int{1, 2, 3},
			expectedErr:   nil,
		},
		{
			name:          "propagate err",
			results:       []Result[int]{Ok(1), Err[int](assert.AnError), Ok(3)},
			expectedValue: nil,
			expectedErr:   assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			result := Sequence(tt.results)

			value, err := result.Unwrap()
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestTraverse(t *testing.T) {
	tests := []struct {
		name          string
		values        []int
		expectedValue []string
		fn            func(int) Result[string]
		expectedErr   error
	}{
		{
			name:          "ok",
			values:        []int{1, 2, 3},
			expectedValue: []string{"1", "2", "3"},
			fn: func(i int) Result[string] {
				return Ok(strconv.Itoa(i))
			},
			expectedErr: nil,
		},
		{
			name:          "propagate err",
			values:        []int{1, 2, 3},
			expectedValue: nil,
			fn: func(i int) Result[string] {
				if i == 2 {
					return Err[string](assert.AnError)
				}
				return Ok(strconv.Itoa(i))
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			result := Traverse(tt.values, tt.fn)

			value, err := result.Unwrap()
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name          string
		results       Result[Result[int]]
		expectedValue int
		expectedErr   error
	}{
		{
			name:          "ok",
			results:       Ok(Ok(1)),
			expectedValue: 1,
			expectedErr:   nil,
		},
		{
			name:          "err",
			results:       Err[Result[int]](assert.AnError),
			expectedValue: 0,
			expectedErr:   assert.AnError,
		},
		{
			name:          "propagate err",
			results:       Ok(Err[int](assert.AnError)),
			expectedValue: 0,
			expectedErr:   assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			result := Join(tt.results)

			value, err := result.Unwrap()
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		name          string
		fn            Result[func(int) Result[int]]
		arg           Result[int]
		expectedValue int
		expectedErr   error
	}{
		{
			name:          "ok",
			fn:            Ok(func(i int) Result[int] { return Ok(i + 1) }),
			arg:           Ok(1),
			expectedValue: 2,
			expectedErr:   nil,
		},
		{
			name:          "err",
			fn:            Err[func(int) Result[int]](assert.AnError),
			arg:           Ok(1),
			expectedValue: 0,
			expectedErr:   assert.AnError,
		},
		{
			name:          "propagate err",
			fn:            Ok(func(i int) Result[int] { return Err[int](assert.AnError) }),
			arg:           Ok(1),
			expectedValue: 0,
			expectedErr:   assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			result := Apply(tt.fn, tt.arg)

			value, err := result.Unwrap()
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
