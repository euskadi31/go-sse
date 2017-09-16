// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkRetry(b *testing.B) {
	b.ReportAllocs()

	msg := Retry{
		Duration: time.Second * 5,
	}

	for n := 0; n < b.N; n++ {
		msg.MarshalEvent()
	}
}

func TestEmptyRetry(t *testing.T) {
	msg := Retry{}

	_, err := msg.MarshalEvent()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrDurationEmpty.Error())
}

func TestRetry(t *testing.T) {
	expected := []byte("retry: 5000\r\n")
	msg := Retry{
		Duration: time.Second * 5,
	}

	actual, err := msg.MarshalEvent()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
