package bingo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
	turboStreamMedia = "text/vnd.turbo-stream.html"
)

var (
	upgrader = websocket.Upgrader{}
)

// EchoRouter This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// Hotwire hotwire handlers which demonstate some of the capabilities
type Hotwire struct{}

// NewHotwire new hotwire handlers
func NewHotwire() *Hotwire {
	return &Hotwire{}
}

// writeMessage this constructs an SSE compatible message with a sequence, and
// line breaks from the output of a template
//
// This looks something like this:
//
//   event: message
//   id: 6
//   data: <turbo-stream action="replace" target="load">
//   data:     <template>
//   data:         <span id="load">04:20:13: 1.9</span>
//   data:     </template>
//   data: </turbo-stream>
//
func writeMessage(w io.Writer, id int, event, message string) error {

	_, err := fmt.Fprintf(w, "event: %s\nid: %d\n", event, id)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewBufferString(message))
	for scanner.Scan() {
		_, err = fmt.Fprintf(w, "data: %s\n", scanner.Text())
		if err != nil {
			return err
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	_, err = fmt.Fprint(w, "\n")
	if err != nil {
		return err
	}

	return nil
}
