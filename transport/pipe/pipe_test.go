package pipe_test

import (
	"io"
	"sync"
	"testing"
	"time"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/task"
	. "v2ray.com/core/transport/pipe"
	. "v2ray.com/ext/assert"
)

func TestPipeReadWrite(t *testing.T) {
	assert := With(t)

	pReader, pWriter := New(WithSizeLimit(1024))
	payload := []byte{'a', 'b', 'c', 'd'}
	b := buf.New()
	b.Write(payload)
	assert(pWriter.WriteMultiBuffer(buf.NewMultiBufferValue(b)), IsNil)

	rb, err := pReader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(rb.String(), Equals, b.String())
}

func TestPipeCloseError(t *testing.T) {
	assert := With(t)

	pReader, pWriter := New(WithSizeLimit(1024))
	payload := []byte{'a', 'b', 'c', 'd'}
	b := buf.New()
	b.Write(payload)
	assert(pWriter.WriteMultiBuffer(buf.NewMultiBufferValue(b)), IsNil)
	pWriter.CloseError()

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
	assert(pWriter.WriteMultiBuffer(buf.NewMultiBufferValue(b)), IsNil)
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
	assert(pWriter.WriteMultiBuffer(buf.NewMultiBufferValue(bb)), IsNil)

	err := task.Run(task.Parallel(func() error {
		b := buf.New()
		b.Write([]byte{'c', 'd'})
		return pWriter.WriteMultiBuffer(buf.NewMultiBufferValue(b))
	}, func() error {
		time.Sleep(time.Second)

		rb, err := pReader.ReadMultiBuffer()
		if err != nil {
			return err
		}
		assert(rb.String(), Equals, "ab")

		rb, err = pReader.ReadMultiBuffer()
		if err != nil {
			return err
		}
		assert(rb.String(), Equals, "cd")
		return nil
	}))()

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
			b.WriteBytes('a', 'b', 'c', 'd')
			pWriter.WriteMultiBuffer(buf.NewMultiBufferValue(b))
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
