package pipe_test

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/task"
	. "v2ray.com/core/transport/pipe"
	. "v2ray.com/ext/assert"
)

func TestPipeReadWrite(t *testing.T) {
	assert := With(t)

	pReader, pWriter := New(WithSizeLimit(1024))

	b := buf.New()
	b.WriteString("abcd")
	assert(pWriter.WriteMultiBuffer(buf.MultiBuffer{b}), IsNil)

	b2 := buf.New()
	b2.WriteString("efg")
	assert(pWriter.WriteMultiBuffer(buf.MultiBuffer{b2}), IsNil)

	rb, err := pReader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(rb.String(), Equals, "abcdefg")
}

func TestPipeInterrupt(t *testing.T) {
	assert := With(t)

	pReader, pWriter := New(WithSizeLimit(1024))
	payload := []byte{'a', 'b', 'c', 'd'}
	b := buf.New()
	b.Write(payload)
	assert(pWriter.WriteMultiBuffer(buf.MultiBuffer{b}), IsNil)
	pWriter.Interrupt()

	rb, err := pReader.ReadMultiBuffer()
	assert(err, Equals, io.ErrClosedPipe)
	assert(rb.IsEmpty(), IsTrue)
}

func TestPipeClose(t *testing.T) {
	assert := With(t)

	pReader, pWriter := New(WithSizeLimit(1024))
	payload := []byte{'a', 'b', 'c', 'd'}
	b := buf.New()
	b.Write(payload)
	assert(pWriter.WriteMultiBuffer(buf.MultiBuffer{b}), IsNil)
	assert(pWriter.Close(), IsNil)

	rb, err := pReader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(rb.String(), Equals, b.String())

	rb, err = pReader.ReadMultiBuffer()
	assert(err, Equals, io.EOF)
	assert(rb.IsEmpty(), IsTrue)
}

func TestPipeLimitZero(t *testing.T) {
	assert := With(t)

	pReader, pWriter := New(WithSizeLimit(0))
	bb := buf.New()
	bb.Write([]byte{'a', 'b'})
	assert(pWriter.WriteMultiBuffer(buf.MultiBuffer{bb}), IsNil)

	err := task.Run(context.Background(), func() error {
		b := buf.New()
		b.Write([]byte{'c', 'd'})
		return pWriter.WriteMultiBuffer(buf.MultiBuffer{b})
	}, func() error {
		time.Sleep(time.Second)

		var container buf.MultiBufferContainer
		if err := buf.Copy(pReader, &container); err != nil {
			return err
		}

		assert(container.String(), Equals, "abcd")
		return nil
	}, func() error {
		time.Sleep(time.Second * 2)
		pWriter.Close()
		return nil
	})

	assert(err, IsNil)
}

func TestPipeWriteMultiThread(t *testing.T) {
	assert := With(t)

	pReader, pWriter := New(WithSizeLimit(0))

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			b := buf.New()
			b.WriteString("abcd")
			pWriter.WriteMultiBuffer(buf.MultiBuffer{b})
			wg.Done()
		}()
	}

	time.Sleep(time.Millisecond * 100)

	pWriter.Close()
	wg.Wait()

	b, err := pReader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(b[0].Bytes(), Equals, []byte{'a', 'b', 'c', 'd'})
}

func TestInterfaces(t *testing.T) {
	assert := With(t)

	assert((*Reader)(nil), Implements, (*buf.Reader)(nil))
	assert((*Reader)(nil), Implements, (*buf.TimeoutReader)(nil))
}

func BenchmarkPipeReadWrite(b *testing.B) {
	reader, writer := New(WithoutSizeLimit())
	a := buf.New()
	a.Extend(buf.Size)
	c := buf.MultiBuffer{a}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		common.Must(writer.WriteMultiBuffer(c))
		d, err := reader.ReadMultiBuffer()
		common.Must(err)
		c = d
	}
}
