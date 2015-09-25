package webBLL

import (
	"fmt"
	"net/http"
)

var (
	// 传输外部消息的Channel
	MessageCh = make(chan string, 100)
)

// 接收需要广播的消息
func ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if message, ok := r.Form["Message"]; ok {
		MessageCh <- message[0]
	}

	fmt.Fprintf(w, "OK")
}
