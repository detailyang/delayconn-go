package delayconn

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type rwconn struct {
	r io.Reader
	w io.Writer
}

func (rwc *rwconn) Read(b []byte) (n int, err error) { return rwc.r.Read(b) }

// Write writes data to the connection.
// Write can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (rwc *rwconn) Write(b []byte) (n int, err error) { return rwc.w.Write(b) }

func (rwc *rwconn) Close() error { return nil }

func (rwc *rwconn) LocalAddr() net.Addr  { return &net.TCPAddr{} }
func (rwc *rwconn) RemoteAddr() net.Addr { return &net.TCPAddr{} }

func (rwc *rwconn) SetDeadline(t time.Time) error {
	return nil
}

func (rwc *rwconn) SetReadDeadline(t time.Time) error {
	return nil
}

func (rwc rwconn) SetWriteDeadline(t time.Time) error { return nil }

func TestReadDelayConn(t *testing.T) {
	r := bytes.NewBuffer(nil)
	w := bytes.NewBuffer(nil)
	rwc := &rwconn{r, w}

	started := time.Now()
	NewReadDelayConn(500*time.Millisecond, rwc).Read(nil)
	require.True(t, time.Since(started).Milliseconds() > 250)
}

func TestWriteDelayConn(t *testing.T) {
	r := bytes.NewBuffer(nil)
	w := bytes.NewBuffer(nil)
	rwc := &rwconn{r, w}

	started := time.Now()
	NewWriteDelayConn(500*time.Millisecond, rwc).Write(nil)
	require.True(t, time.Since(started).Milliseconds() > 250)
}

func TestOneByteWriteConn(t *testing.T) {
	r := bytes.NewBuffer(nil)
	w := bytes.NewBuffer(nil)
	rwc := &rwconn{r, w}

	var c [1024]byte
	n, err := NewOneByteWriteConn(rwc).Write(c[:])
	require.Nil(t, err)
	require.Equal(t, 1, n)
}

func TestOneByteReadConn(t *testing.T) {
	r := bytes.NewBuffer([]byte("haha"))
	w := bytes.NewBuffer(nil)
	rwc := &rwconn{r, w}

	var c [1024]byte
	n, err := NewOneByteReadConn(rwc).Read(c[:])
	require.Nil(t, err)
	require.Equal(t, 1, n)
}

func TestPerWriteDelayConn(t *testing.T) {
	r := bytes.NewBuffer([]byte("haha"))
	w := bytes.NewBuffer(nil)
	rwc := &rwconn{r, w}

	var c [1024]byte
	ww := NewPerWriteDelayConn(1*time.Millisecond, rwc)
	n, err := ww.Write(c[:])
	require.Nil(t, err)
	require.Equal(t, 1024, n)
}
