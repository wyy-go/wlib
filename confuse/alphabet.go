package confuse

import (
	cryptoRand "crypto/rand"
	"math/bits"
	"math/rand/v2"
	"unsafe"
)

// Previous defined bytes, do not change this.
var (
	DefaultAlphabet   = []byte("QWERTYUIOPLKJHGFDSAZXCVBNMabcdefghijklmnopqrstuvwxyz")
	DefaultDigit      = []byte("0123456789")
	DefaultAlphaDigit = []byte("QWERTYUIOPLKJHGFDSAZXCVBNMabcdefghijklmnopqrstuvwxyz0123456789")
	DefaultSymbol     = []byte("QWERTYUIOPLKJHGFDSAZXCVBNMabcdefghijklmnopqrstuvwxyz0123456789!\"#$%&'()*+,-./:;<=>?@[\\]^_{|}~`") //nolint: lll
)

// Alphabet rand alpha with give length, which Contains only letters
func Alphabet(length int) string { return randString(length, DefaultAlphabet) }

// AlphabetBytes rand alpha with give length, which Contains only letters
func AlphabetBytes(length int) []byte { return randBytes(length, DefaultAlphabet) }

// Number rand string with give length, which Contains only number
func Number(length int) string { return randString(length, DefaultDigit) }

// NumberBytes rand string with give length, which Contains only number
func NumberBytes(length int) []byte { return randBytes(length, DefaultDigit) }

// AlphaNumber rand string with give length, which Contains number and letters
func AlphaNumber(length int) string { return randString(length, DefaultAlphaDigit) }

// AlphaNumberBytes rand string with give length, which Contains number and letters
func AlphaNumberBytes(length int) []byte { return randBytes(length, DefaultAlphaDigit) }

// Symbol rand symbol with give length, which Contains number, letters and specific symbol
func Symbol(length int) string { return randString(length, DefaultSymbol) }

// SymbolBytes rand symbol with give length, which Contains number, letters and specific symbol
func SymbolBytes(length int) []byte { return randBytes(length, DefaultSymbol) }

// String rand bytes, if not alphabets, it will use DefaultAlphabet.
func String(length int, alphabets ...byte) string {
	if len(alphabets) == 0 {
		alphabets = DefaultAlphaDigit
	}
	return randString(length, alphabets)
}

// Bytes rand bytes, if not alphabets, it will use DefaultAlphabet.
func Bytes(length int, alphabets ...byte) []byte {
	if len(alphabets) == 0 {
		alphabets = DefaultAlphaDigit
	}
	return randBytes(length, alphabets)
}

func randString(length int, alphabets []byte) string {
	b := randBytes(length, alphabets)
	return *(*string)(unsafe.Pointer(&b))
}

func randBytes(length int, alphabets []byte) []byte {
	b := make([]byte, length)
	if _, err := cryptoRand.Read(b); err == nil {
		for i, v := range b {
			b[i] = alphabets[v%byte(len(alphabets))]
		}
		return b
	}

	bn := bits.Len(uint(len(alphabets)))
	mask := int64(1)<<bn - 1
	max := 63 / bn

	// A rand.Int64() generates 63 random bits, enough for alphabets letters!
	for i, cache, remain := 0, rand.Int64(), max; i < length; {
		if remain == 0 {
			cache, remain = rand.Int64(), max
		}
		if idx := int(cache & mask); idx < len(alphabets) {
			b[i] = alphabets[idx]
			i++
		}
		cache >>= bn
		remain--
	}
	return b
}
