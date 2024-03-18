package fp

import "errors"

func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.IsErr() {
		return r
	}
	return Ok(fn(*r.value))
}

func (r Result[T]) Bind(fn func(T) Result[T]) Result[T] {
	if r.IsErr() {
		return r
	}
	return fn(*r.value)
}

func (r Result[T]) Match(ok func(T), err func(error)) {
	if r.IsErr() {
		err(r.err)
		return
	}
	ok(*r.value)
}

func (r Result[T]) BindFuncs(fns ...func(T) Result[T]) Result[T] {
	if r.IsErr() {
		return r
	}

	for _, fn := range fns {
		newResult := r.Bind(fn)
		if newResult.IsErr() {
			return newResult
		}

		r = newResult
	}

	return r
}

var ErrPredicate = errors.New("predicate is false")

func (r Result[T]) Filter(predicate func(T) bool) Result[T] {
	if r.IsErr() {
		return r
	}

	if !predicate(*r.value) {
		return Err[T](ErrPredicate)
	}

	return r
}

func (r Result[T]) FilterOr(def T, predicate func(T) bool, err error) Result[T] {
	if r.IsErr() {
		return r
	}

	if !predicate(*r.value) {
		return Err[T](err)
	}

	return r
}
func (r Result[T]) Fold(results []Result[T], fn func(T, T) T) Result[T] {
	accumulator := *r.value
	for _, r := range results {
		if r.err != nil {
			return Err[T](r.err) // Propagate the error
		}
		accumulator = fn(accumulator, *r.value)
	}
	return Ok[T](accumulator)
}

func Sequence[T any](results []Result[T]) Result[[]T] {
	var values []T
	for _, r := range results {
		if r.err != nil {
			return Err[[]T](r.err) // Propagate the error
		}
		values = append(values, *r.value)
	}
	return Ok(values)
}

func Traverse[T any, U any](values []T, fn func(T) Result[U]) Result[[]U] {
	var results []Result[U]
	for _, v := range values {
		results = append(results, fn(v))
	}
	return Sequence[U](results)
}

func Join[T any](results Result[Result[T]]) Result[T] {
	if results.IsErr() {
		return Err[T](results.err)
	}
	return *results.value
}

func Apply[T any, U any](fn Result[func(T) Result[U]], arg Result[T]) Result[U] {
	fnUnwrapped, err := fn.Unwrap()
	if err != nil {
		return Err[U](err)
	}

	return fnUnwrapped(*arg.value)
}
