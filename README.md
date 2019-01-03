Server Sent Events for Go [![Last release](https://img.shields.io/github/release/euskadi31/go-sse.svg)](https://github.com/euskadi31/go-sse/releases/latest) [![Documentation](https://godoc.org/github.com/euskadi31/go-sse?status.svg)](https://godoc.org/github.com/euskadi31/go-sse)
=========================

[![Go Report Card](https://goreportcard.com/badge/github.com/euskadi31/go-sse)](https://goreportcard.com/report/github.com/euskadi31/go-sse)

| Branch  | Status | Coverage |
|---------|--------|----------|
| master  | [![Build Status](https://img.shields.io/travis/euskadi31/go-sse/master.svg)](https://travis-ci.org/euskadi31/go-sse) | [![Coveralls](https://img.shields.io/coveralls/euskadi31/go-sse/master.svg)](https://coveralls.io/github/euskadi31/go-sse?branch=master) |


Golang Server Sent Events server

## Example

### Server
```go
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
			case <-rw.CloseNotify:
				log.Println("Done")

				return
			}
		}
	})

	serve.SetRetry(time.Second * 5)

	http.Handle("/events", serve)

	log.Panic(http.ListenAndServe(":1337", nil))
}

```

### Client
```js
var client = new EventSource("http://localhost:1337/events");

client.onmessage = (msg) => {
    console.log(msg);
};
```

## License

go-sse is licensed under [the MIT license](LICENSE.md).
