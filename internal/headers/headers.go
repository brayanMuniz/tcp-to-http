package headers

import (
	"fmt"
	"strings"
	"unicode"
)

// NOTE: \r and \n is each one byte
const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (bytesCounsumed int, finishedParsing bool, err error) {
	dataString := string(data)

	crlfIndex := strings.Index(dataString, crlf)
	// more data needed
	if crlfIndex == -1 {
		return 0, false, nil
	}

	// done
	if crlfIndex == 0 {
		return 0, true, nil
	}

	newHeader := strings.TrimSpace(dataString[:crlfIndex])

	headerParts := strings.Split(newHeader, " ")
	if len(headerParts) > 2 || len(headerParts) < 2 {
		return 0, false, fmt.Errorf("header not formatted correctly")
	}

	headerKey := headerParts[0][:len(headerParts[0])-1] // remove the :
	headerValue := headerParts[1]

	// Check if key is valid
	for _, c := range headerKey {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			continue
		}

		var specialCharacters = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
		for _, s := range specialCharacters {
			if c == s {
				continue
			}
		}

		return 0, false, fmt.Errorf("header key is not valid")
	}

	headerKey = strings.ToLower(headerKey)
	h[headerKey] = headerValue

	return crlfIndex + 2, false, nil // + 2 for the crlf
}
