package primitives

// TODO: this does not exclusively represent a Weierstrass curve; it represents
// anything that cares about two constants, a and b. Find another name for this.

type ABValues[T any] interface {
	A() T
	B() T
}
