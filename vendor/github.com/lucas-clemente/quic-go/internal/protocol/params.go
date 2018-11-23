package protocol

import "time"

// MaxPacketSizeIPv4 is the maximum packet size that we use for sending IPv4 packets.
const MaxPacketSizeIPv4 = 1252

// MaxPacketSizeIPv6 is the maximum packet size that we use for sending IPv6 packets.
const MaxPacketSizeIPv6 = 1232

// MinStatelessResetSize is the minimum size of a stateless reset packet
const MinStatelessResetSize = 1 + 20 + 16

// NonForwardSecurePacketSizeReduction is the number of bytes a non forward-secure packet has to be smaller than a forward-secure packet
// This makes sure that those packets can always be retransmitted without splitting the contained StreamFrames
const NonForwardSecurePacketSizeReduction = 50

const defaultMaxCongestionWindowPackets = 1000

// DefaultMaxCongestionWindow is the default for the max congestion window
const DefaultMaxCongestionWindow ByteCount = defaultMaxCongestionWindowPackets * DefaultTCPMSS

// InitialCongestionWindow is the initial congestion window in QUIC packets
const InitialCongestionWindow ByteCount = 32 * DefaultTCPMSS

// MaxUndecryptablePackets limits the number of undecryptable packets that a
// session queues for later until it sends a public reset.
const MaxUndecryptablePackets = 10

// ConnectionFlowControlMultiplier determines how much larger the connection flow control windows needs to be relative to any stream's flow control window
// This is the value that Chromium is using
const ConnectionFlowControlMultiplier = 1.5

// InitialMaxStreamData is the stream-level flow control window for receiving data
const InitialMaxStreamData = (1 << 10) * 512 // 512 kb

// InitialMaxData is the connection-level flow control window for receiving data
const InitialMaxData = ConnectionFlowControlMultiplier * InitialMaxStreamData

// DefaultMaxReceiveStreamFlowControlWindow is the default maximum stream-level flow control window for receiving data, for the server
const DefaultMaxReceiveStreamFlowControlWindow = 6 * (1 << 20) // 6 MB

// DefaultMaxReceiveConnectionFlowControlWindow is the default connection-level flow control window for receiving data, for the server
const DefaultMaxReceiveConnectionFlowControlWindow = 15 * (1 << 20) // 12 MB

// WindowUpdateThreshold is the fraction of the receive window that has to be consumed before an higher offset is advertised to the client
const WindowUpdateThreshold = 0.25

// DefaultMaxIncomingStreams is the maximum number of streams that a peer may open
const DefaultMaxIncomingStreams = 100

// DefaultMaxIncomingUniStreams is the maximum number of unidirectional streams that a peer may open
const DefaultMaxIncomingUniStreams = 100

// MaxSessionUnprocessedPackets is the max number of packets stored in each session that are not yet processed.
const MaxSessionUnprocessedPackets = defaultMaxCongestionWindowPackets

// SkipPacketAveragePeriodLength is the average period length in which one packet number is skipped to prevent an Optimistic ACK attack
const SkipPacketAveragePeriodLength PacketNumber = 500

// MaxTrackedSkippedPackets is the maximum number of skipped packet numbers the SentPacketHandler keep track of for Optimistic ACK attack mitigation
const MaxTrackedSkippedPackets = 10

// CookieExpiryTime is the valid time of a cookie
const CookieExpiryTime = 24 * time.Hour

// MaxOutstandingSentPackets is maximum number of packets saved for retransmission.
// When reached, it imposes a soft limit on sending new packets:
// Sending ACKs and retransmission is still allowed, but now new regular packets can be sent.
const MaxOutstandingSentPackets = 2 * defaultMaxCongestionWindowPackets

// MaxTrackedSentPackets is maximum number of sent packets saved for retransmission.
// When reached, no more packets will be sent.
// This value *must* be larger than MaxOutstandingSentPackets.
const MaxTrackedSentPackets = MaxOutstandingSentPackets * 5 / 4

// MaxTrackedReceivedAckRanges is the maximum number of ACK ranges tracked
const MaxTrackedReceivedAckRanges = defaultMaxCongestionWindowPackets

// MaxNonRetransmittableAcks is the maximum number of packets containing an ACK, but no retransmittable frames, that we send in a row
const MaxNonRetransmittableAcks = 19

// MaxStreamFrameSorterGaps is the maximum number of gaps between received StreamFrames
// prevents DoS attacks against the streamFrameSorter
const MaxStreamFrameSorterGaps = 1000

// MaxCryptoStreamOffset is the maximum offset allowed on any of the crypto streams.
// This limits the size of the ClientHello and Certificates that can be received.
const MaxCryptoStreamOffset = 16 * (1 << 10)

// MinRemoteIdleTimeout is the minimum value that we accept for the remote idle timeout
const MinRemoteIdleTimeout = 5 * time.Second

// DefaultIdleTimeout is the default idle timeout
const DefaultIdleTimeout = 30 * time.Second

// DefaultHandshakeTimeout is the default timeout for a connection until the crypto handshake succeeds.
const DefaultHandshakeTimeout = 10 * time.Second

// RetiredConnectionIDDeleteTimeout is the time we keep closed sessions around in order to retransmit the CONNECTION_CLOSE.
// after this time all information about the old connection will be deleted
const RetiredConnectionIDDeleteTimeout = 5 * time.Second

// MinStreamFrameSize is the minimum size that has to be left in a packet, so that we add another STREAM frame.
// This avoids splitting up STREAM frames into small pieces, which has 2 advantages:
// 1. it reduces the framing overhead
// 2. it reduces the head-of-line blocking, when a packet is lost
const MinStreamFrameSize ByteCount = 128

// MaxAckFrameSize is the maximum size for an ACK frame that we write
// Due to the varint encoding, ACK frames can grow (almost) indefinitely large.
// The MaxAckFrameSize should be large enough to encode many ACK range,
// but must ensure that a maximum size ACK frame fits into one packet.
const MaxAckFrameSize ByteCount = 1000

// MinPacingDelay is the minimum duration that is used for packet pacing
// If the packet packing frequency is higher, multiple packets might be sent at once.
// Example: For a packet pacing delay of 20 microseconds, we would send 5 packets at once, wait for 100 microseconds, and so forth.
const MinPacingDelay time.Duration = 100 * time.Microsecond

// DefaultConnectionIDLength is the connection ID length that is used for multiplexed connections
// if no other value is configured.
const DefaultConnectionIDLength = 4
