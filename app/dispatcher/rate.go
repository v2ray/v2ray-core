package dispatcher

import (
	"time"

	"v2ray.com/core/common/buf"
)

type Writer struct {
	writer  buf.Writer
	limiter *RateLimiter
}

type RateLimiter struct {
	rate  time.Duration
	count int64
	t     time.Time
}

func NewRateLimiter(rate int64) *RateLimiter {
	return &RateLimiter{
		rate:  time.Duration(rate),
		count: 0,
		t:     time.Now(),
	}
}

func RateWriter(writer buf.Writer, limiter *RateLimiter) buf.Writer {
	return &Writer{
		writer:  writer,
		limiter: limiter,
	}
}

func (l *RateLimiter) RateWait(count int64) {
	l.count += count
	t := time.Duration(l.count)*time.Second/l.rate - time.Since(l.t)
	if t > 0 {
		time.Sleep(t)
	}
}

func (w *Writer) WriteMultiBuffer(mb buf.MultiBuffer) error {
	w.limiter.RateWait(int64(mb.Len()))
	return w.writer.WriteMultiBuffer(mb)
}
