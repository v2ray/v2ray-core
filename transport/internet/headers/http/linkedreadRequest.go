package http

import (
	"bufio"
	"net/http"

	_ "unsafe" // required to use //go:linkname
)

//go:linkname readRequest net/http.readRequest
func readRequest(b *bufio.Reader, deleteHostHeader bool) (req *http.Request, err error)
