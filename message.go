// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sse

import (
	"bytes"
	"errors"
)

// ErrMessageEmpty message.
var ErrMessageEmpty = errors.New("message is empty")

// MessageEvent struct.
type MessageEvent struct {
	ID    string
	Event string
	Data  []byte
}

// MarshalEvent implements the encoding.TextMarshaler interface.
func (m MessageEvent) MarshalEvent() ([]byte, error) {
	if len(m.Data) == 0 {
		return nil, ErrMessageEmpty
	}

	var buffer bytes.Buffer

	if m.ID != "" {
		buffer.WriteString("id: ")
		buffer.WriteString(m.ID)
		buffer.WriteString("\r\n")
	}

	if m.Event != "" {
		buffer.WriteString("event: ")
		buffer.WriteString(m.Event)
		buffer.WriteString("\r\n")
	}

	data := bytes.ReplaceAll(m.Data, []byte("\r"), []byte(""))
	lines := bytes.Split(data, []byte("\n"))

	for _, line := range lines {
		buffer.WriteString("data: ")
		buffer.Write(line)
		buffer.WriteString("\r\n")
	}

	return buffer.Bytes(), nil
}
