// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sse

import (
	"bytes"
	"errors"
	"strconv"
	"time"
)

// ErrDurationEmpty message
var ErrDurationEmpty = errors.New("duration is empty")

// Retry struct
type Retry struct {
	Duration time.Duration
}

// MarshalEvent implements the EventMarshaler interface.
func (r Retry) MarshalEvent() ([]byte, error) {
	if r.Duration == 0 {
		return nil, ErrDurationEmpty
	}

	var buffer bytes.Buffer

	buffer.WriteString("retry: ")
	buffer.WriteString(strconv.Itoa(int(r.Duration / time.Millisecond)))
	buffer.WriteString("\r\n")

	return buffer.Bytes(), nil
}
