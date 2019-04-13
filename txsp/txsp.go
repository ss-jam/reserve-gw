// Handle TX state parks reseervations
package txsp

import (
	"io"
	"net/http"
)

// Return the response to the calling event
func Reply(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to the Texas State Parks alternative booking resource. Let us know what you think.")
}

// Return the XHR response to the calling event
func Xhr(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to the Texas State Parks alternative booking resource. Let us know what you think.")
}
