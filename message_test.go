// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkMessageEvent(b *testing.B) {
	b.ReportAllocs()

	msg := MessageEvent{
		ID:    "1",
		Event: "foo",
		Data:  []byte("bar\nsecond line"),
	}

	for n := 0; n < b.N; n++ {
		msg.MarshalEvent()
	}
}

func TestEmptyMessageEvent(t *testing.T) {
	msg := MessageEvent{}

	_, err := msg.MarshalEvent()
	assert.Error(t, err)
	assert.EqualError(t, err, ErrMessageEmpty.Error())
}

func TestMessageEvent(t *testing.T) {

	var messagesTests = []struct {
		msg      MessageEvent // data
		expected []byte       // expected result
	}{
		{
			MessageEvent{
				ID:    "1",
				Event: "foo",
				Data:  []byte("bar\nsecond line"),
			},
			[]byte("id: 1\r\nevent: foo\r\ndata: bar\r\ndata: second line\r\n"),
		},
		{
			MessageEvent{
				Event: "foo",
				Data:  []byte("bar\nsecond line"),
			},
			[]byte("event: foo\r\ndata: bar\r\ndata: second line\r\n"),
		},
		{
			MessageEvent{
				ID:   "1",
				Data: []byte("bar\nsecond line"),
			},
			[]byte("id: 1\r\ndata: bar\r\ndata: second line\r\n"),
		},
		{
			MessageEvent{
				Data: []byte("bar\nsecond line"),
			},
			[]byte("data: bar\r\ndata: second line\r\n"),
		},
		{
			MessageEvent{
				ID:    "1",
				Event: "foo",
				Data:  []byte("bar"),
			},
			[]byte("id: 1\r\nevent: foo\r\ndata: bar\r\n"),
		},
	}

	for _, tt := range messagesTests {
		actual, err := tt.msg.MarshalEvent()

		assert.NoError(t, err)
		assert.Equal(t, tt.expected, actual)
	}
}
