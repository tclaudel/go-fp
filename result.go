package fp

type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](v T) Result[T] {
	return Result[T]{
		value: v,
		err:   nil,
	}
}

func Err[T any](e error) Result[T] {
	return Result[T]{err: e}
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

func (r Result[T]) UnwrapErr() error {
	return r.err
}

func (r Result[T]) UnwrapOr(def T) T {
	if r.err == nil {
		return r.value
	}

	return def
}

func (r Result[T]) UnwrapOrElse(fn func() T) T {
	if r.err == nil {
		return r.value
	}

	return fn()
}
