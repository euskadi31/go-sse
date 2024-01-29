// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sse

import (
	"net/http"
	"time"
)

// LastEventID header.
const LastEventID = "Last-Event-ID"

// HandlerFunc type.
type HandlerFunc func(ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *http.Request) {
	f(w, r)
}

// Server struct.
type Server struct {
	handle        HandlerFunc
	retryInterval time.Duration
}

// NewServer constructor.
func NewServer(handle HandlerFunc) *Server {
	return &Server{
		handle: handle,
	}
}

// SetRetry duration.
func (s *Server) SetRetry(duration time.Duration) {
	s.retryInterval = duration
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// Make sure that the writer supports flushing.
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	response := ResponseWriter{
		ResponseWriter: rw,
		flusher:        flusher,
	}

	if s.retryInterval > 0 {
		response.Send(&Retry{
			Duration: s.retryInterval,
		})
	}

	s.handle.ServeHTTP(response, r)
}
