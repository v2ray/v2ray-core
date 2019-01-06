package handshake

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

type transportParameterID uint16

const (
	originalConnectionIDParameterID           transportParameterID = 0x0
	idleTimeoutParameterID                    transportParameterID = 0x1
	statelessResetTokenParameterID            transportParameterID = 0x2
	maxPacketSizeParameterID                  transportParameterID = 0x3
	initialMaxDataParameterID                 transportParameterID = 0x4
	initialMaxStreamDataBidiLocalParameterID  transportParameterID = 0x5
	initialMaxStreamDataBidiRemoteParameterID transportParameterID = 0x6
	initialMaxStreamDataUniParameterID        transportParameterID = 0x7
	initialMaxStreamsBidiParameterID          transportParameterID = 0x8
	initialMaxStreamsUniParameterID           transportParameterID = 0x9
	disableMigrationParameterID               transportParameterID = 0xc
)

// TransportParameters are parameters sent to the peer during the handshake
type TransportParameters struct {
	InitialMaxStreamDataBidiLocal  protocol.ByteCount
	InitialMaxStreamDataBidiRemote protocol.ByteCount
	InitialMaxStreamDataUni        protocol.ByteCount
	InitialMaxData                 protocol.ByteCount

	MaxPacketSize protocol.ByteCount

	MaxUniStreams  uint64
	MaxBidiStreams uint64

	IdleTimeout      time.Duration
	DisableMigration bool

	StatelessResetToken  []byte
	OriginalConnectionID protocol.ConnectionID
}

func (p *TransportParameters) unmarshal(data []byte, sentBy protocol.Perspective) error {
	// needed to check that every parameter is only sent at most once
	var parameterIDs []transportParameterID

	r := bytes.NewReader(data)
	for r.Len() >= 4 {
		paramIDInt, _ := utils.BigEndian.ReadUint16(r)
		paramID := transportParameterID(paramIDInt)
		paramLen, _ := utils.BigEndian.ReadUint16(r)
		parameterIDs = append(parameterIDs, paramID)
		switch paramID {
		case initialMaxStreamDataBidiLocalParameterID,
			initialMaxStreamDataBidiRemoteParameterID,
			initialMaxStreamDataUniParameterID,
			initialMaxDataParameterID,
			initialMaxStreamsBidiParameterID,
			initialMaxStreamsUniParameterID,
			idleTimeoutParameterID,
			maxPacketSizeParameterID:
			if err := p.readNumericTransportParameter(r, paramID, int(paramLen)); err != nil {
				return err
			}
		default:
			if r.Len() < int(paramLen) {
				return fmt.Errorf("remaining length (%d) smaller than parameter length (%d)", r.Len(), paramLen)
			}
			switch paramID {
			case disableMigrationParameterID:
				if paramLen != 0 {
					return fmt.Errorf("wrong length for disable_migration: %d (expected empty)", paramLen)
				}
				p.DisableMigration = true
			case statelessResetTokenParameterID:
				if sentBy == protocol.PerspectiveClient {
					return errors.New("client sent a stateless_reset_token")
				}
				if paramLen != 16 {
					return fmt.Errorf("wrong length for stateless_reset_token: %d (expected 16)", paramLen)
				}
				b := make([]byte, 16)
				r.Read(b)
				p.StatelessResetToken = b
			case originalConnectionIDParameterID:
				if sentBy == protocol.PerspectiveClient {
					return errors.New("client sent an original_connection_id")
				}
				p.OriginalConnectionID, _ = protocol.ReadConnectionID(r, int(paramLen))
			default:
				r.Seek(int64(paramLen), io.SeekCurrent)
			}
		}
	}

	// check that every transport parameter was sent at most once
	sort.Slice(parameterIDs, func(i, j int) bool { return parameterIDs[i] < parameterIDs[j] })
	for i := 0; i < len(parameterIDs)-1; i++ {
		if parameterIDs[i] == parameterIDs[i+1] {
			return fmt.Errorf("received duplicate transport parameter %#x", parameterIDs[i])
		}
	}

	if r.Len() != 0 {
		return fmt.Errorf("should have read all data. Still have %d bytes", r.Len())
	}
	return nil
}

func (p *TransportParameters) readNumericTransportParameter(
	r *bytes.Reader,
	paramID transportParameterID,
	expectedLen int,
) error {
	remainingLen := r.Len()
	val, err := utils.ReadVarInt(r)
	if err != nil {
		return fmt.Errorf("error while reading transport parameter %d: %s", paramID, err)
	}
	if remainingLen-r.Len() != expectedLen {
		return fmt.Errorf("inconsistent transport parameter length for %d", paramID)
	}
	switch paramID {
	case initialMaxStreamDataBidiLocalParameterID:
		p.InitialMaxStreamDataBidiLocal = protocol.ByteCount(val)
	case initialMaxStreamDataBidiRemoteParameterID:
		p.InitialMaxStreamDataBidiRemote = protocol.ByteCount(val)
	case initialMaxStreamDataUniParameterID:
		p.InitialMaxStreamDataUni = protocol.ByteCount(val)
	case initialMaxDataParameterID:
		p.InitialMaxData = protocol.ByteCount(val)
	case initialMaxStreamsBidiParameterID:
		p.MaxBidiStreams = val
	case initialMaxStreamsUniParameterID:
		p.MaxUniStreams = val
	case idleTimeoutParameterID:
		p.IdleTimeout = utils.MaxDuration(protocol.MinRemoteIdleTimeout, time.Duration(val)*time.Second)
	case maxPacketSizeParameterID:
		if val < 1200 {
			return fmt.Errorf("invalid value for max_packet_size: %d (minimum 1200)", val)
		}
		p.MaxPacketSize = protocol.ByteCount(val)
	default:
		return fmt.Errorf("TransportParameter BUG: transport parameter %d not found", paramID)
	}
	return nil
}

