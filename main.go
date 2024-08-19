package main

import (
	"crypto/rand"
	"crypto/sha256"
	"elliptic-curve-interfaces/ecc"
	"elliptic-curve-interfaces/ecc/curves/secp256r1"
	"elliptic-curve-interfaces/ecc/ecdsa"
	"encoding/base64"
	"fmt"
	"math/big"
)

func hasher(message []byte) ([]byte, error) {
	hash := sha256.New()
	hash.Write(message)
	return hash.Sum(nil), nil
}

func main() {
	privateKey, err := rand.Int(rand.Reader, secp256r1.GeneratorOrder())
	if err != nil {
		panic(err)
	}

	publicKey := secp256r1.Generator().ScalarMultiply(privateKey)

	message := []byte("Hello, World!")
	b64Message := base64.StdEncoding.EncodeToString(message)

	fmt.Println("Message (raw string)")
	fmt.Println(string(message))
	fmt.Println()

	fmt.Println("Message (raw base64)")
	fmt.Println(b64Message)
	fmt.Println()

	r, s, err := ecdsa.Sign(hasher, secp256r1.Generator(), secp256r1.GeneratorOrder(), privateKey, message)

	if err != nil {
		panic(err)
	}

	b, err := ecc.SerializeUncompressedECPointFormat(publicKey)
	if err != nil {
		panic(err)
	}

	encodedPublicKey := base64.StdEncoding.EncodeToString(b)

	fmt.Println("Public key (raw uncompressd EC Point Format bas64)")
	fmt.Println(encodedPublicKey)
	fmt.Println()

	rawSignature := []byte{}
	rawSignature = append(rawSignature, r.Bytes()...)
	rawSignature = append(rawSignature, s.Bytes()...)

	// Convert byte slice to base64 encoded string
	encodedSignature := base64.StdEncoding.EncodeToString(rawSignature)

	fmt.Println("Signature (raw base64)")
	fmt.Println(encodedSignature)
	fmt.Println()

	verification, err := ecdsa.Verify(hasher, secp256r1.Generator(), secp256r1.GeneratorOrder(), publicKey, [2]*big.Int{r, s}, message)
	if err != nil {
		panic(err)
	}

	fmt.Println("Does verify?")
	fmt.Println(verification)
}
