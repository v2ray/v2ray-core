package pipe

import (
	"errors"
	"io"
	"runtime"
	"sync"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/signal/done"
)

type state byte

const (
	open state = iota
	closed
	errord
)

type pipe struct {
	sync.Mutex
	data            buf.MultiBuffer
	readSignal      *signal.Notifier
	writeSignal     *signal.Notifier
	done            *done.Instance
	limit           int32
	state           state
	discardOverflow bool
}

var errBufferFull = errors.New("buffer full")

func (p *pipe) getState(forRead bool) error {
	switch p.state {
	case open:
		if !forRead && p.limit >= 0 && p.data.Len() > p.limit {
			return errBufferFull
		}
		return nil
	case closed:
		if !forRead {
			return io.ErrClosedPipe
		}
		if !p.data.IsEmpty() {
			return nil
		}
		return io.EOF
	case errord:
		return io.ErrClosedPipe
	default:
		panic("impossible case")
	}
}

func (p *pipe) readMultiBufferInternal() (buf.MultiBuffer, error) {
	p.Lock()
	defer p.Unlock()

	if err := p.getState(true); err != nil {
		return nil, err
	}

	data := p.data
	p.data = nil
	return data, nil
}

func (p *pipe) ReadMultiBuffer() (buf.MultiBuffer, error) {
	for {
		data, err := p.readMultiBufferInternal()
		if data != nil || err != nil {
			p.writeSignal.Signal()
			return data, err
		}

		select {
		case <-p.readSignal.Wait():
		case <-p.done.Wait():
		}
	}
}

func (p *pipe) ReadMultiBufferTimeout(d time.Duration) (buf.MultiBuffer, error) {
	timer := time.NewTimer(d)
	defer timer.Stop()

	for {
		data, err := p.readMultiBufferInternal()
		if data != nil || err != nil {
			p.writeSignal.Signal()
			return data, err
		}

		select {
		case <-p.readSignal.Wait():
		case <-p.done.Wait():
		case <-timer.C:
			return nil, buf.ErrReadTimeout
		}
	}
}

func (p *pipe) writeMultiBufferInternal(mb buf.MultiBuffer) error {
	p.Lock()
	defer p.Unlock()

	if err := p.getState(false); err != nil {
		return err
	}

	p.data.AppendMulti(mb)
	return nil
}

func (p *pipe) WriteMultiBuffer(mb buf.MultiBuffer) error {
	if mb.IsEmpty() {
		return nil
	}

	for {
		err := p.writeMultiBufferInternal(mb)
		switch {
		case err == nil:
			p.readSignal.Signal()

			// Yield current goroutine. Hopefully the reading counterpart can pick up the payload.
			runtime.Gosched()
			return nil
		case err == errBufferFull && p.discardOverflow:
			mb.Release()
			return nil
		case err != errBufferFull:
			mb.Release()
			p.readSignal.Signal()
			return err
		}

		select {
		case <-p.writeSignal.Wait():
		case <-p.done.Wait():
			return io.ErrClosedPipe
		}
	}
}

func (p *pipe) Close() error {
	p.Lock()
	defer p.Unlock()

	if p.state == closed || p.state == errord {
		return nil
	}

	p.state = closed
	common.Must(p.done.Close())
	return nil
}

func (p *pipe) CloseError() {
	p.Lock()
	defer p.Unlock()

	if p.state == closed || p.state == errord {
		return
	}

	p.state = errord

	if !p.data.IsEmpty() {
		p.data.Release()
		p.data = nil
	}

	common.Must(p.done.Close())
}
