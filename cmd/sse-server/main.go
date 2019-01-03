// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/euskadi31/go-sse"
)

func main() {
	serve := sse.NewServer(func(rw sse.ResponseWriter, r *http.Request) {
		tickChan := time.NewTicker(time.Second * 2).C

		// recovery
		lastID := r.Header.Get("Last-Event-ID")
		if lastID != "" {
			log.Printf("Recovery with ID: %s\n", lastID)
		}

		for {
			select {
			case t := <-tickChan:
				eventString := fmt.Sprintf("the time is %v", t)

				log.Println("Send event...")

				rw.Send(&sse.MessageEvent{
					ID:   strconv.Itoa(int(t.Unix())),
					Data: []byte(eventString),
				})
			case <-r.Context().Done():
				log.Println("Done")

				return
			}
		}
	})

	serve.SetRetry(time.Second * 5)

	http.Handle("/events", serve)

	log.Panic(http.ListenAndServe(":1337", nil))
}
