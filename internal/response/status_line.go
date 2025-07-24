package response

import "fmt"

type StatusCode int

const (
	OK                    StatusCode = 200
	BAD_REQUEST           StatusCode = 400
	INTERNAL_SERVER_ERROR StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {
	message := ""
	switch statusCode {
	case OK:
		message = "200 OK"
	case BAD_REQUEST:
		message = "400 BAD REQUEST"
	case INTERNAL_SERVER_ERROR:
		message = "500 INTERNAL SERVER ERROR"
	}

	return []byte(fmt.Sprintf("HTTP/1.1 %s", message))

}

