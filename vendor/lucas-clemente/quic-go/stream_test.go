package quic

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/lucas-clemente/quic-go/internal/mocks"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

// in the tests for the stream deadlines we set a deadline
// and wait to make an assertion when Read / Write was unblocked
// on the CIs, the timing is a lot less precise, so scale every duration by this factor
func scaleDuration(t time.Duration) time.Duration {
	scaleFactor := 1
	if f, err := strconv.Atoi(os.Getenv("TIMESCALE_FACTOR")); err == nil { // parsing "" errors, so this works fine if the env is not set
		scaleFactor = f
	}
	Expect(scaleFactor).ToNot(BeZero())
	return time.Duration(scaleFactor) * t
}

var _ = Describe("Stream", func() {
	const streamID protocol.StreamID = 1337

	var (
		str            *stream
		strWithTimeout io.ReadWriter // str wrapped with gbytes.Timeout{Reader,Writer}
		mockFC         *mocks.MockStreamFlowController
		mockSender     *MockStreamSender
	)

	BeforeEach(func() {
		mockSender = NewMockStreamSender(mockCtrl)
		mockFC = mocks.NewMockStreamFlowController(mockCtrl)
		str = newStream(streamID, mockSender, mockFC, protocol.VersionWhatever)

		timeout := scaleDuration(250 * time.Millisecond)
		strWithTimeout = struct {
			io.Reader
			io.Writer
		}{
			gbytes.TimeoutReader(str, timeout),
			gbytes.TimeoutWriter(str, timeout),
		}
	})

	It("gets stream id", func() {
		Expect(str.StreamID()).To(Equal(protocol.StreamID(1337)))
	})

	// need some stream cancelation tests here, since gQUIC doesn't cleanly separate the two stream halves
	Context("stream cancelations", func() {
		Context("for gQUIC", func() {
			BeforeEach(func() {
				str.version = versionGQUICFrames
				str.receiveStream.version = versionGQUICFrames
				str.sendStream.version = versionGQUICFrames
			})

			It("unblocks Write when receiving a RST_STREAM frame with non-zero error code", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				mockSender.EXPECT().queueControlFrame(&wire.RstStreamFrame{
					StreamID:   streamID,
					ByteOffset: 1000,
					ErrorCode:  errorCodeStoppingGQUIC,
				})
				mockSender.EXPECT().onStreamCompleted(streamID)
				mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(6), true)
				str.writeOffset = 1000
				f := &wire.RstStreamFrame{
					StreamID:   streamID,
					ByteOffset: 6,
					ErrorCode:  123,
				}
				writeReturned := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					_, err := strWithTimeout.Write([]byte("foobar"))
					Expect(err).To(MatchError("Stream 1337 was reset with error code 123"))
					Expect(err).To(BeAssignableToTypeOf(streamCanceledError{}))
					Expect(err.(streamCanceledError).Canceled()).To(BeTrue())
					Expect(err.(streamCanceledError).ErrorCode()).To(Equal(protocol.ApplicationErrorCode(123)))
					close(writeReturned)
				}()
				Consistently(writeReturned).ShouldNot(BeClosed())
				err := str.handleRstStreamFrame(f)
				Expect(err).ToNot(HaveOccurred())
				Eventually(writeReturned).Should(BeClosed())
			})

			It("unblocks Write when receiving a RST_STREAM frame with error code 0", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				mockSender.EXPECT().queueControlFrame(&wire.RstStreamFrame{
					StreamID:   streamID,
					ByteOffset: 1000,
					ErrorCode:  errorCodeStoppingGQUIC,
				})
				mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(6), true)
				str.writeOffset = 1000
				f := &wire.RstStreamFrame{
					StreamID:   streamID,
					ByteOffset: 6,
					ErrorCode:  0,
				}
				writeReturned := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					_, err := strWithTimeout.Write([]byte("foobar"))
					Expect(err).To(MatchError("Stream 1337 was reset with error code 0"))
					Expect(err).To(BeAssignableToTypeOf(streamCanceledError{}))
					Expect(err.(streamCanceledError).Canceled()).To(BeTrue())
					Expect(err.(streamCanceledError).ErrorCode()).To(Equal(protocol.ApplicationErrorCode(0)))
					close(writeReturned)
				}()
				Consistently(writeReturned).ShouldNot(BeClosed())
				err := str.handleRstStreamFrame(f)
				Expect(err).ToNot(HaveOccurred())
				Eventually(writeReturned).Should(BeClosed())
			})

			It("sends a RST_STREAM with error code 0, after the stream is closed", func() {
				str.version = versionGQUICFrames
				mockSender.EXPECT().onHasStreamData(streamID).Times(2) // once for the Write, once for the Close
				mockFC.EXPECT().SendWindowSize().Return(protocol.MaxByteCount).AnyTimes()
				mockFC.EXPECT().AddBytesSent(protocol.ByteCount(6))
				err := str.CancelRead(1234)
				Expect(err).ToNot(HaveOccurred())
				writeReturned := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					_, err := strWithTimeout.Write([]byte("foobar"))
					Expect(err).ToNot(HaveOccurred())
					close(writeReturned)
				}()
				Eventually(func() *wire.StreamFrame {
					frame, _ := str.popStreamFrame(1000)
					return frame
				}).ShouldNot(BeNil())
				Eventually(writeReturned).Should(BeClosed())
				mockSender.EXPECT().queueControlFrame(&wire.RstStreamFrame{
					StreamID:   streamID,
					ByteOffset: 6,
					ErrorCode:  0,
				})
				Expect(str.Close()).To(Succeed())
			})
		})

		Context("for IETF QUIC", func() {
			It("doesn't queue a RST_STREAM after closing the stream", func() { // this is what it does for gQUIC
				mockSender.EXPECT().queueControlFrame(&wire.StopSendingFrame{
					StreamID:  streamID,
					ErrorCode: 1234,
				})
				mockSender.EXPECT().onHasStreamData(streamID)
				err := str.CancelRead(1234)
				Expect(err).ToNot(HaveOccurred())
				Expect(str.Close()).To(Succeed())
			})
		})
	})

	Context("deadlines", func() {
		It("sets a write deadline, when SetDeadline is called", func() {
			str.SetDeadline(time.Now().Add(-time.Second))
			n, err := strWithTimeout.Write([]byte("foobar"))
			Expect(err).To(MatchError(errDeadline))
			Expect(n).To(BeZero())
		})

		It("sets a read deadline, when SetDeadline is called", func() {
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(6), false).AnyTimes()
			f := &wire.StreamFrame{Data: []byte("foobar")}
			err := str.handleStreamFrame(f)
			Expect(err).ToNot(HaveOccurred())
			str.SetDeadline(time.Now().Add(-time.Second))
			b := make([]byte, 6)
			n, err := strWithTimeout.Read(b)
			Expect(err).To(MatchError(errDeadline))
			Expect(n).To(BeZero())
		})
	})

	Context("completing", func() {
		It("is not completed when only the receive side is completed", func() {
			// don't EXPECT a call to mockSender.onStreamCompleted()
			str.receiveStream.sender.onStreamCompleted(streamID)
		})

		It("is not completed when only the send side is completed", func() {
			// don't EXPECT a call to mockSender.onStreamCompleted()
			str.sendStream.sender.onStreamCompleted(streamID)
		})

		It("is completed when both sides are completed", func() {
			mockSender.EXPECT().onStreamCompleted(streamID)
			str.sendStream.sender.onStreamCompleted(streamID)
			str.receiveStream.sender.onStreamCompleted(streamID)
		})
	})
})
