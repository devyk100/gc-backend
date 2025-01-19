package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"gc.yashk.dev/gc-backend/internal/gc_ws"
	"gc.yashk.dev/gc-backend/internal/globals"
	"gc.yashk.dev/gc-backend/utils/types"
	"github.com/pion/webrtc/v4"
)

var addr = ":8085"
var ListLock sync.RWMutex

func main() {
	log.Println("Hello world from the orchestrator")
	flag.Parse()
	globals.TrackLocals = map[string]*webrtc.TrackLocalStaticRTP{}
	globals.PeerConnections = []types.PeerConnectionState{}
	ListLock = sync.RWMutex{}
	http.HandleFunc("/ws", gc_ws.WebsocketHandler)

	if err := http.ListenAndServe(addr, nil); err != nil { //nolint: gosec
		fmt.Errorf("Failed to start http server: %v", err)
	}
}
