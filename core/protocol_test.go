package core

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"
)

// StreamConn is a fake to net.Conn for test purpose only
type StreamConn struct {
	buffer bytes.Buffer
}

func (stream *StreamConn) LocalAddr() net.Addr {
	return nil
}

func (stream *StreamConn) RemoteAddr() net.Addr {
	return nil
}

func (stream *StreamConn) Close() error {
	return fmt.Errorf("stream conn fake: not implemented, for test purpose only)")
}

func (stream *StreamConn) SetDeadline(t time.Time) error {
	return fmt.Errorf("stream conn fake: not implemented, for test purpose only)")
}

func (stream *StreamConn) SetReadDeadline(t time.Time) error {
	return fmt.Errorf("stream conn fake: not implemented, for test purpose only)")
}

func (stream *StreamConn) SetWriteDeadline(t time.Time) error {
	return fmt.Errorf("stream conn fake: not implemented, for test purpose only)")
}

func NewStreamConn() *StreamConn {
	return &StreamConn{buffer: bytes.Buffer{}}
}

func NewStreamConnWithData(defaultMessage string) *StreamConn {
	conn := &StreamConn{buffer: bytes.Buffer{}}
	WriteMessage(defaultMessage, conn)

	return conn
}

func (stream *StreamConn) Read(b []byte) (n int, err error) {
	var buffer []byte
	currentBufferSize := len(b)

	// Will simulate limit size to read buffer. For test purpose only
	if currentBufferSize > 20 {
		buffer = make([]byte, 20)

	} else {
		buffer = make([]byte, currentBufferSize)
	}

	sizes, err := stream.buffer.Read(buffer)

	copy(b, buffer)

	return sizes, err

}

func (stream *StreamConn) Write(b []byte) (n int, err error) {
	var buffer []byte
	currentMessageSize := len(b)

	// Will simulate limit size to write buffer. For test purpose only
	if currentMessageSize > 20 {
		buffer = make([]byte, 20)

	} else {
		buffer = make([]byte, currentMessageSize)
	}

	copy(buffer, b)

	return stream.buffer.Write(buffer)
}

func TestWriteShortMessage(t *testing.T) {
	stream := NewStreamConn()

	bytesWritter, _ := WriteMessage("hellou!", stream)

	if bytesWritter != 7 {
		t.FailNow()
	}
}

func TestWriteLargeMessage(t *testing.T) {
	stream := NewStreamConn()

	bytesWritter, _ := WriteMessage("message to write buffer", stream)

	if bytesWritter != 23 {
		t.FailNow()
	}
}

func TestReadShortMessage(t *testing.T) {
	stream := NewStreamConnWithData("Hello!")

	message, _ := ReadMessage(stream)

	if message != "Hello!" {
		t.FailNow()
	}
}

func TestReadLargeMessage(t *testing.T) {
	stream := NewStreamConnWithData("message to read buffer")

	message, _ := ReadMessage(stream)

	if message != "message to read buffer" {
		t.FailNow()
	}
}
