package gc_webrtc

import (
	"encoding/json"

	"gc.yashk.dev/gc-backend/internal/globals"
	"gc.yashk.dev/gc-backend/utils/types"
	"github.com/pion/webrtc/v4"
)

func AnswerPeer(peer *webrtc.PeerConnection, message *types.WebsocketMessage) {
	answer := webrtc.SessionDescription{}
	if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
		globals.Log.Errorf("Failed to unmarshal this SDP %v", err)
		return
	}

	globals.Log.Infof("Got the answer")

	if err := peer.SetRemoteDescription(answer); err != nil {
		globals.Log.Errorf("Errr in setting SDP of the remote peer %v", answer)
		return
	}
}
