package gc_webrtc

import (
	"encoding/json"
	"fmt"
	"time"

	"gc.yashk.dev/gc-backend/internal/globals"
	"gc.yashk.dev/gc-backend/utils/types"
	"github.com/pion/webrtc/v4"
)

func attemptSync() bool {
	for i := range globals.PeerConnections {
		if globals.PeerConnections[i].PeerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
			// remove the ith peer connection from this list
			globals.PeerConnections = append(globals.PeerConnections[:i], globals.PeerConnections[i+1:]...)
			return true // this list was modified so we must restart to properly do a signal to all webrtc connections
		}

		existingSenders := map[string]bool{}

		for _, sender := range globals.PeerConnections[i].PeerConnection.GetSenders() {
			if sender.Track() == nil {
				continue
			}

			existingSenders[sender.Track().ID()] = true

			if _, ok := globals.TrackLocals[sender.Track().ID()]; !ok {
				if err := globals.PeerConnections[i].PeerConnection.RemoveTrack(sender); err != nil {
					return true
				}
			}
		}

		for _, receiver := range globals.PeerConnections[i].PeerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			existingSenders[receiver.Track().ID()] = true
		}

		for trackID := range globals.TrackLocals {
			if _, ok := existingSenders[trackID]; !ok {
				fmt.Println("Adding the track", trackID, globals.TrackLocals[trackID])
				if _, err := globals.PeerConnections[i].PeerConnection.AddTrack(globals.TrackLocals[trackID]); err != nil {
					fmt.Print("Error in adding tracks, ", err.Error())
					return true
				}
			}
		}

		offer, err := globals.PeerConnections[i].PeerConnection.CreateOffer(nil)
		if err != nil {
			return true
		}

		if err := globals.PeerConnections[i].PeerConnection.SetLocalDescription(offer); err != nil {
			return true
		}

		offerString, err := json.Marshal(offer)
		if err != nil {
			globals.Log.Errorf("Failed to marshal offer to json %v", err)
			return true
		}

		globals.Log.Infof("Send offer to client: %v", offer)

		if err = globals.PeerConnections[i].Websocket.WriteJSON(&types.WebsocketMessage{
			Event: "offer",
			Data:  string(offerString),
		}); err != nil {
			return true
		}
	}

	return false
}

func SignalPeerConnections() {
	globals.ListLock.Lock()
	defer func() {
		globals.ListLock.Unlock()
		DispatchKeyframe()
	}()
	for syncAttemmpt := 0; ; syncAttemmpt++ {
		if syncAttemmpt == 25 {
			go func() {
				time.Sleep(time.Second * 3)
				SignalPeerConnections()
			}()
			return
		}

		fmt.Println("This is attempt", syncAttemmpt, " from signal peer connections")
		if !attemptSync() {
			break
		}
	}
}
