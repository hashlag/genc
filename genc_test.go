package main

import (
	"genc/cryptoshield"
	"testing"
)

func BenchmarkEncryptFile(b *testing.B) {
	enc := cryptoshield.NewEncryptor()

	for n := 0; n < b.N; n++ {
		enc.EncryptFile("localres/testdata.bin", "localres/testdata.enc.bin", "00000000000000000000", false)
	}
}
