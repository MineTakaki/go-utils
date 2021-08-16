package text

import "io"

type (
	bomWriter struct {
		w io.Writer
		n int64
	}
)

func WithBom(w io.Writer) io.Writer {
	return &bomWriter{w: w}
}

func (w *bomWriter) Write(p []byte) (n int, err error) {
	if w.n == 0 {
		if n, err = w.w.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
			return
		}
		w.n += int64(n)
	}
	if n, err = w.w.Write(p); err != nil {
		return
	}
	w.n += int64(n)
	return
}
