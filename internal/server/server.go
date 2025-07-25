package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/brayanMuniz/tcp-to-https/internal/request"
	"github.com/brayanMuniz/tcp-to-https/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		handler:  handler,
		listener: listener,
	}

	go s.listen()
	return s, nil

}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	w := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.BAD_REQUEST)
		log.Println("Error reading the request", err)
		body := []byte("Bad request fam")
		defaultHeaders := response.GetDefaultHeaders(len(body))
		w.WriteHeaders(defaultHeaders)
		w.WriteBody(body)
		return
	}

	s.handler(w, req)

	return
}
