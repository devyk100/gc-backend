package globals

import (
	"sync"

	"gc.yashk.dev/gc-backend/utils/types"
	"github.com/pion/logging"
	"github.com/pion/webrtc/v4"
)

var ListLock sync.RWMutex
var PeerConnections []types.PeerConnectionState
var TrackLocals map[string]*webrtc.TrackLocalStaticRTP
var Log = logging.NewDefaultLoggerFactory().NewLogger("sfu-ws")
