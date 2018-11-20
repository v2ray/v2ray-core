package quic

import (
	"bytes"

	"github.com/golang/mock/gomock"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stream Framer", func() {
	const (
		id1 = protocol.StreamID(10)
		id2 = protocol.StreamID(11)
	)

	var (
		framer           *streamFramer
		cryptoStream     *MockCryptoStream
		stream1, stream2 *MockSendStreamI
		streamGetter     *MockStreamGetter
	)

	BeforeEach(func() {
		streamGetter = NewMockStreamGetter(mockCtrl)
		stream1 = NewMockSendStreamI(mockCtrl)
		stream1.EXPECT().StreamID().Return(protocol.StreamID(5)).AnyTimes()
		stream2 = NewMockSendStreamI(mockCtrl)
		stream2.EXPECT().StreamID().Return(protocol.StreamID(6)).AnyTimes()
		cryptoStream = NewMockCryptoStream(mockCtrl)
		framer = newStreamFramer(cryptoStream, streamGetter, versionGQUICFrames)
	})

	Context("handling the crypto stream", func() {
		It("says if it has crypto stream data", func() {
			Expect(framer.HasCryptoStreamData()).To(BeFalse())
			framer.AddActiveStream(framer.version.CryptoStreamID())
			Expect(framer.HasCryptoStreamData()).To(BeTrue())
		})

		It("says that it doesn't have crypto stream data after popping all data", func() {
			streamID := framer.version.CryptoStreamID()
			f := &wire.StreamFrame{
				StreamID: streamID,
				Data:     []byte("foobar"),
			}
			cryptoStream.EXPECT().popStreamFrame(protocol.ByteCount(1000)).Return(f, false)
			framer.AddActiveStream(streamID)
			Expect(framer.PopCryptoStreamFrame(1000)).To(Equal(f))
			Expect(framer.HasCryptoStreamData()).To(BeFalse())
		})

		It("says that it has more crypto stream data if not all data was popped", func() {
			streamID := framer.version.CryptoStreamID()
			f := &wire.StreamFrame{
				StreamID: streamID,
				Data:     []byte("foobar"),
			}
			cryptoStream.EXPECT().popStreamFrame(protocol.ByteCount(1000)).Return(f, true)
			framer.AddActiveStream(streamID)
			Expect(framer.PopCryptoStreamFrame(1000)).To(Equal(f))
			Expect(framer.HasCryptoStreamData()).To(BeTrue())
		})
	})

	Context("Popping", func() {
		It("returns nil when popping an empty framer", func() {
			Expect(framer.PopStreamFrames(1000)).To(BeEmpty())
		})

		It("returns STREAM frames", func() {
			streamGetter.EXPECT().GetOrOpenSendStream(id1).Return(stream1, nil)
			f := &wire.StreamFrame{
				StreamID: id1,
				Data:     []byte("foobar"),
				Offset:   42,
			}
			stream1.EXPECT().popStreamFrame(gomock.Any()).Return(f, false)
			framer.AddActiveStream(id1)
			fs := framer.PopStreamFrames(1000)
			Expect(fs).To(Equal([]*wire.StreamFrame{f}))
		})

		It("skips a stream that was reported active, but was completed shortly after", func() {
			streamGetter.EXPECT().GetOrOpenSendStream(id1).Return(nil, nil)
			streamGetter.EXPECT().GetOrOpenSendStream(id2).Return(stream2, nil)
			f := &wire.StreamFrame{
				StreamID: id2,
				Data:     []byte("foobar"),
			}
			stream2.EXPECT().popStreamFrame(gomock.Any()).Return(f, false)
			framer.AddActiveStream(id1)
			framer.AddActiveStream(id2)
			Expect(framer.PopStreamFrames(1000)).To(Equal([]*wire.StreamFrame{f}))
		})

		It("skips a stream that was reported active, but doesn't have any data", func() {
			streamGetter.EXPECT().GetOrOpenSendStream(id1).Return(stream1, nil)
			streamGetter.EXPECT().GetOrOpenSendStream(id2).Return(stream2, nil)
			f := &wire.StreamFrame{
				StreamID: id2,
				Data:     []byte("foobar"),
			}
			stream1.EXPECT().popStreamFrame(gomock.Any()).Return(nil, false)
			stream2.EXPECT().popStreamFrame(gomock.Any()).Return(f, false)
			framer.AddActiveStream(id1)
			framer.AddActiveStream(id2)
			Expect(framer.PopStreamFrames(1000)).To(Equal([]*wire.StreamFrame{f}))
		})

		It("pops from a stream multiple times, if it has enough data", func() {
			streamGetter.EXPECT().GetOrOpenSendStream(id1).Return(stream1, nil).Times(2)
			f1 := &wire.StreamFrame{StreamID: id1, Data: []byte("foobar")}
			f2 := &wire.StreamFrame{StreamID: id1, Data: []byte("foobaz")}
			stream1.EXPECT().popStreamFrame(gomock.Any()).Return(f1, true)
			stream1.EXPECT().popStreamFrame(gomock.Any()).Return(f2, false)
			framer.AddActiveStream(id1) // only add it once
			Expect(framer.PopStreamFrames(protocol.MinStreamFrameSize)).To(Equal([]*wire.StreamFrame{f1}))
			Expect(framer.PopStreamFrames(protocol.MinStreamFrameSize)).To(Equal([]*wire.StreamFrame{f2}))
			// no further calls to popStreamFrame, after popStreamFrame said there's no more data
			Expect(framer.PopStreamFrames(protocol.MinStreamFrameSize)).To(BeNil())
		})

		It("re-queues a stream at the end, if it has enough data", func() {
			streamGetter.EXPECT().GetOrOpenSendStream(id1).Return(stream1, nil).Times(2)
			streamGetter.EXPECT().GetOrOpenSendStream(id2).Return(stream2, nil)
			f11 := &wire.StreamFrame{StreamID: id1, Data: []byte("foobar")}
			f12 := &wire.StreamFrame{StreamID: id1, Data: []byte("foobaz")}
			f2 := &wire.StreamFrame{StreamID: id2, Data: []byte("raboof")}
			stream1.EXPECT().popStreamFrame(gomock.Any()).Return(f11, true)
			stream1.EXPECT().popStreamFrame(gomock.Any()).Return(f12, false)
			stream2.EXPECT().popStreamFrame(gomock.Any()).Return(f2, false)
			framer.AddActiveStream(id1) // only add it once
			framer.AddActiveStream(id2)
			Expect(framer.PopStreamFrames(protocol.MinStreamFrameSize)).To(Equal([]*wire.StreamFrame{f11})) // first a frame from stream 1
			Expect(framer.PopStreamFrames(protocol.MinStreamFrameSize)).To(Equal([]*wire.StreamFrame{f2}))  // then a frame from stream 2
			Expect(framer.PopStreamFrames(protocol.MinStreamFrameSize)).To(Equal([]*wire.StreamFrame{f12})) // then another frame from stream 1
		})

		It("only dequeues data from each stream once per packet", func() {
			streamGetter.EXPECT().GetOrOpenSendStream(id1).Return(stream1, nil)
			streamGetter.EXPECT().GetOrOpenSendStream(id2).Return(stream2, nil)
			f1 := &wire.StreamFrame{StreamID: id1, Data: []byte("foobar")}
			f2 := &wire.StreamFrame{StreamID: id2, Data: []byte("raboof")}
			// both streams have more data, and will be re-queued
			stream1.EXPECT().popStreamFrame(gomock.Any()).Return(f1, true)
			stream2.EXPECT().popStreamFrame(gomock.Any()).Return(f2, true)
			framer.AddActiveStream(id1)
			framer.AddActiveStream(id2)
			Expect(framer.PopStreamFrames(1000)).To(Equal([]*wire.StreamFrame{f1, f2}))
		})

		It("returns multiple normal frames in the order they were reported active", func() {
			streamGetter.EXPECT().GetOrOpenSendStream(id1).Return(stream1, nil)
			streamGetter.EXPECT().GetOrOpenSendStream(id2).Return(stream2, nil)
			f1 := &wire.StreamFrame{Data: []byte("foobar")}
			f2 := &wire.StreamFrame{Data: []byte("foobaz")}
			stream1.EXPECT().popStreamFrame(gomock.Any()).Return(f1, false)
			stream2.EXPECT().popStreamFrame(gomock.Any()).Return(f2, false)
			framer.AddActiveStream(id2)
			framer.AddActiveStream(id1)
			Expect(framer.PopStreamFrames(1000)).To(Equal([]*wire.StreamFrame{f2, f1}))
		})

		It("only asks a stream for data once, even if it was reported active multiple times", func() {
			streamGetter.EXPECT().GetOrOpenSendStream(id1).Return(stream1, nil)
			f := &wire.StreamFrame{Data: []byte("foobar")}
			stream1.EXPECT().popStreamFrame(gomock.Any()).Return(f, false) // only one call to this function
			framer.AddActiveStream(id1)
			framer.AddActiveStream(id1)
			Expect(framer.PopStreamFrames(1000)).To(HaveLen(1))
		})

		It("does not pop empty frames", func() {
			fs := framer.PopStreamFrames(500)
			Expect(fs).To(BeEmpty())
		})

		It("pops frames that have the minimum size", func() {
			streamGetter.EXPECT().GetOrOpenSendStream(id1).Return(stream1, nil)
			stream1.EXPECT().popStreamFrame(protocol.MinStreamFrameSize).Return(&wire.StreamFrame{Data: []byte("foobar")}, false)
			framer.AddActiveStream(id1)
			framer.PopStreamFrames(protocol.MinStreamFrameSize)
		})

		It("does not pop frames smaller than the minimum size", func() {
			// don't expect a call to PopStreamFrame()
			framer.PopStreamFrames(protocol.MinStreamFrameSize - 1)
		})

		It("stops iterating when the remaining size is smaller than the minimum STREAM frame size", func() {
			streamGetter.EXPECT().GetOrOpenSendStream(id1).Return(stream1, nil)
			// pop a frame such that the remaining size is one byte less than the minimum STREAM frame size
			f := &wire.StreamFrame{
				StreamID: id1,
				Data:     bytes.Repeat([]byte("f"), int(500-protocol.MinStreamFrameSize)),
			}
			stream1.EXPECT().popStreamFrame(protocol.ByteCount(500)).Return(f, false)
			framer.AddActiveStream(id1)
			fs := framer.PopStreamFrames(500)
			Expect(fs).To(Equal([]*wire.StreamFrame{f}))
		})
	})
})
