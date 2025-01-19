package gc_webrtc

import (
	"gc.yashk.dev/gc-backend/internal/globals"
	"github.com/pion/webrtc/v4"
)

func AddTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	globals.ListLock.Lock()
	defer func() {
		globals.ListLock.Unlock()
		SignalPeerConnections()
	}()

	// Create a new TrackLocal with the same codec as our incoming
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		panic(err)
	}

	globals.TrackLocals[t.ID()] = trackLocal
	return trackLocal
}
