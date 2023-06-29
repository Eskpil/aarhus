package helpers

import "crypto/rand"

type StringType int64

const (
	StringTypeUnknown StringType = iota
	StringTypeAlnum
	StringTypeAlpha
	StringTypeNumber
)

func RandomString(strSize int, randType StringType) string {
	var dictionary string

	if randType == StringTypeUnknown {
		return ""
	}

	if randType == StringTypeAlnum {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == StringTypeAlpha {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == StringTypeNumber {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}
