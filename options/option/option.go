package option

type Option[T any] struct {
	value  T
	isSome bool
}

func None[T any]() Option[T] {
	return Option[T]{value: zero[T](), isSome: false}
}

func Some[T any](t T) Option[T] {
	return Option[T]{value: t, isSome: true}
}

func zero[T any]() T {
	var zero T
	return zero
}
