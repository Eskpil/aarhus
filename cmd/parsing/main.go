package main

import (
	"fmt"
	"strings"
)

func main() {
	input := "java -Xmx%MEMORY%G -Xmc%CPU% -jar fabricmc-1.20.1.jar nogui"

	keys := extractKeys(input)

	for _, key := range keys {
		fmt.Println(key)
	}
}

func extractKeys(input string) []string {
	var keys []string
	var currentKey strings.Builder
	inKey := false

	for i := 0; i < len(input); i++ {
		if input[i] == '%' {
			if inKey {
				keys = append(keys, currentKey.String())
				currentKey.Reset()
				inKey = false
			} else {
				inKey = true
			}
		} else if inKey {
			currentKey.WriteByte(input[i])
		}
	}

	if inKey {
		keys = append(keys, currentKey.String())
	}

	return keys
}
