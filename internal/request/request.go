package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/brayanMuniz/tcp-to-https/internal/headers"
)

const bufferSize = 8
const rn = "\r\n"

type State int

const (
	parsingRequestLine State = iota
	parsingHeaders
	parsingBody
	parsingDone
)

type Request struct {
	RequestLine  RequestLine
	Headers      headers.Headers
	Body         []byte
	bodyRead     int // this is to keep asking for more data while we read the body
	CurrentState State
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var request Request
	request.Headers = headers.NewHeaders()
	request.CurrentState = parsingRequestLine

	readToIdx := 0
	buffer := make([]byte, bufferSize, bufferSize)

	for request.CurrentState != parsingDone {

		// increase buffer if it exceeded original size
		for readToIdx >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		amountRead, err := reader.Read(buffer[readToIdx:])
		if err != nil {
			if err == io.EOF {
				if request.CurrentState != parsingDone {
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
	bytesParsed := 0
	for r.CurrentState != parsingDone {
		amountParsed, err := r.parseLine(data[bytesParsed:])
		if err != nil {
			return 0, err
		}

		bytesParsed += amountParsed
		if amountParsed == 0 {
			break
		}

	}

	return bytesParsed, nil
}

func (r *Request) parseLine(data []byte) (amountParsed int, err error) {
	switch r.CurrentState {
	case parsingRequestLine:
		requestLine, amountParsed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		// need more data
		if amountParsed == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.CurrentState = parsingHeaders

		return amountParsed + len(rn), nil

	case parsingHeaders:
		bytesConsumed, finishedParsing, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if finishedParsing {
			r.CurrentState = parsingBody
			return bytesConsumed + len(rn), nil
		}

		// need more data
		if bytesConsumed == 0 {
			return 0, nil
		}

		return bytesConsumed, nil

	case parsingBody:
		l, ok := r.Headers.GetValue("Content-Length")

		// no content length provided, assuming done state
		if !ok {
			r.CurrentState = parsingDone
			return len(data), nil
		}

		cLength, err := strconv.Atoi(l)
		if err != nil {
			return 0, fmt.Errorf("content length value is not a number")
		}

		r.Body = append(r.Body, data...)
		r.bodyRead += len(data)

		if r.bodyRead > cLength {
			return 0, fmt.Errorf("data is more than content length")
		}

		if r.bodyRead == cLength {
			r.CurrentState = parsingDone
		}

		return len(data), nil
	}

	return 0, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	rnIdx := bytes.Index(data, []byte(rn))
	if rnIdx == -1 {
		return nil, 0, nil
	}

	requestString := string(data[:rnIdx])
	requestLine, err := requestLineFromString(requestString)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, len(requestString), nil
}

func requestLineFromString(requestLineString string) (*RequestLine, error) {
	var requestLine RequestLine
	parts := strings.Split(requestLineString, " ")
	if len(parts) < 3 {
		return nil, fmt.Errorf("Not enough parts in the request line")
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
		return nil, fmt.Errorf("not a valid method")
	}

	versionSlashIdx := strings.Index(parts[2], "/")
	if versionSlashIdx == -1 {
		return nil, fmt.Errorf("/ not found in http version")
	}
	httpVersion := parts[2][versionSlashIdx+1 : len(parts[2])]

	requestLine.Method = httpMethod
	requestLine.HttpVersion = httpVersion
	requestLine.RequestTarget = parts[1]

	return &requestLine, nil

}
