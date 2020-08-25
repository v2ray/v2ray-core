package protocol

// A StreamID in QUIC
type StreamID uint64

// StreamType encodes if this is a unidirectional or bidirectional stream
type StreamType uint8

const (
	// StreamTypeUni is a unidirectional stream
	StreamTypeUni StreamType = iota
	// StreamTypeBidi is a bidirectional stream
	StreamTypeBidi
)

// InitiatedBy says if the stream was initiated by the client or by the server
func (s StreamID) InitiatedBy() Perspective {
	if s%2 == 0 {
		return PerspectiveClient
	}
	return PerspectiveServer
}

//Type says if this is a unidirectional or bidirectional stream
func (s StreamID) Type() StreamType {
	if s%4 >= 2 {
		return StreamTypeUni
	}
	return StreamTypeBidi
}

// StreamNum returns how many streams in total are below this
// Example: for stream 9 it returns 3 (i.e. streams 1, 5 and 9)
func (s StreamID) StreamNum() uint64 {
	return uint64(s/4) + 1
}

// MaxStreamID is the highest stream ID that a peer is allowed to open,
// when it is allowed to open numStreams.
func MaxStreamID(stype StreamType, numStreams uint64, pers Perspective) StreamID {
	if numStreams == 0 {
		return 0
	}
	var first StreamID
	switch stype {
	case StreamTypeBidi:
		switch pers {
		case PerspectiveClient:
			first = 0
		case PerspectiveServer:
			first = 1
		}
	case StreamTypeUni:
		switch pers {
		case PerspectiveClient:
			first = 2
		case PerspectiveServer:
			first = 3
		}
	}
	return first + 4*StreamID(numStreams-1)
}

// FirstStream returns the first valid stream ID
func FirstStream(stype StreamType, pers Perspective) StreamID {
	return MaxStreamID(stype, 1, pers)
}
