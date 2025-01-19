package gc_ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"gc.yashk.dev/gc-backend/internal/gc_webrtc"
	"gc.yashk.dev/gc-backend/internal/globals"
	"gc.yashk.dev/gc-backend/utils/types"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	unsafeConn, err := Upgrader.Upgrade(w, r, nil)

	fmt.Println("Tried connecting to the WS")

	if err != nil {
		fmt.Errorf("Error in creating the ws connection from the http connection: %w", err)
	}

	c := &types.ThreadSafeWriter{unsafeConn, sync.Mutex{}}
	defer c.Close()

	PeerConnection := gc_webrtc.ConnectWebRTC(c)
	defer PeerConnection.Close()

	message := types.WebsocketMessage{}

	for {
		_, raw, err := c.ReadMessage()
		if err != nil {
			fmt.Errorf("Failed to read the message: %v", err)
			return
		}

		// fmt.Println("Got this message, ", string(raw))

		if err := json.Unmarshal(raw, &message); err != nil {
			fmt.Errorf("Error in unmarshalling this json %v", err)
			return
		}

		fmt.Println(message.Event, "event")

		switch message.Event {
		case "candidate":
			globals.Log.Infof("The candidate event is triggered")
			gc_webrtc.AddPeerCandidate(PeerConnection, &message)
		case "answer":
			gc_webrtc.AnswerPeer(PeerConnection, &message)
		default:
			fmt.Println("This is the default case meaning nothing matched of the sorts here")
		}
	}
}
