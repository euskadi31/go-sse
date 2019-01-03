// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sse

import (
	"fmt"
	"net/http"
)

// ResponseWriter struct
type ResponseWriter struct {
	http.ResponseWriter
	flusher http.Flusher
}

// Send data to client
func (rw *ResponseWriter) Send(data EventMarshaler) {
	b, err := data.MarshalEvent()
	if err != nil {
		return
	}

	fmt.Fprintf(rw, "%s\n", b)

	// Flush the data immediately instead of buffering it for later.
	rw.flusher.Flush()
}
