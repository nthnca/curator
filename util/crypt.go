package util

import "log"

func MD5(a []byte) [16]byte {
	var rv [16]byte
	if len(a) != 16 {
		log.Fatalf("Byte array has incorrect size for MD5 (%d)", len(a))
	}

	copy(rv[:], a)
	return rv
}

func Sha256(a []byte) [32]byte {
	var rv [32]byte
	if len(a) != 32 {
		log.Fatalf("Byte array has incorrect size for SHA256 (%d)", len(a))
	}

	copy(rv[:], a)
	return rv
}
