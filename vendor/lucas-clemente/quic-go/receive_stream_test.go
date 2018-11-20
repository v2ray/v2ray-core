package quic

import (
	"errors"
	"io"
	"runtime"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lucas-clemente/quic-go/internal/mocks"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Receive Stream", func() {
	const streamID protocol.StreamID = 1337

	var (
		str            *receiveStream
		strWithTimeout io.Reader // str wrapped with gbytes.TimeoutReader
		mockFC         *mocks.MockStreamFlowController
		mockSender     *MockStreamSender
	)

	BeforeEach(func() {
		mockSender = NewMockStreamSender(mockCtrl)
		mockFC = mocks.NewMockStreamFlowController(mockCtrl)
		str = newReceiveStream(streamID, mockSender, mockFC, versionIETFFrames)

		timeout := scaleDuration(250 * time.Millisecond)
		strWithTimeout = gbytes.TimeoutReader(str, timeout)
	})

	It("gets stream id", func() {
		Expect(str.StreamID()).To(Equal(protocol.StreamID(1337)))
	})

	Context("reading", func() {
		It("reads a single STREAM frame", func() {
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(4), false)
			mockFC.EXPECT().AddBytesRead(protocol.ByteCount(4))
			mockFC.EXPECT().MaybeQueueWindowUpdate()
			frame := wire.StreamFrame{
				Offset: 0,
				Data:   []byte{0xDE, 0xAD, 0xBE, 0xEF},
			}
			err := str.handleStreamFrame(&frame)
			Expect(err).ToNot(HaveOccurred())
			b := make([]byte, 4)
			n, err := strWithTimeout.Read(b)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(4))
			Expect(b).To(Equal([]byte{0xDE, 0xAD, 0xBE, 0xEF}))
		})

		It("reads a single STREAM frame in multiple goes", func() {
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(4), false)
			mockFC.EXPECT().AddBytesRead(protocol.ByteCount(2))
			mockFC.EXPECT().AddBytesRead(protocol.ByteCount(2))
			mockFC.EXPECT().MaybeQueueWindowUpdate().Times(2)
			frame := wire.StreamFrame{
				Offset: 0,
				Data:   []byte{0xDE, 0xAD, 0xBE, 0xEF},
			}
			err := str.handleStreamFrame(&frame)
			Expect(err).ToNot(HaveOccurred())
			b := make([]byte, 2)
			n, err := strWithTimeout.Read(b)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(2))
			Expect(b).To(Equal([]byte{0xDE, 0xAD}))
			n, err = strWithTimeout.Read(b)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(2))
			Expect(b).To(Equal([]byte{0xBE, 0xEF}))
		})

		It("reads all data available", func() {
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(2), false)
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(4), false)
			mockFC.EXPECT().AddBytesRead(protocol.ByteCount(2)).Times(2)
			mockFC.EXPECT().MaybeQueueWindowUpdate().Times(2)
			frame1 := wire.StreamFrame{
				Offset: 0,
				Data:   []byte{0xDE, 0xAD},
			}
			frame2 := wire.StreamFrame{
				Offset: 2,
				Data:   []byte{0xBE, 0xEF},
			}
			err := str.handleStreamFrame(&frame1)
			Expect(err).ToNot(HaveOccurred())
			err = str.handleStreamFrame(&frame2)
			Expect(err).ToNot(HaveOccurred())
			b := make([]byte, 6)
			n, err := strWithTimeout.Read(b)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(4))
			Expect(b).To(Equal([]byte{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x00}))
		})

		It("assembles multiple STREAM frames", func() {
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(2), false)
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(4), false)
			mockFC.EXPECT().AddBytesRead(protocol.ByteCount(2)).Times(2)
			mockFC.EXPECT().MaybeQueueWindowUpdate().Times(2)
			frame1 := wire.StreamFrame{
				Offset: 0,
				Data:   []byte{0xDE, 0xAD},
			}
			frame2 := wire.StreamFrame{
				Offset: 2,
				Data:   []byte{0xBE, 0xEF},
			}
			err := str.handleStreamFrame(&frame1)
			Expect(err).ToNot(HaveOccurred())
			err = str.handleStreamFrame(&frame2)
			Expect(err).ToNot(HaveOccurred())
			b := make([]byte, 4)
			n, err := strWithTimeout.Read(b)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(4))
			Expect(b).To(Equal([]byte{0xDE, 0xAD, 0xBE, 0xEF}))
		})

		It("waits until data is available", func() {
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(2), false)
			mockFC.EXPECT().AddBytesRead(protocol.ByteCount(2))
			mockFC.EXPECT().MaybeQueueWindowUpdate()
			go func() {
				defer GinkgoRecover()
				frame := wire.StreamFrame{Data: []byte{0xDE, 0xAD}}
				time.Sleep(10 * time.Millisecond)
				err := str.handleStreamFrame(&frame)
				Expect(err).ToNot(HaveOccurred())
			}()
			b := make([]byte, 2)
			n, err := strWithTimeout.Read(b)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(2))
		})

		It("handles STREAM frames in wrong order", func() {
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(2), false)
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(4), false)
			mockFC.EXPECT().AddBytesRead(protocol.ByteCount(2)).Times(2)
			mockFC.EXPECT().MaybeQueueWindowUpdate().Times(2)
			frame1 := wire.StreamFrame{
				Offset: 2,
				Data:   []byte{0xBE, 0xEF},
			}
			frame2 := wire.StreamFrame{
				Offset: 0,
				Data:   []byte{0xDE, 0xAD},
			}
			err := str.handleStreamFrame(&frame1)
			Expect(err).ToNot(HaveOccurred())
			err = str.handleStreamFrame(&frame2)
			Expect(err).ToNot(HaveOccurred())
			b := make([]byte, 4)
			n, err := strWithTimeout.Read(b)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(4))
			Expect(b).To(Equal([]byte{0xDE, 0xAD, 0xBE, 0xEF}))
		})

		It("ignores duplicate STREAM frames", func() {
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(2), false)
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(2), false)
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(4), false)
			mockFC.EXPECT().AddBytesRead(protocol.ByteCount(2)).Times(2)
			mockFC.EXPECT().MaybeQueueWindowUpdate().Times(2)
			frame1 := wire.StreamFrame{
				Offset: 0,
				Data:   []byte{0xDE, 0xAD},
			}
			frame2 := wire.StreamFrame{
				Offset: 0,
				Data:   []byte{0x13, 0x37},
			}
			frame3 := wire.StreamFrame{
				Offset: 2,
				Data:   []byte{0xBE, 0xEF},
			}
			err := str.handleStreamFrame(&frame1)
			Expect(err).ToNot(HaveOccurred())
			err = str.handleStreamFrame(&frame2)
			Expect(err).ToNot(HaveOccurred())
			err = str.handleStreamFrame(&frame3)
			Expect(err).ToNot(HaveOccurred())
			b := make([]byte, 4)
			n, err := strWithTimeout.Read(b)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(4))
			Expect(b).To(Equal([]byte{0xDE, 0xAD, 0xBE, 0xEF}))
		})

		It("doesn't rejects a STREAM frames with an overlapping data range", func() {
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(4), false)
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(6), false)
			mockFC.EXPECT().AddBytesRead(protocol.ByteCount(2))
			mockFC.EXPECT().AddBytesRead(protocol.ByteCount(4))
			mockFC.EXPECT().MaybeQueueWindowUpdate().Times(2)
			frame1 := wire.StreamFrame{
				Offset: 0,
				Data:   []byte("foob"),
			}
			frame2 := wire.StreamFrame{
				Offset: 2,
				Data:   []byte("obar"),
			}
			err := str.handleStreamFrame(&frame1)
			Expect(err).ToNot(HaveOccurred())
			err = str.handleStreamFrame(&frame2)
			Expect(err).ToNot(HaveOccurred())
			b := make([]byte, 6)
			n, err := strWithTimeout.Read(b)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(6))
			Expect(b).To(Equal([]byte("foobar")))
		})

		Context("deadlines", func() {
			It("the deadline error has the right net.Error properties", func() {
				Expect(errDeadline.Temporary()).To(BeTrue())
				Expect(errDeadline.Timeout()).To(BeTrue())
				Expect(errDeadline).To(MatchError("deadline exceeded"))
			})

			It("returns an error when Read is called after the deadline", func() {
				mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(6), false).AnyTimes()
				f := &wire.StreamFrame{Data: []byte("foobar")}
				err := str.handleStreamFrame(f)
				Expect(err).ToNot(HaveOccurred())
				str.SetReadDeadline(time.Now().Add(-time.Second))
				b := make([]byte, 6)
				n, err := strWithTimeout.Read(b)
				Expect(err).To(MatchError(errDeadline))
				Expect(n).To(BeZero())
			})

			It("unblocks after the deadline", func() {
				deadline := time.Now().Add(scaleDuration(50 * time.Millisecond))
				str.SetReadDeadline(deadline)
				b := make([]byte, 6)
				n, err := strWithTimeout.Read(b)
				Expect(err).To(MatchError(errDeadline))
				Expect(n).To(BeZero())
				Expect(time.Now()).To(BeTemporally("~", deadline, scaleDuration(10*time.Millisecond)))
			})

			It("doesn't unblock if the deadline is changed before the first one expires", func() {
				deadline1 := time.Now().Add(scaleDuration(50 * time.Millisecond))
				deadline2 := time.Now().Add(scaleDuration(100 * time.Millisecond))
				str.SetReadDeadline(deadline1)
				go func() {
					defer GinkgoRecover()
					time.Sleep(scaleDuration(20 * time.Millisecond))
					str.SetReadDeadline(deadline2)
					// make sure that this was actually execute before the deadline expires
					Expect(time.Now()).To(BeTemporally("<", deadline1))
				}()
				runtime.Gosched()
				b := make([]byte, 10)
				n, err := strWithTimeout.Read(b)
				Expect(err).To(MatchError(errDeadline))
				Expect(n).To(BeZero())
				Expect(time.Now()).To(BeTemporally("~", deadline2, scaleDuration(20*time.Millisecond)))
			})

			It("unblocks earlier, when a new deadline is set", func() {
				deadline1 := time.Now().Add(scaleDuration(200 * time.Millisecond))
				deadline2 := time.Now().Add(scaleDuration(50 * time.Millisecond))
				go func() {
					defer GinkgoRecover()
					time.Sleep(scaleDuration(10 * time.Millisecond))
					str.SetReadDeadline(deadline2)
					// make sure that this was actually execute before the deadline expires
					Expect(time.Now()).To(BeTemporally("<", deadline2))
				}()
				str.SetReadDeadline(deadline1)
				runtime.Gosched()
				b := make([]byte, 10)
				_, err := strWithTimeout.Read(b)
				Expect(err).To(MatchError(errDeadline))
				Expect(time.Now()).To(BeTemporally("~", deadline2, scaleDuration(25*time.Millisecond)))
			})
		})

		Context("closing", func() {
			Context("with FIN bit", func() {
				It("returns EOFs", func() {
					mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(4), true)
					mockFC.EXPECT().AddBytesRead(protocol.ByteCount(4))
					mockFC.EXPECT().MaybeQueueWindowUpdate()
					str.handleStreamFrame(&wire.StreamFrame{
						Offset: 0,
						Data:   []byte{0xDE, 0xAD, 0xBE, 0xEF},
						FinBit: true,
					})
					mockSender.EXPECT().onStreamCompleted(streamID)
					b := make([]byte, 4)
					n, err := strWithTimeout.Read(b)
					Expect(err).To(MatchError(io.EOF))
					Expect(n).To(Equal(4))
					Expect(b).To(Equal([]byte{0xDE, 0xAD, 0xBE, 0xEF}))
					n, err = strWithTimeout.Read(b)
					Expect(n).To(BeZero())
					Expect(err).To(MatchError(io.EOF))
				})

				It("handles out-of-order frames", func() {
					mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(2), false)
					mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(4), true)
					mockFC.EXPECT().AddBytesRead(protocol.ByteCount(2)).Times(2)
					mockFC.EXPECT().MaybeQueueWindowUpdate().Times(2)
					frame1 := wire.StreamFrame{
						Offset: 2,
						Data:   []byte{0xBE, 0xEF},
						FinBit: true,
					}
					frame2 := wire.StreamFrame{
						Offset: 0,
						Data:   []byte{0xDE, 0xAD},
					}
					err := str.handleStreamFrame(&frame1)
					Expect(err).ToNot(HaveOccurred())
					err = str.handleStreamFrame(&frame2)
					Expect(err).ToNot(HaveOccurred())
					mockSender.EXPECT().onStreamCompleted(streamID)
					b := make([]byte, 4)
					n, err := strWithTimeout.Read(b)
					Expect(err).To(MatchError(io.EOF))
					Expect(n).To(Equal(4))
					Expect(b).To(Equal([]byte{0xDE, 0xAD, 0xBE, 0xEF}))
					n, err = strWithTimeout.Read(b)
					Expect(n).To(BeZero())
					Expect(err).To(MatchError(io.EOF))
				})

				It("returns EOFs with partial read", func() {
					mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(2), true)
					mockFC.EXPECT().AddBytesRead(protocol.ByteCount(2))
					mockFC.EXPECT().MaybeQueueWindowUpdate()
					err := str.handleStreamFrame(&wire.StreamFrame{
						Offset: 0,
						Data:   []byte{0xde, 0xad},
						FinBit: true,
					})
					Expect(err).ToNot(HaveOccurred())
					mockSender.EXPECT().onStreamCompleted(streamID)
					b := make([]byte, 4)
					n, err := strWithTimeout.Read(b)
					Expect(err).To(MatchError(io.EOF))
					Expect(n).To(Equal(2))
					Expect(b[:n]).To(Equal([]byte{0xde, 0xad}))
				})

				It("handles immediate FINs", func() {
					mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(0), true)
					mockFC.EXPECT().AddBytesRead(protocol.ByteCount(0))
					mockFC.EXPECT().MaybeQueueWindowUpdate()
					err := str.handleStreamFrame(&wire.StreamFrame{
						Offset: 0,
						FinBit: true,
					})
					Expect(err).ToNot(HaveOccurred())
					mockSender.EXPECT().onStreamCompleted(streamID)
					b := make([]byte, 4)
					n, err := strWithTimeout.Read(b)
					Expect(n).To(BeZero())
					Expect(err).To(MatchError(io.EOF))
				})
			})

			It("closes when CloseRemote is called", func() {
				mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(0), true)
				mockFC.EXPECT().AddBytesRead(protocol.ByteCount(0))
				mockFC.EXPECT().MaybeQueueWindowUpdate()
				str.CloseRemote(0)
				mockSender.EXPECT().onStreamCompleted(streamID)
				b := make([]byte, 8)
				n, err := strWithTimeout.Read(b)
				Expect(n).To(BeZero())
				Expect(err).To(MatchError(io.EOF))
			})
		})

		Context("closing for shutdown", func() {
			testErr := errors.New("test error")

			It("immediately returns all reads", func() {
				done := make(chan struct{})
				b := make([]byte, 4)
				go func() {
					defer GinkgoRecover()
					n, err := strWithTimeout.Read(b)
					Expect(n).To(BeZero())
					Expect(err).To(MatchError(testErr))
					close(done)
				}()
				Consistently(done).ShouldNot(BeClosed())
				str.closeForShutdown(testErr)
				Eventually(done).Should(BeClosed())
			})

			It("errors for all following reads", func() {
				str.closeForShutdown(testErr)
				b := make([]byte, 1)
				n, err := strWithTimeout.Read(b)
				Expect(n).To(BeZero())
				Expect(err).To(MatchError(testErr))
			})
		})
	})

	Context("stream cancelations", func() {
		Context("canceling read", func() {
			It("unblocks Read", func() {
				mockSender.EXPECT().queueControlFrame(gomock.Any())
				done := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					_, err := strWithTimeout.Read([]byte{0})
					Expect(err).To(MatchError("Read on stream 1337 canceled with error code 1234"))
					close(done)
				}()
				Consistently(done).ShouldNot(BeClosed())
				err := str.CancelRead(1234)
				Expect(err).ToNot(HaveOccurred())
				Eventually(done).Should(BeClosed())
			})

			It("doesn't allow further calls to Read", func() {
				mockSender.EXPECT().queueControlFrame(gomock.Any())
				err := str.CancelRead(1234)
				Expect(err).ToNot(HaveOccurred())
				_, err = strWithTimeout.Read([]byte{0})
				Expect(err).To(MatchError("Read on stream 1337 canceled with error code 1234"))
			})

			It("does nothing when CancelRead is called twice", func() {
				mockSender.EXPECT().queueControlFrame(gomock.Any())
				err := str.CancelRead(1234)
				Expect(err).ToNot(HaveOccurred())
				err = str.CancelRead(2345)
				Expect(err).ToNot(HaveOccurred())
				_, err = strWithTimeout.Read([]byte{0})
				Expect(err).To(MatchError("Read on stream 1337 canceled with error code 1234"))
			})

			It("doesn't send a RST_STREAM frame, if the FIN was already read", func() {
				mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(6), true)
				mockFC.EXPECT().AddBytesRead(protocol.ByteCount(6))
				mockFC.EXPECT().MaybeQueueWindowUpdate()
				// no calls to mockSender.queueControlFrame
				err := str.handleStreamFrame(&wire.StreamFrame{
					StreamID: streamID,
					Data:     []byte("foobar"),
					FinBit:   true,
				})
				Expect(err).ToNot(HaveOccurred())
				mockSender.EXPECT().onStreamCompleted(streamID)
				_, err = strWithTimeout.Read(make([]byte, 100))
				Expect(err).To(MatchError(io.EOF))
				err = str.CancelRead(1234)
				Expect(err).ToNot(HaveOccurred())
			})

			It("queues a STOP_SENDING frame, for IETF QUIC", func() {
				str.version = versionIETFFrames
				mockSender.EXPECT().queueControlFrame(&wire.StopSendingFrame{
					StreamID:  streamID,
					ErrorCode: 1234,
				})
				err := str.CancelRead(1234)
				Expect(err).ToNot(HaveOccurred())
			})

			It("doesn't queue a STOP_SENDING frame, for gQUIC", func() {
				str.version = versionGQUICFrames
				// no calls to mockSender.queueControlFrame
				err := str.CancelRead(1234)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("receiving RST_STREAM frames", func() {
			rst := &wire.RstStreamFrame{
				StreamID:   streamID,
				ByteOffset: 42,
				ErrorCode:  1234,
			}

			It("unblocks Read", func() {
				mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(42), true)
				done := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					_, err := strWithTimeout.Read([]byte{0})
					Expect(err).To(MatchError("Stream 1337 was reset with error code 1234"))
					Expect(err).To(BeAssignableToTypeOf(streamCanceledError{}))
					Expect(err.(streamCanceledError).Canceled()).To(BeTrue())
					Expect(err.(streamCanceledError).ErrorCode()).To(Equal(protocol.ApplicationErrorCode(1234)))
					close(done)
				}()
				Consistently(done).ShouldNot(BeClosed())
				mockSender.EXPECT().onStreamCompleted(streamID)
				str.handleRstStreamFrame(rst)
				Eventually(done).Should(BeClosed())
			})

			It("doesn't allow further calls to Read", func() {
				mockSender.EXPECT().onStreamCompleted(streamID)
				mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(42), true)
				err := str.handleRstStreamFrame(rst)
				Expect(err).ToNot(HaveOccurred())
				_, err = strWithTimeout.Read([]byte{0})
				Expect(err).To(MatchError("Stream 1337 was reset with error code 1234"))
				Expect(err).To(BeAssignableToTypeOf(streamCanceledError{}))
				Expect(err.(streamCanceledError).Canceled()).To(BeTrue())
				Expect(err.(streamCanceledError).ErrorCode()).To(Equal(protocol.ApplicationErrorCode(1234)))
			})

			It("errors when receiving a RST_STREAM with an inconsistent offset", func() {
				testErr := errors.New("already received a different final offset before")
				mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(42), true).Return(testErr)
				err := str.handleRstStreamFrame(rst)
				Expect(err).To(MatchError(testErr))
			})

			It("ignores duplicate RST_STREAM frames", func() {
				mockSender.EXPECT().onStreamCompleted(streamID)
				mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(42), true).Times(2)
				err := str.handleRstStreamFrame(rst)
				Expect(err).ToNot(HaveOccurred())
				err = str.handleRstStreamFrame(rst)
				Expect(err).ToNot(HaveOccurred())
			})

			It("doesn't do anyting when it was closed for shutdown", func() {
				str.closeForShutdown(nil)
				err := str.handleRstStreamFrame(rst)
				Expect(err).ToNot(HaveOccurred())
			})

			Context("for gQUIC", func() {
				BeforeEach(func() {
					str.version = versionGQUICFrames
				})

				It("unblocks Read when receiving a RST_STREAM frame with non-zero error code", func() {
					mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(42), true)
					readReturned := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						_, err := strWithTimeout.Read([]byte{0})
						Expect(err).To(MatchError("Stream 1337 was reset with error code 1234"))
						Expect(err).To(BeAssignableToTypeOf(streamCanceledError{}))
						Expect(err.(streamCanceledError).Canceled()).To(BeTrue())
						Expect(err.(streamCanceledError).ErrorCode()).To(Equal(protocol.ApplicationErrorCode(1234)))
						close(readReturned)
					}()
					Consistently(readReturned).ShouldNot(BeClosed())
					mockSender.EXPECT().onStreamCompleted(streamID)
					err := str.handleRstStreamFrame(rst)
					Expect(err).ToNot(HaveOccurred())
					Eventually(readReturned).Should(BeClosed())
				})

				It("continues reading until the end when receiving a RST_STREAM frame with error code 0", func() {
					mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(6), true).Times(2)
					gomock.InOrder(
						mockFC.EXPECT().AddBytesRead(protocol.ByteCount(4)),
						mockFC.EXPECT().AddBytesRead(protocol.ByteCount(2)),
						mockSender.EXPECT().onStreamCompleted(streamID),
					)
					mockFC.EXPECT().MaybeQueueWindowUpdate().Times(2)
					readReturned := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						n, err := strWithTimeout.Read(make([]byte, 4))
						Expect(err).ToNot(HaveOccurred())
						Expect(n).To(Equal(4))
						n, err = strWithTimeout.Read(make([]byte, 4))
						Expect(err).To(MatchError(io.EOF))
						Expect(n).To(Equal(2))
						close(readReturned)
					}()
					Consistently(readReturned).ShouldNot(BeClosed())
					err := str.handleStreamFrame(&wire.StreamFrame{
						StreamID: streamID,
						Data:     []byte("foobar"),
						FinBit:   true,
					})
					Expect(err).ToNot(HaveOccurred())
					err = str.handleRstStreamFrame(&wire.RstStreamFrame{
						StreamID:   streamID,
						ByteOffset: 6,
						ErrorCode:  0,
					})
					Expect(err).ToNot(HaveOccurred())
					Eventually(readReturned).Should(BeClosed())
				})
			})
		})
	})

	Context("flow control", func() {
		It("errors when a STREAM frame causes a flow control violation", func() {
			testErr := errors.New("flow control violation")
			mockFC.EXPECT().UpdateHighestReceived(protocol.ByteCount(8), false).Return(testErr)
			frame := wire.StreamFrame{
				Offset: 2,
				Data:   []byte("foobar"),
			}
			err := str.handleStreamFrame(&frame)
			Expect(err).To(MatchError(testErr))
		})

		It("gets a window update", func() {
			mockFC.EXPECT().GetWindowUpdate().Return(protocol.ByteCount(0x100))
			Expect(str.getWindowUpdate()).To(Equal(protocol.ByteCount(0x100)))
		})
	})
})
