package dispatcher

import (
	"io"
	"time"
)

type Reader struct {
	r io.Reader
	l *RateLimiter
}

type Writer struct {
	w io.Writer
	l *RateLimiter
}

type RateLimiter struct {
	rate  time.Duration
	count int64
	t     time.Time
}

type readSeeker struct {
	Reader
	s io.Seeker
}

func NewRateLimiter(rate int64) *RateLimiter {
	return &RateLimiter{
		rate:  time.Duration(rate),
		count: 0,
		t:     time.Now(),
	}
}

func RateReader(r io.Reader, l *RateLimiter) io.Reader {
	return &Reader{
		r: r,
		l: l,
	}
}

func ReadSeeker(rs io.ReadSeeker, l *RateLimiter) io.ReadSeeker {
	return &readSeeker{
		Reader: Reader{
			r: rs,
			l: l,
		},
		s: rs,
	}
}

func RateWriter(w io.Writer, l *RateLimiter) io.Writer {
	return &Writer{
		w: w,
		l: l,
	}
}

func (l *RateLimiter) RateWait(count int) {
	l.count += int64(count)
	t := time.Duration(l.count)*time.Second/l.rate - time.Since(l.t)
	if t > 0 {
		time.Sleep(t)
	}
}

func (r *Reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	r.l.RateWait(n)
	return n, err
}

func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return rs.s.Seek(offset, whence)
}

func (w *Writer) Write(buf []byte) (int, error) {
	w.l.RateWait(len(buf))
	return w.w.Write(buf)
}
