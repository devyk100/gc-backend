package gc_webrtc

import (
	"encoding/json"

	"gc.yashk.dev/gc-backend/internal/globals"
	"gc.yashk.dev/gc-backend/utils/types"
	"github.com/pion/webrtc/v4"
)

func AddPeerCandidate(peer *webrtc.PeerConnection, message *types.WebsocketMessage) {
	candidate := webrtc.ICECandidateInit{}
	if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
		globals.Log.Errorf("failed to add this ICE candidate %v", err)
		return
	}

	globals.Log.Infof("Got a candidate")

	if err := peer.AddICECandidate(candidate); err != nil {
		globals.Log.Errorf("Failed to add ICE candidate: %v", err)
		return
	}

}
