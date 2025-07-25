package response

import (
	"fmt"
	"io"

	"github.com/brayanMuniz/tcp-to-https/internal/headers"
)

const rn = "\r\n"

type writerState int

const (
	stateStatusLine writerState = iota
	stateHeaders
	stateBody
	stateDoneWriting
)

type Writer struct {
	writerState writerState
	writer      io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerState: stateStatusLine,
		writer:      w,
	}
}

func (w *Writer) WriteStatusLine(s StatusCode) error {
	if w.writerState != stateStatusLine {
		return fmt.Errorf("Not in the correct state: status line")
	}
	defer func() { w.writerState = stateHeaders }()

	statusLine := getStatusLine(s)
	_, err := w.writer.Write([]byte(fmt.Sprintf("%s%s\r\n", statusLine, rn)))

	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != stateHeaders {
		return fmt.Errorf("Not in the correct state: headers")
	}
	defer func() { w.writerState = stateBody }()

	for k, v := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(body []byte) error {
	if w.writerState != stateBody {
		return fmt.Errorf("Not in the correct state: body")
	}
	defer func() { w.writerState = stateDoneWriting }()

	_, err := w.writer.Write(body)
	if err != nil {
		return err
	}

	return nil
}
