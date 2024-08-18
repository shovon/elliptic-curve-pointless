package main

import (
	"elliptic-curve-interfaces/ecc/curves/dummy"
	"elliptic-curve-interfaces/math/finitefield"
	"elliptic-curve-interfaces/math/primitives"
	"elliptic-curve-interfaces/maybe"
	"fmt"
	"math/big"
)

func main() {
	// privateKey, err := rand.Int(rand.Reader, secp256r1.GeneratorOrder())
	// if err != nil {
	// 	panic(err)
	// }

	// // publicKey := secp256r1.Generator().ScalarMultiply(privateKey)

	// rs, err := ecdsa.Sign(secp256r1.Generator(), secp256r1.GeneratorOrder(), privateKey, []byte("Hello, World!"))

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(rs[0])
	// fmt.Println(rs[1])

	for k := 0; k < 25; k++ {
		v := dummy.Generator().ScalarMultiply(big.NewInt(int64(k)))
		p, ok := maybe.Extract(maybe.Maybe[finitefield.NotInfinity[dummy.Curve]](v))
		fmt.Printf("%d ", k)
		if !ok {
			fmt.Println("Point at infinity")
		} else {
			fmt.Printf("%s %s\n", primitives.Point2D[*big.Int](p).X().String(), primitives.Point2D[*big.Int](p).Y().String())
		}
	}
}
