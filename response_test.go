// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sse

import (
	"errors"
	"testing"
)

type event struct{}

func (event) MarshalEvent() ([]byte, error) {
	return nil, errors.New("fail")
}

func TestResponseWriter(t *testing.T) {
	wr := &ResponseWriter{}

	wr.Send(&event{})
}
