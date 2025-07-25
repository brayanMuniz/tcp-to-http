# tcp-to-http

An implementation of HTTP from TCP

[demo video](./demo.gif) 

## Overview
Raw TCP. No [net/http](https://pkg.go.dev/net/http) package. Uses a finite state machine approach in order to validate the protocol from incoming TCP packets. 

## Run and Test

Run it with: `go run ./cmd/httpserver/main.go`  
Test it with:   
`curl http://localhost:42069/ -v`  
`curl http://localhost:42069/give400 -v`  
`curl http://localhost:42069/give500 -v`  

## Features

* **Raw TCP Handling:** Listens for incoming connections directly on a TCP port.
* **HTTP/1.1 Parsing:** Capable of parsing essential HTTP/1.1 request components like method, path, and headers.
* **HTTP Response Generation:** Constructs valid HTTP/1.1 responses, including status lines, headers 
* **Request Routing:** Includes a basic routing mechanism to handle different request paths 
* **Error Handling:** Provides basic error handling for malformed requests or server issues.
