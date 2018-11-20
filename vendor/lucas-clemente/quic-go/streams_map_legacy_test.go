package quic

import (
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/lucas-clemente/quic-go/internal/handshake"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"github.com/lucas-clemente/quic-go/qerr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Streams Map (for gQUIC)", func() {
	var m *streamsMapLegacy

	newStream := func(id protocol.StreamID) streamI {
		str := NewMockStreamI(mockCtrl)
		str.EXPECT().StreamID().Return(id).AnyTimes()
		return str
	}

	setNewStreamsMap := func(p protocol.Perspective) {
		m = newStreamsMapLegacy(newStream, protocol.DefaultMaxIncomingStreams, p).(*streamsMapLegacy)
	}

	deleteStream := func(id protocol.StreamID) {
		ExpectWithOffset(1, m.DeleteStream(id)).To(Succeed())
	}

	It("applies the max stream limit for small number of streams", func() {
		sm := newStreamsMapLegacy(newStream, 1, protocol.PerspectiveServer).(*streamsMapLegacy)
		Expect(sm.maxIncomingStreams).To(BeEquivalentTo(1 + protocol.MaxStreamsMinimumIncrement))
	})

	It("applies the max stream limit for big number of streams", func() {
		sm := newStreamsMapLegacy(newStream, 1000, protocol.PerspectiveServer).(*streamsMapLegacy)
		Expect(sm.maxIncomingStreams).To(BeEquivalentTo(1000 * protocol.MaxStreamsMultiplier))
	})

	Context("getting and creating streams", func() {
		Context("as a server", func() {
			BeforeEach(func() {
				setNewStreamsMap(protocol.PerspectiveServer)
			})

			Context("client-side streams", func() {
				It("gets new streams", func() {
					s, err := m.getOrOpenStream(3)
					Expect(err).NotTo(HaveOccurred())
					Expect(s).ToNot(BeNil())
					Expect(s.StreamID()).To(Equal(protocol.StreamID(3)))
					Expect(m.streams).To(HaveLen(1))
					Expect(m.numIncomingStreams).To(BeEquivalentTo(1))
					Expect(m.numOutgoingStreams).To(BeZero())
				})

				It("rejects streams with even IDs", func() {
					_, err := m.getOrOpenStream(6)
					Expect(err).To(MatchError("InvalidStreamID: peer attempted to open stream 6"))
				})

				It("rejects streams with even IDs, which are lower thatn the highest client-side stream", func() {
					_, err := m.getOrOpenStream(5)
					Expect(err).NotTo(HaveOccurred())
					_, err = m.getOrOpenStream(4)
					Expect(err).To(MatchError("InvalidStreamID: peer attempted to open stream 4"))
				})

				It("gets existing streams", func() {
					s, err := m.getOrOpenStream(5)
					Expect(err).NotTo(HaveOccurred())
					Expect(s.StreamID()).To(Equal(protocol.StreamID(5)))
					numStreams := m.numIncomingStreams
					s, err = m.getOrOpenStream(5)
					Expect(err).NotTo(HaveOccurred())
					Expect(s.StreamID()).To(Equal(protocol.StreamID(5)))
					Expect(m.numIncomingStreams).To(Equal(numStreams))
				})

				It("returns nil for closed streams", func() {
					_, err := m.getOrOpenStream(5)
					Expect(err).NotTo(HaveOccurred())
					deleteStream(5)
					s, err := m.getOrOpenStream(5)
					Expect(err).NotTo(HaveOccurred())
					Expect(s).To(BeNil())
				})

				It("opens skipped streams", func() {
					_, err := m.getOrOpenStream(7)
					Expect(err).NotTo(HaveOccurred())
					Expect(m.streams).To(HaveKey(protocol.StreamID(3)))
					Expect(m.streams).To(HaveKey(protocol.StreamID(5)))
					Expect(m.streams).To(HaveKey(protocol.StreamID(7)))
				})

				It("doesn't reopen an already closed stream", func() {
					_, err := m.getOrOpenStream(5)
					Expect(err).ToNot(HaveOccurred())
					deleteStream(5)
					Expect(err).ToNot(HaveOccurred())
					str, err := m.getOrOpenStream(5)
					Expect(err).ToNot(HaveOccurred())
					Expect(str).To(BeNil())
				})

				Context("counting streams", func() {
					It("errors when too many streams are opened", func() {
						for i := uint32(0); i < m.maxIncomingStreams; i++ {
							_, err := m.getOrOpenStream(protocol.StreamID(i*2 + 1))
							Expect(err).NotTo(HaveOccurred())
						}
						_, err := m.getOrOpenStream(protocol.StreamID(2*m.maxIncomingStreams + 3))
						Expect(err).To(MatchError(qerr.TooManyOpenStreams))
					})

					It("errors when too many streams are opened implicitly", func() {
						_, err := m.getOrOpenStream(protocol.StreamID(m.maxIncomingStreams*2 + 3))
						Expect(err).To(MatchError(qerr.TooManyOpenStreams))
					})

					It("does not error when many streams are opened and closed", func() {
						for i := uint32(2); i < 10*m.maxIncomingStreams; i++ {
							str, err := m.getOrOpenStream(protocol.StreamID(i*2 + 1))
							Expect(err).NotTo(HaveOccurred())
							deleteStream(str.StreamID())
						}
					})
				})
			})

			Context("server-side streams", func() {
				It("doesn't allow opening streams before receiving the transport parameters", func() {
					_, err := m.OpenStream()
					Expect(err).To(MatchError(qerr.TooManyOpenStreams))
				})

				It("opens a stream 2 first", func() {
					m.UpdateLimits(&handshake.TransportParameters{MaxStreams: 10000})
					s, err := m.OpenStream()
					Expect(err).ToNot(HaveOccurred())
					Expect(s).ToNot(BeNil())
					Expect(s.StreamID()).To(Equal(protocol.StreamID(2)))
					Expect(m.numIncomingStreams).To(BeZero())
					Expect(m.numOutgoingStreams).To(BeEquivalentTo(1))
				})

				It("returns the error when the streamsMap was closed", func() {
					testErr := errors.New("test error")
					m.CloseWithError(testErr)
					_, err := m.OpenStream()
					Expect(err).To(MatchError(testErr))
				})

				It("doesn't reopen an already closed stream", func() {
					m.UpdateLimits(&handshake.TransportParameters{MaxStreams: 10000})
					str, err := m.OpenStream()
					Expect(err).ToNot(HaveOccurred())
					Expect(str.StreamID()).To(Equal(protocol.StreamID(2)))
					deleteStream(2)
					Expect(err).ToNot(HaveOccurred())
					str, err = m.getOrOpenStream(2)
					Expect(err).ToNot(HaveOccurred())
					Expect(str).To(BeNil())
				})

				Context("counting streams", func() {
					const maxOutgoingStreams = 50

					BeforeEach(func() {
						m.UpdateLimits(&handshake.TransportParameters{MaxStreams: maxOutgoingStreams})
					})

					It("errors when too many streams are opened", func() {
						for i := 1; i <= maxOutgoingStreams; i++ {
							_, err := m.OpenStream()
							Expect(err).NotTo(HaveOccurred())
						}
						_, err := m.OpenStream()
						Expect(err).To(MatchError(qerr.TooManyOpenStreams))
					})

					It("does not error when many streams are opened and closed", func() {
						for i := 2; i < 10*maxOutgoingStreams; i++ {
							str, err := m.OpenStream()
							Expect(err).NotTo(HaveOccurred())
							deleteStream(str.StreamID())
						}
					})

					It("allows many server- and client-side streams at the same time", func() {
						for i := 1; i < maxOutgoingStreams; i++ {
							_, err := m.OpenStream()
							Expect(err).ToNot(HaveOccurred())
						}
						for i := 0; i < maxOutgoingStreams; i++ {
							_, err := m.getOrOpenStream(protocol.StreamID(2*i + 1))
							Expect(err).ToNot(HaveOccurred())
						}
					})
				})

				Context("opening streams synchronously", func() {
					const maxOutgoingStreams = 10

					BeforeEach(func() {
						m.UpdateLimits(&handshake.TransportParameters{MaxStreams: maxOutgoingStreams})
					})

					openMaxNumStreams := func() {
						for i := 1; i <= maxOutgoingStreams; i++ {
							_, err := m.OpenStream()
							Expect(err).NotTo(HaveOccurred())
						}
						_, err := m.OpenStream()
						Expect(err).To(MatchError(qerr.TooManyOpenStreams))
					}

					It("waits until another stream is closed", func() {
						openMaxNumStreams()
						var str Stream
						done := make(chan struct{})
						go func() {
							defer GinkgoRecover()
							var err error
							str, err = m.OpenStreamSync()
							Expect(err).ToNot(HaveOccurred())
							close(done)
						}()
						Consistently(done).ShouldNot(BeClosed())
						deleteStream(6)
						Eventually(done).Should(BeClosed())
						Expect(str.StreamID()).To(Equal(protocol.StreamID(2*maxOutgoingStreams + 2)))
					})

					It("stops waiting when an error is registered", func() {
						testErr := errors.New("test error")
						openMaxNumStreams()
						for _, str := range m.streams {
							str.(*MockStreamI).EXPECT().closeForShutdown(testErr)
						}

						done := make(chan struct{})
						go func() {
							defer GinkgoRecover()
							_, err := m.OpenStreamSync()
							Expect(err).To(MatchError(testErr))
							close(done)
						}()

						Consistently(done).ShouldNot(BeClosed())
						m.CloseWithError(testErr)
						Eventually(done).Should(BeClosed())
					})

					It("immediately returns when OpenStreamSync is called after an error was registered", func() {
						testErr := errors.New("test error")
						m.CloseWithError(testErr)
						_, err := m.OpenStreamSync()
						Expect(err).To(MatchError(testErr))
					})
				})
			})

			Context("accepting streams", func() {
				It("does nothing if no stream is opened", func() {
					var accepted bool
					go func() {
						_, _ = m.AcceptStream()
						accepted = true
					}()
					Consistently(func() bool { return accepted }).Should(BeFalse())
				})

				It("starts with stream 3", func() {
					var str Stream
					done := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						var err error
						str, err = m.AcceptStream()
						Expect(err).ToNot(HaveOccurred())
						close(done)
					}()
					_, err := m.getOrOpenStream(3)
					Expect(err).ToNot(HaveOccurred())
					Eventually(done).Should(BeClosed())
					Expect(str.StreamID()).To(Equal(protocol.StreamID(3)))
				})

				It("returns an implicitly opened stream, if a stream number is skipped", func() {
					var str Stream
					done := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						var err error
						str, err = m.AcceptStream()
						Expect(err).ToNot(HaveOccurred())
						close(done)
					}()
					_, err := m.getOrOpenStream(5)
					Expect(err).ToNot(HaveOccurred())
					Eventually(done).Should(BeClosed())
					Expect(str.StreamID()).To(Equal(protocol.StreamID(3)))
				})

				It("returns to multiple accepts", func() {
					var str1, str2 Stream
					done1 := make(chan struct{})
					done2 := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						var err error
						str1, err = m.AcceptStream()
						Expect(err).ToNot(HaveOccurred())
						close(done1)
					}()
					go func() {
						defer GinkgoRecover()
						var err error
						str2, err = m.AcceptStream()
						Expect(err).ToNot(HaveOccurred())
						close(done2)
					}()
					_, err := m.getOrOpenStream(5) // opens stream 3 and 5
					Expect(err).ToNot(HaveOccurred())
					Eventually(done1).Should(BeClosed())
					Eventually(done2).Should(BeClosed())
					Expect(str1.StreamID()).ToNot(Equal(str2.StreamID()))
					Expect(str1.StreamID() + str2.StreamID()).To(BeEquivalentTo(3 + 5))
				})

				It("waits until a new stream is available", func() {
					var str Stream
					done := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						var err error
						str, err = m.AcceptStream()
						Expect(err).ToNot(HaveOccurred())
						close(done)
					}()
					Consistently(done).ShouldNot(BeClosed())
					_, err := m.getOrOpenStream(3)
					Expect(err).ToNot(HaveOccurred())
					Eventually(done).Should(BeClosed())
					Expect(str.StreamID()).To(Equal(protocol.StreamID(3)))
				})

				It("returns multiple streams on subsequent Accept calls, if available", func() {
					var str Stream
					done := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						var err error
						str, err = m.AcceptStream()
						Expect(err).ToNot(HaveOccurred())
						close(done)
					}()
					_, err := m.getOrOpenStream(5)
					Expect(err).ToNot(HaveOccurred())
					Eventually(done).Should(BeClosed())
					Expect(str.StreamID()).To(Equal(protocol.StreamID(3)))
					str, err = m.AcceptStream()
					Expect(err).ToNot(HaveOccurred())
					Expect(str.StreamID()).To(Equal(protocol.StreamID(5)))
				})

				It("blocks after accepting a stream", func() {
					_, err := m.getOrOpenStream(3)
					Expect(err).ToNot(HaveOccurred())
					str, err := m.AcceptStream()
					Expect(err).ToNot(HaveOccurred())
					Expect(str.StreamID()).To(Equal(protocol.StreamID(3)))
					done := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						_, _ = m.AcceptStream()
						close(done)
					}()
					Consistently(done).ShouldNot(BeClosed())
					// make the go routine return
					str.(*MockStreamI).EXPECT().closeForShutdown(gomock.Any())
					m.CloseWithError(errors.New("shut down"))
					Eventually(done).Should(BeClosed())
				})

				It("stops waiting when an error is registered", func() {
					testErr := errors.New("testErr")
					done := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						_, err := m.AcceptStream()
						Expect(err).To(MatchError(testErr))
						close(done)
					}()
					Consistently(done).ShouldNot(BeClosed())
					m.CloseWithError(testErr)
					Eventually(done).Should(BeClosed())
				})
				It("immediately returns when Accept is called after an error was registered", func() {
					testErr := errors.New("testErr")
					m.CloseWithError(testErr)
					_, err := m.AcceptStream()
					Expect(err).To(MatchError(testErr))
				})
			})
		})

		Context("as a client", func() {
			BeforeEach(func() {
				setNewStreamsMap(protocol.PerspectiveClient)
				m.UpdateLimits(&handshake.TransportParameters{MaxStreams: 10000})
			})

			Context("server-side streams", func() {
				It("rejects streams with odd IDs", func() {
					_, err := m.getOrOpenStream(5)
					Expect(err).To(MatchError("InvalidStreamID: peer attempted to open stream 5"))
				})

				It("rejects streams with odds IDs, which are lower than the highest server-side stream", func() {
					_, err := m.getOrOpenStream(6)
					Expect(err).NotTo(HaveOccurred())
					_, err = m.getOrOpenStream(5)
					Expect(err).To(MatchError("InvalidStreamID: peer attempted to open stream 5"))
				})

				It("gets new streams", func() {
					s, err := m.getOrOpenStream(2)
					Expect(err).NotTo(HaveOccurred())
					Expect(s.StreamID()).To(Equal(protocol.StreamID(2)))
					Expect(m.streams).To(HaveLen(1))
					Expect(m.numOutgoingStreams).To(BeZero())
					Expect(m.numIncomingStreams).To(BeEquivalentTo(1))
				})

				It("opens skipped streams", func() {
					_, err := m.getOrOpenStream(6)
					Expect(err).NotTo(HaveOccurred())
					Expect(m.streams).To(HaveKey(protocol.StreamID(2)))
					Expect(m.streams).To(HaveKey(protocol.StreamID(4)))
					Expect(m.streams).To(HaveKey(protocol.StreamID(6)))
					Expect(m.numOutgoingStreams).To(BeZero())
					Expect(m.numIncomingStreams).To(BeEquivalentTo(3))
				})

				It("doesn't reopen an already closed stream", func() {
					str, err := m.OpenStream()
					Expect(err).ToNot(HaveOccurred())
					Expect(str.StreamID()).To(Equal(protocol.StreamID(3)))
					deleteStream(3)
					Expect(err).ToNot(HaveOccurred())
					str, err = m.getOrOpenStream(3)
					Expect(err).ToNot(HaveOccurred())
					Expect(str).To(BeNil())
				})
			})

			Context("client-side streams", func() {
				It("starts with stream 3", func() {
					s, err := m.OpenStream()
					Expect(err).ToNot(HaveOccurred())
					Expect(s).ToNot(BeNil())
					Expect(s.StreamID()).To(BeEquivalentTo(3))
					Expect(m.numOutgoingStreams).To(BeEquivalentTo(1))
					Expect(m.numIncomingStreams).To(BeZero())
				})

				It("opens multiple streams", func() {
					s1, err := m.OpenStream()
					Expect(err).ToNot(HaveOccurred())
					s2, err := m.OpenStream()
					Expect(err).ToNot(HaveOccurred())
					Expect(s2.StreamID()).To(Equal(s1.StreamID() + 2))
				})

				It("doesn't reopen an already closed stream", func() {
					_, err := m.getOrOpenStream(4)
					Expect(err).ToNot(HaveOccurred())
					deleteStream(4)
					Expect(err).ToNot(HaveOccurred())
					str, err := m.getOrOpenStream(4)
					Expect(err).ToNot(HaveOccurred())
					Expect(str).To(BeNil())
				})
			})

			Context("accepting streams", func() {
				It("accepts stream 2 first", func() {
					var str Stream
					done := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						var err error
						str, err = m.AcceptStream()
						Expect(err).ToNot(HaveOccurred())
						close(done)
					}()
					_, err := m.getOrOpenStream(2)
					Expect(err).ToNot(HaveOccurred())
					Eventually(done).Should(BeClosed())
					Expect(str.StreamID()).To(Equal(protocol.StreamID(2)))
				})
			})
		})
	})

	Context("deleting streams", func() {
		BeforeEach(func() {
			setNewStreamsMap(protocol.PerspectiveServer)
		})

		It("deletes an incoming stream", func() {
			_, err := m.getOrOpenStream(5) // open stream 3 and 5
			Expect(err).ToNot(HaveOccurred())
			Expect(m.numIncomingStreams).To(BeEquivalentTo(2))
			err = m.DeleteStream(3)
			Expect(err).ToNot(HaveOccurred())
			Expect(m.streams).To(HaveLen(1))
			Expect(m.streams).To(HaveKey(protocol.StreamID(5)))
			Expect(m.numIncomingStreams).To(BeEquivalentTo(1))
		})

		It("deletes an outgoing stream", func() {
			m.UpdateLimits(&handshake.TransportParameters{MaxStreams: 10000})
			_, err := m.OpenStream() // open stream 2
			Expect(err).ToNot(HaveOccurred())
			_, err = m.OpenStream()
			Expect(err).ToNot(HaveOccurred())
			Expect(m.numOutgoingStreams).To(BeEquivalentTo(2))
			err = m.DeleteStream(2)
			Expect(err).ToNot(HaveOccurred())
			Expect(m.numOutgoingStreams).To(BeEquivalentTo(1))
		})

		It("errors when the stream doesn't exist", func() {
			err := m.DeleteStream(1337)
			Expect(err).To(MatchError(errMapAccess))
		})
	})

	It("sets the flow control limit", func() {
		setNewStreamsMap(protocol.PerspectiveServer)
		_, err := m.getOrOpenStream(5)
		Expect(err).ToNot(HaveOccurred())
		m.streams[3].(*MockStreamI).EXPECT().handleMaxStreamDataFrame(&wire.MaxStreamDataFrame{
			StreamID:   3,
			ByteOffset: 321,
		})
		m.streams[5].(*MockStreamI).EXPECT().handleMaxStreamDataFrame(&wire.MaxStreamDataFrame{
			StreamID:   5,
			ByteOffset: 321,
		})
		m.UpdateLimits(&handshake.TransportParameters{StreamFlowControlWindow: 321})
	})

	It("doesn't accept MAX_STREAM_ID frames", func() {
		Expect(m.HandleMaxStreamIDFrame(&wire.MaxStreamIDFrame{})).ToNot(Succeed())
	})
})
