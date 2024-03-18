package fp

type Result[T any] struct {
	value *T
	err   error
}

func Ok[T any](v T) Result[T] {
	return Result[T]{value: &v}
}

func Err[T any](e error) Result[T] {
	return Result[T]{err: e}
}

func (r Result[T]) IsOk() bool {
	return r.value != nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) Unwrap() (T, error) {
	if r.value == nil {
		return *new(T), r.err
	}
	return *r.value, nil
}

func (r Result[T]) UnwrapErr() error {
	return r.err
}

func (r Result[T]) UnwrapOr(def T) T {
	if r.value == nil {
		return def
	}
	return *r.value
}

func (r Result[T]) UnwrapOrElse(fn func() T) T {
	if r.value == nil {
		return fn()
	}
	return *r.value
}
