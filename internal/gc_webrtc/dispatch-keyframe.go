package gc_webrtc

import (
	"gc.yashk.dev/gc-backend/internal/globals"
	"github.com/pion/rtcp"
)

func DispatchKeyframe() {
	globals.ListLock.Lock()
	defer globals.ListLock.Unlock()

	for i := range globals.PeerConnections {
		for _, receiver := range globals.PeerConnections[i].PeerConnection.GetReceivers() {
			if receiver.Track() == nil {
				return
			}

			_ = globals.PeerConnections[i].PeerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}
