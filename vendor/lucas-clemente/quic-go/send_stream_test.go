package quic

import (
	"bytes"
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

var _ = Describe("Send Stream", func() {
	const streamID protocol.StreamID = 1337

	var (
		str            *sendStream
		strWithTimeout io.Writer // str wrapped with gbytes.TimeoutWriter
		mockFC         *mocks.MockStreamFlowController
		mockSender     *MockStreamSender
	)

	BeforeEach(func() {
		mockSender = NewMockStreamSender(mockCtrl)
		mockFC = mocks.NewMockStreamFlowController(mockCtrl)
		str = newSendStream(streamID, mockSender, mockFC, protocol.VersionWhatever)

		timeout := scaleDuration(250 * time.Millisecond)
		strWithTimeout = gbytes.TimeoutWriter(str, timeout)
	})

	waitForWrite := func() {
		EventuallyWithOffset(0, func() []byte {
			str.mutex.Lock()
			data := str.dataForWriting
			str.mutex.Unlock()
			return data
		}).ShouldNot(BeEmpty())
	}

	It("gets stream id", func() {
		Expect(str.StreamID()).To(Equal(protocol.StreamID(1337)))
	})

	Context("writing", func() {
		It("writes and gets all data at once", func() {
			mockSender.EXPECT().onHasStreamData(streamID)
			mockFC.EXPECT().SendWindowSize().Return(protocol.ByteCount(9999))
			mockFC.EXPECT().AddBytesSent(protocol.ByteCount(6))
			done := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				n, err := strWithTimeout.Write([]byte("foobar"))
				Expect(err).ToNot(HaveOccurred())
				Expect(n).To(Equal(6))
				close(done)
			}()
			waitForWrite()
			f, _ := str.popStreamFrame(1000)
			Expect(f.Data).To(Equal([]byte("foobar")))
			Expect(f.FinBit).To(BeFalse())
			Expect(f.Offset).To(BeZero())
			Expect(f.DataLenPresent).To(BeTrue())
			Expect(str.writeOffset).To(Equal(protocol.ByteCount(6)))
			Expect(str.dataForWriting).To(BeNil())
			Eventually(done).Should(BeClosed())
		})

		It("writes and gets data in two turns", func() {
			mockSender.EXPECT().onHasStreamData(streamID)
			frameHeaderLen := protocol.ByteCount(4)
			mockFC.EXPECT().SendWindowSize().Return(protocol.ByteCount(9999)).Times(2)
			mockFC.EXPECT().AddBytesSent(gomock.Any() /* protocol.ByteCount(3)*/).Times(2)
			done := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				n, err := strWithTimeout.Write([]byte("foobar"))
				Expect(err).ToNot(HaveOccurred())
				Expect(n).To(Equal(6))
				close(done)
			}()
			waitForWrite()
			f, _ := str.popStreamFrame(3 + frameHeaderLen)
			Expect(f.Data).To(Equal([]byte("foo")))
			Expect(f.FinBit).To(BeFalse())
			Expect(f.Offset).To(BeZero())
			Expect(f.DataLenPresent).To(BeTrue())
			f, _ = str.popStreamFrame(100)
			Expect(f.Data).To(Equal([]byte("bar")))
			Expect(f.FinBit).To(BeFalse())
			Expect(f.Offset).To(Equal(protocol.ByteCount(3)))
			Expect(f.DataLenPresent).To(BeTrue())
			Expect(str.popStreamFrame(1000)).To(BeNil())
			Eventually(done).Should(BeClosed())
		})

		It("popStreamFrame returns nil if no data is available", func() {
			frame, hasMoreData := str.popStreamFrame(1000)
			Expect(frame).To(BeNil())
			Expect(hasMoreData).To(BeFalse())
		})

		It("says if it has more data for writing", func() {
			mockSender.EXPECT().onHasStreamData(streamID)
			mockFC.EXPECT().SendWindowSize().Return(protocol.ByteCount(9999)).Times(2)
			mockFC.EXPECT().AddBytesSent(gomock.Any()).Times(2)
			done := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				n, err := strWithTimeout.Write(bytes.Repeat([]byte{0}, 100))
				Expect(err).ToNot(HaveOccurred())
				Expect(n).To(Equal(100))
				close(done)
			}()
			waitForWrite()
			frame, hasMoreData := str.popStreamFrame(50)
			Expect(frame).ToNot(BeNil())
			Expect(hasMoreData).To(BeTrue())
			frame, hasMoreData = str.popStreamFrame(1000)
			Expect(frame).ToNot(BeNil())
			Expect(hasMoreData).To(BeFalse())
			frame, _ = str.popStreamFrame(1000)
			Expect(frame).To(BeNil())
			Eventually(done).Should(BeClosed())
		})

		It("copies the slice while writing", func() {
			mockSender.EXPECT().onHasStreamData(streamID)
			frameHeaderSize := protocol.ByteCount(4)
			mockFC.EXPECT().SendWindowSize().Return(protocol.ByteCount(9999)).Times(2)
			mockFC.EXPECT().AddBytesSent(protocol.ByteCount(1))
			mockFC.EXPECT().AddBytesSent(protocol.ByteCount(2))
			s := []byte("foo")
			done := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				n, err := strWithTimeout.Write(s)
				Expect(err).ToNot(HaveOccurred())
				Expect(n).To(Equal(3))
				close(done)
			}()
			waitForWrite()
			frame, _ := str.popStreamFrame(frameHeaderSize + 1)
			Expect(frame.Data).To(Equal([]byte("f")))
			s[1] = 'e'
			f, _ := str.popStreamFrame(100)
			Expect(f).ToNot(BeNil())
			Expect(f.Data).To(Equal([]byte("oo")))
			Eventually(done).Should(BeClosed())
		})

		It("returns when given a nil input", func() {
			n, err := strWithTimeout.Write(nil)
			Expect(n).To(BeZero())
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns when given an empty slice", func() {
			n, err := strWithTimeout.Write([]byte(""))
			Expect(n).To(BeZero())
			Expect(err).ToNot(HaveOccurred())
		})

		It("cancels the context when Close is called", func() {
			mockSender.EXPECT().onHasStreamData(streamID)
			Expect(str.Context().Done()).ToNot(BeClosed())
			str.Close()
			Expect(str.Context().Done()).To(BeClosed())
		})

		Context("flow control blocking", func() {
			It("queues a BLOCKED frame if the stream is flow control blocked", func() {
				mockFC.EXPECT().SendWindowSize().Return(protocol.ByteCount(0))
				mockFC.EXPECT().IsNewlyBlocked().Return(true, protocol.ByteCount(12))
				mockSender.EXPECT().queueControlFrame(&wire.StreamBlockedFrame{
					StreamID: streamID,
					Offset:   12,
				})
				mockSender.EXPECT().onHasStreamData(streamID)
				done := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					_, err := str.Write([]byte("foobar"))
					Expect(err).ToNot(HaveOccurred())
					close(done)
				}()
				waitForWrite()
				f, hasMoreData := str.popStreamFrame(1000)
				Expect(f).To(BeNil())
				Expect(hasMoreData).To(BeFalse())
				// make the Write go routine return
				str.closeForShutdown(nil)
				Eventually(done).Should(BeClosed())
			})

			It("says that it doesn't have any more data, when it is flow control blocked", func() {
				frameHeaderSize := protocol.ByteCount(4)
				mockSender.EXPECT().onHasStreamData(streamID)

				done := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					_, err := str.Write([]byte("foobar"))
					Expect(err).ToNot(HaveOccurred())
					close(done)
				}()
				waitForWrite()

				// first pop a STREAM frame of the maximum size allowed by flow control
				mockFC.EXPECT().SendWindowSize().Return(protocol.ByteCount(3))
				mockFC.EXPECT().AddBytesSent(protocol.ByteCount(3))
				f, hasMoreData := str.popStreamFrame(frameHeaderSize + 3)
				Expect(f).ToNot(BeNil())
				Expect(hasMoreData).To(BeTrue())

				// try to pop again, this time noticing that we're blocked
				mockFC.EXPECT().SendWindowSize()
				// don't use offset 3 here, to make sure the BLOCKED frame contains the number returned by the flow controller
				mockFC.EXPECT().IsNewlyBlocked().Return(true, protocol.ByteCount(10))
				mockSender.EXPECT().queueControlFrame(&wire.StreamBlockedFrame{
					StreamID: streamID,
					Offset:   10,
				})
				f, hasMoreData = str.popStreamFrame(1000)
				Expect(f).To(BeNil())
				Expect(hasMoreData).To(BeFalse())
				// make the Write go routine return
				str.closeForShutdown(nil)
				Eventually(done).Should(BeClosed())
			})
		})

		Context("deadlines", func() {
			It("returns an error when Write is called after the deadline", func() {
				str.SetWriteDeadline(time.Now().Add(-time.Second))
				n, err := strWithTimeout.Write([]byte("foobar"))
				Expect(err).To(MatchError(errDeadline))
				Expect(n).To(BeZero())
			})

			It("unblocks after the deadline", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				deadline := time.Now().Add(scaleDuration(50 * time.Millisecond))
				str.SetWriteDeadline(deadline)
				n, err := strWithTimeout.Write([]byte("foobar"))
				Expect(err).To(MatchError(errDeadline))
				Expect(n).To(BeZero())
				Expect(time.Now()).To(BeTemporally("~", deadline, scaleDuration(20*time.Millisecond)))
			})

			It("returns the number of bytes written, when the deadline expires", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				mockFC.EXPECT().SendWindowSize().Return(protocol.ByteCount(10000)).AnyTimes()
				mockFC.EXPECT().AddBytesSent(gomock.Any())
				deadline := time.Now().Add(scaleDuration(50 * time.Millisecond))
				str.SetWriteDeadline(deadline)
				var n int
				writeReturned := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					var err error
					n, err = strWithTimeout.Write(bytes.Repeat([]byte{0}, 100))
					Expect(err).To(MatchError(errDeadline))
					Expect(time.Now()).To(BeTemporally("~", deadline, scaleDuration(20*time.Millisecond)))
					close(writeReturned)
				}()
				waitForWrite()
				frame, hasMoreData := str.popStreamFrame(50)
				Expect(frame).ToNot(BeNil())
				Expect(hasMoreData).To(BeTrue())
				Eventually(writeReturned, scaleDuration(80*time.Millisecond)).Should(BeClosed())
				Expect(n).To(BeEquivalentTo(frame.DataLen()))
			})

			It("doesn't pop any data after the deadline expired", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				mockFC.EXPECT().SendWindowSize().Return(protocol.ByteCount(10000)).AnyTimes()
				mockFC.EXPECT().AddBytesSent(gomock.Any())
				deadline := time.Now().Add(scaleDuration(50 * time.Millisecond))
				str.SetWriteDeadline(deadline)
				writeReturned := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					_, err := strWithTimeout.Write(bytes.Repeat([]byte{0}, 100))
					Expect(err).To(MatchError(errDeadline))
					close(writeReturned)
				}()
				waitForWrite()
				frame, hasMoreData := str.popStreamFrame(50)
				Expect(frame).ToNot(BeNil())
				Expect(hasMoreData).To(BeTrue())
				Eventually(writeReturned, scaleDuration(80*time.Millisecond)).Should(BeClosed())
				frame, hasMoreData = str.popStreamFrame(50)
				Expect(frame).To(BeNil())
				Expect(hasMoreData).To(BeFalse())
			})

			It("doesn't unblock if the deadline is changed before the first one expires", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				deadline1 := time.Now().Add(scaleDuration(50 * time.Millisecond))
				deadline2 := time.Now().Add(scaleDuration(100 * time.Millisecond))
				str.SetWriteDeadline(deadline1)
				done := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					time.Sleep(scaleDuration(20 * time.Millisecond))
					str.SetWriteDeadline(deadline2)
					// make sure that this was actually execute before the deadline expires
					Expect(time.Now()).To(BeTemporally("<", deadline1))
					close(done)
				}()
				runtime.Gosched()
				n, err := strWithTimeout.Write([]byte("foobar"))
				Expect(err).To(MatchError(errDeadline))
				Expect(n).To(BeZero())
				Expect(time.Now()).To(BeTemporally("~", deadline2, scaleDuration(20*time.Millisecond)))
				Eventually(done).Should(BeClosed())
			})

			It("unblocks earlier, when a new deadline is set", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				deadline1 := time.Now().Add(scaleDuration(200 * time.Millisecond))
				deadline2 := time.Now().Add(scaleDuration(50 * time.Millisecond))
				done := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					time.Sleep(scaleDuration(10 * time.Millisecond))
					str.SetWriteDeadline(deadline2)
					// make sure that this was actually execute before the deadline expires
					Expect(time.Now()).To(BeTemporally("<", deadline2))
					close(done)
				}()
				str.SetWriteDeadline(deadline1)
				runtime.Gosched()
				_, err := strWithTimeout.Write([]byte("foobar"))
				Expect(err).To(MatchError(errDeadline))
				Expect(time.Now()).To(BeTemporally("~", deadline2, scaleDuration(20*time.Millisecond)))
				Eventually(done).Should(BeClosed())
			})
		})

		Context("closing", func() {
			It("doesn't allow writes after it has been closed", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				str.Close()
				_, err := strWithTimeout.Write([]byte("foobar"))
				Expect(err).To(MatchError("write on closed stream 1337"))
			})

			It("allows FIN", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				mockSender.EXPECT().onStreamCompleted(streamID)
				str.Close()
				f, hasMoreData := str.popStreamFrame(1000)
				Expect(f).ToNot(BeNil())
				Expect(f.Data).To(BeEmpty())
				Expect(f.FinBit).To(BeTrue())
				Expect(hasMoreData).To(BeFalse())
			})

			It("doesn't send a FIN when there's still data", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				frameHeaderLen := protocol.ByteCount(4)
				mockFC.EXPECT().SendWindowSize().Return(protocol.ByteCount(9999)).Times(2)
				mockFC.EXPECT().AddBytesSent(gomock.Any()).Times(2)
				str.dataForWriting = []byte("foobar")
				Expect(str.Close()).To(Succeed())
				f, _ := str.popStreamFrame(3 + frameHeaderLen)
				Expect(f).ToNot(BeNil())
				Expect(f.Data).To(Equal([]byte("foo")))
				Expect(f.FinBit).To(BeFalse())
				mockSender.EXPECT().onStreamCompleted(streamID)
				f, _ = str.popStreamFrame(100)
				Expect(f.Data).To(Equal([]byte("bar")))
				Expect(f.FinBit).To(BeTrue())
			})

			It("doesn't allow FIN after it is closed for shutdown", func() {
				str.closeForShutdown(errors.New("test"))
				f, hasMoreData := str.popStreamFrame(1000)
				Expect(f).To(BeNil())
				Expect(hasMoreData).To(BeFalse())
			})

			It("doesn't allow FIN twice", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				mockSender.EXPECT().onStreamCompleted(streamID)
				str.Close()
				f, _ := str.popStreamFrame(1000)
				Expect(f).ToNot(BeNil())
				Expect(f.Data).To(BeEmpty())
				Expect(f.FinBit).To(BeTrue())
				f, hasMoreData := str.popStreamFrame(1000)
				Expect(f).To(BeNil())
				Expect(hasMoreData).To(BeFalse())
			})
		})

		Context("closing for shutdown", func() {
			testErr := errors.New("test")

			It("returns errors when the stream is cancelled", func() {
				str.closeForShutdown(testErr)
				n, err := strWithTimeout.Write([]byte("foo"))
				Expect(n).To(BeZero())
				Expect(err).To(MatchError(testErr))
			})

			It("doesn't get data for writing if an error occurred", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				mockFC.EXPECT().SendWindowSize().Return(protocol.ByteCount(9999))
				mockFC.EXPECT().AddBytesSent(gomock.Any())
				done := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					_, err := strWithTimeout.Write(bytes.Repeat([]byte{0}, 500))
					Expect(err).To(MatchError(testErr))
					close(done)
				}()
				waitForWrite()
				frame, hasMoreData := str.popStreamFrame(50) // get a STREAM frame containing some data, but not all
				Expect(frame).ToNot(BeNil())
				Expect(hasMoreData).To(BeTrue())
				str.closeForShutdown(testErr)
				frame, hasMoreData = str.popStreamFrame(1000)
				Expect(frame).To(BeNil())
				Expect(hasMoreData).To(BeFalse())
				Eventually(done).Should(BeClosed())
			})

			It("cancels the context", func() {
				Expect(str.Context().Done()).ToNot(BeClosed())
				str.closeForShutdown(testErr)
				Expect(str.Context().Done()).To(BeClosed())
			})
		})
	})

	Context("handling MAX_STREAM_DATA frames", func() {
		It("informs the flow controller", func() {
			mockFC.EXPECT().UpdateSendWindow(protocol.ByteCount(0x1337))
			str.handleMaxStreamDataFrame(&wire.MaxStreamDataFrame{
				StreamID:   streamID,
				ByteOffset: 0x1337,
			})
		})

		It("says when it has data for sending", func() {
			mockFC.EXPECT().UpdateSendWindow(gomock.Any())
			mockSender.EXPECT().onHasStreamData(streamID).Times(2) // once for Write, once for the MAX_STREAM_DATA frame
			done := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				_, err := str.Write([]byte("foobar"))
				Expect(err).ToNot(HaveOccurred())
				close(done)
			}()
			waitForWrite()
			str.handleMaxStreamDataFrame(&wire.MaxStreamDataFrame{
				StreamID:   streamID,
				ByteOffset: 42,
			})
			// make sure the Write go routine returns
			str.closeForShutdown(nil)
			Eventually(done).Should(BeClosed())
		})
	})

	Context("stream cancelations", func() {
		Context("canceling writing", func() {
			It("queues a RST_STREAM frame", func() {
				mockSender.EXPECT().queueControlFrame(&wire.RstStreamFrame{
					StreamID:   streamID,
					ByteOffset: 1234,
					ErrorCode:  9876,
				})
				mockSender.EXPECT().onStreamCompleted(streamID)
				str.writeOffset = 1234
				err := str.CancelWrite(9876)
				Expect(err).ToNot(HaveOccurred())
			})

			It("unblocks Write", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				mockSender.EXPECT().onStreamCompleted(streamID)
				mockSender.EXPECT().queueControlFrame(gomock.Any())
				mockFC.EXPECT().SendWindowSize().Return(protocol.MaxByteCount)
				mockFC.EXPECT().AddBytesSent(gomock.Any())
				writeReturned := make(chan struct{})
				var n int
				go func() {
					defer GinkgoRecover()
					var err error
					n, err = strWithTimeout.Write(bytes.Repeat([]byte{0}, 100))
					Expect(err).To(MatchError("Write on stream 1337 canceled with error code 1234"))
					close(writeReturned)
				}()
				waitForWrite()
				frame, _ := str.popStreamFrame(50)
				Expect(frame).ToNot(BeNil())
				err := str.CancelWrite(1234)
				Expect(err).ToNot(HaveOccurred())
				Eventually(writeReturned).Should(BeClosed())
				Expect(n).To(BeEquivalentTo(frame.DataLen()))
			})

			It("cancels the context", func() {
				mockSender.EXPECT().queueControlFrame(gomock.Any())
				mockSender.EXPECT().onStreamCompleted(streamID)
				Expect(str.Context().Done()).ToNot(BeClosed())
				str.CancelWrite(1234)
				Expect(str.Context().Done()).To(BeClosed())
			})

			It("doesn't allow further calls to Write", func() {
				mockSender.EXPECT().queueControlFrame(gomock.Any())
				mockSender.EXPECT().onStreamCompleted(streamID)
				err := str.CancelWrite(1234)
				Expect(err).ToNot(HaveOccurred())
				_, err = strWithTimeout.Write([]byte("foobar"))
				Expect(err).To(MatchError("Write on stream 1337 canceled with error code 1234"))
			})

			It("only cancels once", func() {
				mockSender.EXPECT().queueControlFrame(gomock.Any())
				mockSender.EXPECT().onStreamCompleted(streamID)
				err := str.CancelWrite(1234)
				Expect(err).ToNot(HaveOccurred())
				err = str.CancelWrite(4321)
				Expect(err).ToNot(HaveOccurred())
			})

			It("doesn't cancel when the stream was already closed", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				err := str.Close()
				Expect(err).ToNot(HaveOccurred())
				err = str.CancelWrite(123)
				Expect(err).To(MatchError("CancelWrite for closed stream 1337"))
			})
		})

		Context("receiving STOP_SENDING frames", func() {
			It("queues a RST_STREAM frames with error code Stopping", func() {
				mockSender.EXPECT().queueControlFrame(&wire.RstStreamFrame{
					StreamID:  streamID,
					ErrorCode: errorCodeStopping,
				})
				mockSender.EXPECT().onStreamCompleted(streamID)
				str.handleStopSendingFrame(&wire.StopSendingFrame{
					StreamID:  streamID,
					ErrorCode: 101,
				})
			})

			It("unblocks Write", func() {
				mockSender.EXPECT().onHasStreamData(streamID)
				mockSender.EXPECT().queueControlFrame(gomock.Any())
				done := make(chan struct{})
				go func() {
					defer GinkgoRecover()
					_, err := str.Write([]byte("foobar"))
					Expect(err).To(MatchError("Stream 1337 was reset with error code 123"))
					Expect(err).To(BeAssignableToTypeOf(streamCanceledError{}))
					Expect(err.(streamCanceledError).Canceled()).To(BeTrue())
					Expect(err.(streamCanceledError).ErrorCode()).To(Equal(protocol.ApplicationErrorCode(123)))
					close(done)
				}()
				waitForWrite()
				mockSender.EXPECT().onStreamCompleted(streamID)
				str.handleStopSendingFrame(&wire.StopSendingFrame{
					StreamID:  streamID,
					ErrorCode: 123,
				})
				Eventually(done).Should(BeClosed())
			})

			It("doesn't allow further calls to Write", func() {
				mockSender.EXPECT().queueControlFrame(gomock.Any())
				mockSender.EXPECT().onStreamCompleted(streamID)
				str.handleStopSendingFrame(&wire.StopSendingFrame{
					StreamID:  streamID,
					ErrorCode: 123,
				})
				_, err := str.Write([]byte("foobar"))
				Expect(err).To(MatchError("Stream 1337 was reset with error code 123"))
				Expect(err).To(BeAssignableToTypeOf(streamCanceledError{}))
				Expect(err.(streamCanceledError).Canceled()).To(BeTrue())
				Expect(err.(streamCanceledError).ErrorCode()).To(Equal(protocol.ApplicationErrorCode(123)))
			})
		})
	})
})