func (p *TransportParameters) marshal(b *bytes.Buffer) {
	// initial_max_stream_data_bidi_local
	utils.BigEndian.WriteUint16(b, uint16(initialMaxStreamDataBidiLocalParameterID))
	utils.BigEndian.WriteUint16(b, uint16(utils.VarIntLen(uint64(p.InitialMaxStreamDataBidiLocal))))
	utils.WriteVarInt(b, uint64(p.InitialMaxStreamDataBidiLocal))
	// initial_max_stream_data_bidi_remote
	utils.BigEndian.WriteUint16(b, uint16(initialMaxStreamDataBidiRemoteParameterID))
	utils.BigEndian.WriteUint16(b, uint16(utils.VarIntLen(uint64(p.InitialMaxStreamDataBidiRemote))))
	utils.WriteVarInt(b, uint64(p.InitialMaxStreamDataBidiRemote))
	// initial_max_stream_data_uni
	utils.BigEndian.WriteUint16(b, uint16(initialMaxStreamDataUniParameterID))
	utils.BigEndian.WriteUint16(b, uint16(utils.VarIntLen(uint64(p.InitialMaxStreamDataUni))))
	utils.WriteVarInt(b, uint64(p.InitialMaxStreamDataUni))
	// initial_max_data
	utils.BigEndian.WriteUint16(b, uint16(initialMaxDataParameterID))
	utils.BigEndian.WriteUint16(b, uint16(utils.VarIntLen(uint64(p.InitialMaxData))))
	utils.WriteVarInt(b, uint64(p.InitialMaxData))
	// initial_max_bidi_streams
	utils.BigEndian.WriteUint16(b, uint16(initialMaxStreamsBidiParameterID))
	utils.BigEndian.WriteUint16(b, uint16(utils.VarIntLen(p.MaxBidiStreams)))
	utils.WriteVarInt(b, p.MaxBidiStreams)
	// initial_max_uni_streams
	utils.BigEndian.WriteUint16(b, uint16(initialMaxStreamsUniParameterID))
	utils.BigEndian.WriteUint16(b, uint16(utils.VarIntLen(p.MaxUniStreams)))
	utils.WriteVarInt(b, p.MaxUniStreams)
	// idle_timeout
	utils.BigEndian.WriteUint16(b, uint16(idleTimeoutParameterID))
	utils.BigEndian.WriteUint16(b, uint16(utils.VarIntLen(uint64(p.IdleTimeout/time.Second))))
	utils.WriteVarInt(b, uint64(p.IdleTimeout/time.Second))
	// max_packet_size
	utils.BigEndian.WriteUint16(b, uint16(maxPacketSizeParameterID))
	utils.BigEndian.WriteUint16(b, uint16(utils.VarIntLen(uint64(protocol.MaxReceivePacketSize))))
	utils.WriteVarInt(b, uint64(protocol.MaxReceivePacketSize))
	// disable_migration
	if p.DisableMigration {
		utils.BigEndian.WriteUint16(b, uint16(disableMigrationParameterID))
		utils.BigEndian.WriteUint16(b, 0)
	}
	if len(p.StatelessResetToken) > 0 {
		utils.BigEndian.WriteUint16(b, uint16(statelessResetTokenParameterID))
		utils.BigEndian.WriteUint16(b, uint16(len(p.StatelessResetToken))) // should always be 16 bytes
		b.Write(p.StatelessResetToken)
	}
	// original_connection_id
	if p.OriginalConnectionID.Len() > 0 {
		utils.BigEndian.WriteUint16(b, uint16(originalConnectionIDParameterID))
		utils.BigEndian.WriteUint16(b, uint16(p.OriginalConnectionID.Len()))
		b.Write(p.OriginalConnectionID.Bytes())
	}
}

// String returns a string representation, intended for logging.
func (p *TransportParameters) String() string {
	return fmt.Sprintf("&handshake.TransportParameters{OriginalConnectionID: %s, InitialMaxStreamDataBidiLocal: %#x, InitialMaxStreamDataBidiRemote: %#x, InitialMaxStreamDataUni: %#x, InitialMaxData: %#x, MaxBidiStreams: %d, MaxUniStreams: %d, IdleTimeout: %s}", p.OriginalConnectionID, p.InitialMaxStreamDataBidiLocal, p.InitialMaxStreamDataBidiRemote, p.InitialMaxStreamDataUni, p.InitialMaxData, p.MaxBidiStreams, p.MaxUniStreams, p.IdleTimeout)
}
