Server Sent Events for Go
=========================

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
