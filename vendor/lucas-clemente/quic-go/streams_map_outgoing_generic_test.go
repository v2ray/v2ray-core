package quic

import (
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
	"github.com/lucas-clemente/quic-go/qerr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Streams Map (outgoing)", func() {
	const firstNewStream protocol.StreamID = 10
	var (
		m          *outgoingItemsMap
		newItem    func(id protocol.StreamID) item
		mockSender *MockStreamSender
	)

	BeforeEach(func() {
		newItem = func(id protocol.StreamID) item {
			return &mockGenericStream{id: id}
		}
		mockSender = NewMockStreamSender(mockCtrl)
		m = newOutgoingItemsMap(firstNewStream, newItem, mockSender.queueControlFrame)
	})

	Context("no stream ID limit", func() {
		BeforeEach(func() {
			m.SetMaxStream(0xffffffff)
		})

		It("opens streams", func() {
			str, err := m.OpenStream()
			Expect(err).ToNot(HaveOccurred())
			Expect(str.(*mockGenericStream).id).To(Equal(firstNewStream))
			str, err = m.OpenStream()
			Expect(err).ToNot(HaveOccurred())
			Expect(str.(*mockGenericStream).id).To(Equal(firstNewStream + 4))
		})

		It("doesn't open streams after it has been closed", func() {
			testErr := errors.New("close")
			m.CloseWithError(testErr)
			_, err := m.OpenStream()
			Expect(err).To(MatchError(testErr))
		})

		It("gets streams", func() {
			_, err := m.OpenStream()
			Expect(err).ToNot(HaveOccurred())
			str, err := m.GetStream(firstNewStream)
			Expect(err).ToNot(HaveOccurred())
			Expect(str.(*mockGenericStream).id).To(Equal(firstNewStream))
		})

		It("errors when trying to get a stream that has not yet been opened", func() {
			_, err := m.GetStream(10)
			Expect(err).To(MatchError(qerr.Error(qerr.InvalidStreamID, "peer attempted to open stream 10")))
		})

		It("deletes streams", func() {
			_, err := m.OpenStream() // opens stream 10
			Expect(err).ToNot(HaveOccurred())
			err = m.DeleteStream(10)
			Expect(err).ToNot(HaveOccurred())
			str, err := m.GetStream(10)
			Expect(err).ToNot(HaveOccurred())
			Expect(str).To(BeNil())
		})

		It("errors when deleting a non-existing stream", func() {
			err := m.DeleteStream(1337)
			Expect(err).To(MatchError("Tried to delete unknown stream 1337"))
		})

		It("errors when deleting a stream twice", func() {
			_, err := m.OpenStream() // opens stream 10
			Expect(err).ToNot(HaveOccurred())
			err = m.DeleteStream(10)
			Expect(err).ToNot(HaveOccurred())
			err = m.DeleteStream(10)
			Expect(err).To(MatchError("Tried to delete unknown stream 10"))
		})

		It("closes all streams when CloseWithError is called", func() {
			str1, err := m.OpenStream()
			Expect(err).ToNot(HaveOccurred())
			str2, err := m.OpenStream()
			Expect(err).ToNot(HaveOccurred())
			testErr := errors.New("test err")
			m.CloseWithError(testErr)
			Expect(str1.(*mockGenericStream).closed).To(BeTrue())
			Expect(str1.(*mockGenericStream).closeErr).To(MatchError(testErr))
			Expect(str2.(*mockGenericStream).closed).To(BeTrue())
			Expect(str2.(*mockGenericStream).closeErr).To(MatchError(testErr))
		})
	})

	Context("with stream ID limits", func() {
		It("errors when no stream can be opened immediately", func() {
			mockSender.EXPECT().queueControlFrame(gomock.Any())
			_, err := m.OpenStream()
			Expect(err).To(MatchError(qerr.TooManyOpenStreams))
		})

		It("blocks until a stream can be opened synchronously", func() {
			mockSender.EXPECT().queueControlFrame(gomock.Any())
			done := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				str, err := m.OpenStreamSync()
				Expect(err).ToNot(HaveOccurred())
				Expect(str.(*mockGenericStream).id).To(Equal(firstNewStream))
				close(done)
			}()

			Consistently(done).ShouldNot(BeClosed())
			m.SetMaxStream(firstNewStream)
			Eventually(done).Should(BeClosed())
		})

		It("stops opening synchronously when it is closed", func() {
			mockSender.EXPECT().queueControlFrame(gomock.Any())
			testErr := errors.New("test error")
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

		It("doesn't reduce the stream limit", func() {
			m.SetMaxStream(firstNewStream)
			m.SetMaxStream(firstNewStream - 4)
			str, err := m.OpenStream()
			Expect(err).ToNot(HaveOccurred())
			Expect(str.(*mockGenericStream).id).To(Equal(firstNewStream))
		})

		It("queues a STREAM_ID_BLOCKED frame if no stream can be opened", func() {
			m.SetMaxStream(firstNewStream)
			mockSender.EXPECT().queueControlFrame(&wire.StreamIDBlockedFrame{StreamID: firstNewStream})
			_, err := m.OpenStream()
			Expect(err).ToNot(HaveOccurred())
			_, err = m.OpenStream()
			Expect(err).To(MatchError(qerr.TooManyOpenStreams))
		})

		It("only sends one STREAM_ID_BLOCKED frame for one stream ID", func() {
			m.SetMaxStream(firstNewStream)
			mockSender.EXPECT().queueControlFrame(&wire.StreamIDBlockedFrame{StreamID: firstNewStream})
			_, err := m.OpenStream()
			Expect(err).ToNot(HaveOccurred())
			// try to open a stream twice, but expect only one STREAM_ID_BLOCKED to be sent
			_, err = m.OpenStream()
			Expect(err).To(MatchError(qerr.TooManyOpenStreams))
			_, err = m.OpenStream()
			Expect(err).To(MatchError(qerr.TooManyOpenStreams))
		})
	})
})
