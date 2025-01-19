package types

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
type ThreadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func (t *ThreadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(v)
}
