package delayconn

import "io"

type oneByteWriter struct {
	w io.Writer
}

// OneByteWriter returns a writer
func OneByteWriter(w io.Writer) io.Writer {
	return &oneByteWriter{w: w}
}

func (w *oneByteWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return w.w.Write(p[0:1])
}
