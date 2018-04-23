package pipe

import (
	"errors"
	"io"
	"sync"
	"time"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/signal"
)

type state byte

const (
	open state = iota
	closed
	errord
)

type pipe struct {
	sync.Mutex
	data        buf.MultiBuffer
	readSignal  *signal.Notifier
	writeSignal *signal.Notifier
	limit       int32
	state       state
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

		<-p.readSignal.Wait()
	}
}

func (p *pipe) ReadMultiBufferWithTimeout(d time.Duration) (buf.MultiBuffer, error) {
	timer := time.After(d)
	for {
		data, err := p.readMultiBufferInternal()
		if data != nil || err != nil {
			p.writeSignal.Signal()
			return data, err
		}

		select {
		case <-p.readSignal.Wait():
		case <-timer:
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
		if err == nil || err != errBufferFull {
			p.readSignal.Signal()
			return err
		}

		<-p.writeSignal.Wait()
	}
}

func (p *pipe) Close() error {
	p.Lock()
	defer p.Unlock()

	if p.state == closed || p.state == errord {
		return nil
	}

	p.state = closed
	p.readSignal.Signal()
	p.writeSignal.Signal()
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

	p.readSignal.Signal()
	p.writeSignal.Signal()
}
