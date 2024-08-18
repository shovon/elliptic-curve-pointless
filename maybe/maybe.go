package maybe

type Maybe[T any] struct {
	exists bool
	value  T
}

func Something[T any](value T) Maybe[T] {
	return Maybe[T]{exists: true, value: value}
}

func Nothing[T any]() Maybe[T] {
	return Maybe[T]{}
}

func Then[T comparable, V comparable](m Maybe[T], cb func(value T) Maybe[V]) Maybe[V] {
	if m == Nothing[T]() {
		return Nothing[V]()
	}

	return cb(m.value)
}
