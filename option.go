package fp

type Option[T any] struct {
	value *T
}

func Some[T any](v T) Option[T] {
	return Option[T]{value: &v}
}

func None[T any]() Option[T] {
	return Option[T]{value: nil}
}

func (o Option[T]) IsSome() bool {
	return o.value != nil
}

func (o Option[T]) IsNone() bool {
	return o.value == nil
}

func (o Option[T]) Unwrap() (T, bool) {
	if o.value == nil {
		return *new(T), false
	}
	return *o.value, true
}

func (o Option[T]) UnwrapOr(def T) T {
	if o.value == nil {
		return def
	}
	return *o.value
}

func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.value == nil {
		return fn()
	}
	return *o.value
}
