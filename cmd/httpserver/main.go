package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/brayanMuniz/tcp-to-http/internal/request"
	"github.com/brayanMuniz/tcp-to-http/internal/response"
	"github.com/brayanMuniz/tcp-to-http/internal/server"
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
	if req.RequestLine.RequestTarget == "/give400" {
		handler400(w, req)
		return
	}

	if req.RequestLine.RequestTarget == "/give500" {
		handler500(w, req)
		return
	}

	handler200(w, req)
	return
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.OK)

	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>すごい！</p>
</body>
</html>
`)

	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)

	return
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.BAD_REQUEST)

	body := []byte(`
	<html>
	  <head>
	    <title>400 Bad Request</title>
	  </head>
	  <body>
	    <h1>Bad Request</h1>
	    <p>にがい</p>
	  </body>
	</html>
	`)

	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)

	return
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.INTERNAL_SERVER_ERROR)

	body := []byte(`
	<html>
	  <head>
	    <title>500 Internal Server Error</title>
	  </head>
	  <body>
	    <h1>Internal Server Error</h1>
	    <p>すいません</p>
	  </body>
	</html>
	`)

	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)

	return
}
