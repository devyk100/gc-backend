package gc_webrtc

import (
	"gc.yashk.dev/gc-backend/internal/globals"
	"github.com/pion/webrtc/v4"
)

func RemoveTrack(t *webrtc.TrackLocalStaticRTP) {
	globals.ListLock.Lock()
	defer func() {
		globals.ListLock.Unlock()
		SignalPeerConnections()
	}()
	// delete(globals.TrackLocals, t.ID())
}
