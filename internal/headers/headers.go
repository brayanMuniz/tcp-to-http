package headers

import (
	"fmt"
	"strings"
)

const crlf = "\r\n" // This seperates headers

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (bytesCounsumed int, finishedParsing bool, err error) {
	dataString := string(data)

	crlfIndex := strings.Index(dataString, crlf)
	if crlfIndex == -1 {
		return 0, false, nil
	}

	// done finding heders
	if crlfIndex == 0 {
		return 0, true, nil
	}

	newHeader := strings.TrimSpace(dataString[:crlfIndex])

	headerParts := strings.Split(newHeader, " ")
	if len(headerParts) > 2 || len(headerParts) < 2 {
		return 0, false, fmt.Errorf("header not formatted correctly")
	}

	headerKey := headerParts[0][:len(headerParts[0])-1]
	headerValue := headerParts[1]

	h[headerKey] = headerValue

	return crlfIndex + 2, false, nil

}
