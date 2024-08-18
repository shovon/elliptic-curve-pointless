package maybe

import "sync"

func Extract[T comparable](maybe Maybe[T]) (T, bool) {
	var result T

	if maybe == Nothing[T]() {
		return result, false
	}

	var wg sync.WaitGroup
	wg.Add(1)
	Then(maybe, func(t T) Maybe[any] {
		result = t
		wg.Done()
		return Something[any](nil)
	})
	wg.Wait()

	return result, true
}
