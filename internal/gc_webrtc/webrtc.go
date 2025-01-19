package gc_webrtc

import (
	"encoding/json"
	"fmt"

	"gc.yashk.dev/gc-backend/internal/globals"
	"gc.yashk.dev/gc-backend/utils/types"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

func ConnectWebRTC(conn *types.ThreadSafeWriter) *webrtc.PeerConnection {
	PeerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		globals.Log.Errorf("Failed to create peer conection: %v", err)
	}

	// accept one video, and audio track
	// MAYBE ONE MORE SCREEN PRESENTATION
	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeAudio, webrtc.RTPCodecTypeVideo} {
		if _, err := PeerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			globals.Log.Errorf("Error occurred while accepting the transceiver tracks: %v", err)
		}
	}

	globals.ListLock.Lock()
	globals.PeerConnections = append(globals.PeerConnections, types.PeerConnectionState{PeerConnection, conn})
	globals.ListLock.Unlock()

	PeerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			globals.Log.Errorf("Failed to marshal the candidate to json: %v", err)
		}

		globals.Log.Infof("Sending the candidate to client %s", candidateString)

		if writeErr := conn.WriteJSON(&types.WebsocketMessage{
			Event: "candidate",
			Data:  string(candidateString),
		}); writeErr != nil {
			globals.Log.Errorf("Error while sending candidate through websocket, %v", err)
		}
	})

	PeerConnection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		globals.Log.Infof("Connection state changed %s", p)

		switch p {
		case webrtc.PeerConnectionStateFailed:
			if err := PeerConnection.Close(); err != nil {
				globals.Log.Errorf("failed to close the peer connection %v", err)
			}
		case webrtc.PeerConnectionStateClosed:
			SignalPeerConnections()
		default:
		}
	})

	PeerConnection.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		globals.Log.Infof("Got the track, Kind=%s, ID=%s, PayloadType=%d", tr.Kind(), tr.ID(), tr.PayloadType())
		fmt.Print("Got the track, Kind=%s, ID=%s, PayloadType=%d", tr.Kind(), tr.ID(), tr.PayloadType())
		trackLocal := AddTrack(tr)

		buf := make([]byte, 2000)
		rtpPkt := &rtp.Packet{}
		defer RemoveTrack(trackLocal)

		for {
			i, _, err := tr.Read(buf)
			if err != nil {
				return
			}

			if err = rtpPkt.Unmarshal(buf[:i]); err != nil {
				globals.Log.Errorf("Failed to unmarshal incoming RTP packets: %v", err)
				return
			}

			rtpPkt.Extension = false
			rtpPkt.Extensions = nil
			err = trackLocal.WriteRTP(rtpPkt)
			if err != nil {
				fmt.Errorf("Error is ", err.Error())
				return
			}
		}
	})

	PeerConnection.OnICEConnectionStateChange(func(is webrtc.ICEConnectionState) {
		globals.Log.Infof("ICE Connection state changed: %s", is)
	})

	SignalPeerConnections()

	return PeerConnection
}
