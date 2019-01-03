// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sse

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3200*time.Millisecond)
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	req = req.WithContext(ctx)

	serve := NewServer(func(rw ResponseWriter, r *http.Request) {
		tickChan := time.NewTicker(1 * time.Second).C

		i := 1

		// recovery
		lastID := r.Header.Get(LastEventID)
		if lastID != "" {
			i, _ = strconv.Atoi(lastID)
		}

		for {
			select {
			case <-tickChan:
				eventString := fmt.Sprintf("tick %d", i)

				rw.Send(&MessageEvent{
					ID:   strconv.Itoa(i),
					Data: []byte(eventString),
				})

				i++
			case <-r.Context().Done():
				return
			}
		}
	})

	serve.SetRetry(time.Second * 5)

	serve.ServeHTTP(w, req)

	resp := w.Result()

	reader := bufio.NewReader(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))

	i := 0

	for {
		line, err := reader.ReadBytes('\r')
		if err == io.EOF {
			assert.Equal(t, 7, i)
			break
		}
		line = bytes.TrimSpace(line)

		switch i {
		case 0:
			assert.Equal(t, "retry: 5000", string(line))

		case 1:
			assert.Equal(t, "id: 1", string(line))

		case 2:
			assert.Equal(t, "data: tick 1", string(line))

		case 3:
			assert.Equal(t, "id: 2", string(line))

		case 4:
			assert.Equal(t, "data: tick 2", string(line))

		case 5:
			assert.Equal(t, "id: 3", string(line))

		case 6:
			assert.Equal(t, "data: tick 3", string(line))
		}

		i++
	}
}

func TestServerWithLastEventID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3200*time.Millisecond)
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)

	req.Header.Set(LastEventID, "2")

	w := httptest.NewRecorder()

	req = req.WithContext(ctx)

	serve := NewServer(func(rw ResponseWriter, r *http.Request) {
		tickChan := time.NewTicker(1 * time.Second).C

		i := 1

		// recovery
		lastID := r.Header.Get(LastEventID)
		if lastID != "" {
			i, _ = strconv.Atoi(lastID)
		}

		for {
			select {
			case <-tickChan:
				eventString := fmt.Sprintf("tick %d", i)

				rw.Send(&MessageEvent{
					ID:   strconv.Itoa(i),
					Data: []byte(eventString),
				})

				i++
			case <-r.Context().Done():
				return
			}
		}
	})

	serve.SetRetry(time.Second * 5)

	serve.ServeHTTP(w, req)

	resp := w.Result()

	reader := bufio.NewReader(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))

	i := 0

	for {
		line, err := reader.ReadBytes('\r')
		if err == io.EOF {
			assert.Equal(t, 7, i)
			break
		}
		line = bytes.TrimSpace(line)

		switch i {
		case 0:
			assert.Equal(t, "retry: 5000", string(line))

		case 1:
			assert.Equal(t, "id: 2", string(line))

		case 2:
			assert.Equal(t, "data: tick 2", string(line))

		case 3:
			assert.Equal(t, "id: 3", string(line))

		case 4:
			assert.Equal(t, "data: tick 3", string(line))

		case 5:
			assert.Equal(t, "id: 4", string(line))

		case 6:
			assert.Equal(t, "data: tick 4", string(line))
		}

		i++
	}
}

type responseWriter struct {
	StatusCode int
}

func (responseWriter) Header() http.Header {
	return http.Header{}
}

func (responseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
}

func TestServerWithoutFlusher(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := &responseWriter{}

	serve := NewServer(func(rw ResponseWriter, r *http.Request) {})

	serve.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.StatusCode)
}
