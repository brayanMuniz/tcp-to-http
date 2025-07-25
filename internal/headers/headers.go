package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// NOTE: \r and \n is each one byte
const rn = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (bytesCounsumed int, finishedParsing bool, err error) {
	rnIdx := bytes.Index(data, []byte(rn))

	// more data needed
	if rnIdx == -1 {
		return 0, false, nil
	}

	// done
	if rnIdx == 0 {
		return 0, true, nil
	}

	dataString := string(data)
	newHeader := strings.TrimSpace(dataString[:rnIdx])

	parts := strings.Split(newHeader, " ")
	if len(parts) < 2 || len(parts) > 2 {
		return 0, false, fmt.Errorf("Too many parts in header")
	}

	colonIdx := strings.Index(newHeader, ":")
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("Missing colon")
	}

	headerKey := strings.TrimSpace(newHeader[:colonIdx])
	headerValue := strings.TrimSpace(newHeader[colonIdx+1:])

	// Check if key is valid
	for _, c := range headerKey {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			continue
		}

		found := false
		var specialCharacters = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
		for _, s := range specialCharacters {
			if c == s {
				found = true
				break
			}
		}

		if !found {
			return 0, false, fmt.Errorf("header key is not valid")
		}

	}

	h.Set(headerKey, headerValue)

	return rnIdx + len(rn), false, nil
}

func (h Headers) GetValue(key string) (string, bool) {
	value, ok := h[strings.ToLower(key)]
	if !ok {
		return "", false
	}
	return value, ok
}

// append and seperate with a comma if it exist
func (h Headers) Set(key string, value string) {
	key = strings.ToLower(key)

	if val, ok := h[key]; ok {
		h[key] = val + ", " + value
	} else {
		h[key] = value
	}
}

func (h Headers) Remove(key string) {
	key = strings.ToLower(key)
	if _, exist := h[key]; exist {
		delete(h, key)
	}
}

func (h Headers) Override(key string, value string) {
	key = strings.ToLower(key)
	h[key] = value
}
