package main

import (
	"github.com/brayanMuniz/tcp-to-https/internal/request"
	"github.com/brayanMuniz/tcp-to-https/internal/response"
	"github.com/brayanMuniz/tcp-to-https/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.BAD_REQUEST)
	body := []byte("Bad request fam")
	h := response.GetDefaultHeaders(len(body))
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}
