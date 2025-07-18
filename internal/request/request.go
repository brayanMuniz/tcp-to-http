package request

import (
	"fmt"
	"io"
	"strings"
)

const bufferSize = 8
const rn = "\r\n"

type State int

const (
	initialized State = iota
	done
)

type Request struct {
	RequestLine  RequestLine
	CurrentState State
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var request Request
	request.CurrentState = initialized

	readToIdx := 0
	buffer := make([]byte, bufferSize, bufferSize)

	for request.CurrentState != done {

		// increase buffer if it exceeded original size
		for readToIdx >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		amountRead, err := reader.Read(buffer[readToIdx:])
		if err != nil {
			if err == io.EOF {
				if request.CurrentState != done {
					return nil, fmt.Errorf("reached end while not done")
				}
			}

			return nil, err
		}

		readToIdx += amountRead
		parsedAmount, err := request.parse(buffer[:readToIdx])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[parsedAmount:])
		readToIdx -= parsedAmount
	}

	return &request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	bytesParsed, err := r.parseLine(data)
	if err != nil {
		return 0, err
	}

	return bytesParsed, nil
}

// if there is a line parse it and return the amount parsed
// if need more data to parse a line return 0, nil
func (r *Request) parseLine(data []byte) (int, error) {

	rnIdx := strings.Index(string(data), rn)
	if rnIdx != -1 {
		requestLine, amountParsed, err := parseRequestLine(data[:rnIdx])
		if err != nil {
			return 0, err
		}

		r.RequestLine = *requestLine
		r.CurrentState = done // for now since testing just the request line this works

		return amountParsed, nil
	}

	return 0, nil
}

func parseRequestLine(requestLineBytes []byte) (*RequestLine, int, error) {
	var requestLine RequestLine

	requestLineString := string(requestLineBytes)
	if strings.Contains(requestLineString, rn) {
		return nil, 0, nil
	}

	parts := strings.Split(requestLineString, " ")
	if len(parts) < 3 {
		return nil, len(requestLineString), fmt.Errorf("Not enough parts in the request line")
	}

	httpMethod := parts[0]
	isMethodValid := false
	validMethods := []string{"GET", "POST", "PUT"}
	for _, v := range validMethods {
		if parts[0] == v {
			isMethodValid = true
		}
	}
	if !isMethodValid {
		return nil, len(requestLineString), fmt.Errorf("not a valid method")
	}

	versionSlashIdx := strings.Index(parts[2], "/")
	if versionSlashIdx == -1 {
		return nil, len(requestLineString), fmt.Errorf("/ not found in http version")
	}
	httpVersion := parts[2][versionSlashIdx+1 : len(parts[2])]

	requestLine.Method = httpMethod
	requestLine.HttpVersion = httpVersion
	requestLine.RequestTarget = parts[1]

	return &requestLine, len(requestLineString), nil
}
