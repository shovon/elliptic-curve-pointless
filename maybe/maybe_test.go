package maybe

import "testing"

func TestThen(t *testing.T) {
	type something struct{}
	s := Something(something{})
	Then(s, func(v something) Maybe[something] {
		return Something(something{})
	})
}
