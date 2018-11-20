package quic

import (
	"github.com/lucas-clemente/quic-go/internal/mocks"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Window Update Queue", func() {
	var (
		q            *windowUpdateQueue
		streamGetter *MockStreamGetter
		connFC       *mocks.MockConnectionFlowController
		queuedFrames []wire.Frame
		cryptoStream *MockCryptoStream
	)

	BeforeEach(func() {
		streamGetter = NewMockStreamGetter(mockCtrl)
		cryptoStream = NewMockCryptoStream(mockCtrl)
		connFC = mocks.NewMockConnectionFlowController(mockCtrl)
		cryptoStream.EXPECT().StreamID().Return(protocol.StreamID(0)).AnyTimes()
		queuedFrames = queuedFrames[:0]
		q = newWindowUpdateQueue(streamGetter, cryptoStream, connFC, func(f wire.Frame) {
			queuedFrames = append(queuedFrames, f)
		})
	})

	It("adds stream offsets and gets MAX_STREAM_DATA frames", func() {
		stream1 := NewMockStreamI(mockCtrl)
		stream1.EXPECT().getWindowUpdate().Return(protocol.ByteCount(10))
		stream3 := NewMockStreamI(mockCtrl)
		stream3.EXPECT().getWindowUpdate().Return(protocol.ByteCount(30))
		streamGetter.EXPECT().GetOrOpenReceiveStream(protocol.StreamID(3)).Return(stream3, nil)
		streamGetter.EXPECT().GetOrOpenReceiveStream(protocol.StreamID(1)).Return(stream1, nil)
		q.AddStream(3)
		q.AddStream(1)
		q.QueueAll()
		Expect(queuedFrames).To(ContainElement(&wire.MaxStreamDataFrame{StreamID: 1, ByteOffset: 10}))
		Expect(queuedFrames).To(ContainElement(&wire.MaxStreamDataFrame{StreamID: 3, ByteOffset: 30}))
	})

	It("deletes the entry after getting the MAX_STREAM_DATA frame", func() {
		stream10 := NewMockStreamI(mockCtrl)
		stream10.EXPECT().getWindowUpdate().Return(protocol.ByteCount(100))
		streamGetter.EXPECT().GetOrOpenReceiveStream(protocol.StreamID(10)).Return(stream10, nil)
		q.AddStream(10)
		q.QueueAll()
		Expect(queuedFrames).To(HaveLen(1))
		q.QueueAll()
		Expect(queuedFrames).To(HaveLen(1))
	})

	It("doesn't queue a MAX_STREAM_DATA for a closed stream", func() {
		streamGetter.EXPECT().GetOrOpenReceiveStream(protocol.StreamID(12)).Return(nil, nil)
		q.AddStream(12)
		q.QueueAll()
		Expect(queuedFrames).To(BeEmpty())
	})

	It("doesn't queue a MAX_STREAM_DATA if the flow controller returns an offset of 0", func() {
		stream5 := NewMockStreamI(mockCtrl)
		stream5.EXPECT().getWindowUpdate().Return(protocol.ByteCount(0))
		streamGetter.EXPECT().GetOrOpenReceiveStream(protocol.StreamID(5)).Return(stream5, nil)
		q.AddStream(5)
		q.QueueAll()
		Expect(queuedFrames).To(BeEmpty())
	})

	It("adds MAX_STREAM_DATA frames for the crypto stream", func() {
		cryptoStream.EXPECT().getWindowUpdate().Return(protocol.ByteCount(42))
		q.AddStream(0)
		q.QueueAll()
		Expect(queuedFrames).To(Equal([]wire.Frame{
			&wire.MaxStreamDataFrame{StreamID: 0, ByteOffset: 42},
		}))
	})

	It("queues MAX_DATA frames", func() {
		connFC.EXPECT().GetWindowUpdate().Return(protocol.ByteCount(0x1337))
		q.AddConnection()
		q.QueueAll()
		Expect(queuedFrames).To(Equal([]wire.Frame{
			&wire.MaxDataFrame{ByteOffset: 0x1337},
		}))
	})

	It("deduplicates", func() {
		stream10 := NewMockStreamI(mockCtrl)
		stream10.EXPECT().getWindowUpdate().Return(protocol.ByteCount(200))
		streamGetter.EXPECT().GetOrOpenReceiveStream(protocol.StreamID(10)).Return(stream10, nil)
		q.AddStream(10)
		q.AddStream(10)
		q.QueueAll()
		Expect(queuedFrames).To(Equal([]wire.Frame{
			&wire.MaxStreamDataFrame{StreamID: 10, ByteOffset: 200},
		}))
	})
})
