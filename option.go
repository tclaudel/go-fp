package fp

type Option[T any] struct {
	value T
	ok    bool
}

func Some[T any](v T) Option[T] {
	return Option[T]{
		value: v,
		ok:    true,
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		ok: false,
	}
}

func (o Option[T]) IsSome() bool {
	return o.ok
}

func (o Option[T]) IsNone() bool {
	return !o.ok
}

func (o Option[T]) Unwrap() (T, bool) {
	return o.value, o.ok
}

func (o Option[T]) UnwrapOr(fallback T) T {
	if value, ok := o.Unwrap(); ok {
		return value
	}

	return fallback
}

func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if value, ok := o.Unwrap(); ok {
		return value
	}

	return fn()
}
