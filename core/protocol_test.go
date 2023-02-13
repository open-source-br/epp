package core

import (
	"bytes"
	"testing"
)

// Stream is a fake for test purpose only
type Stream struct {
	buffer      bytes.Buffer
	limitReader int
	limitWriter int
}

func NewStream() *Stream {
	return &Stream{buffer: bytes.Buffer{}, limitReader: 20, limitWriter: 20}
}

func (stream *Stream) Read(b []byte) (n int, err error) {
	var buffer []byte
	currentBufferSize := len(b)

	// Will simulate limit size to read buffer. For test purpose only
	if currentBufferSize > stream.limitReader {
		buffer = make([]byte, stream.limitReader)

	} else {
		buffer = make([]byte, currentBufferSize)
	}

	sizes, err := stream.buffer.Read(buffer)

	copy(b, buffer)

	return sizes, err

}

func (stream *Stream) Write(b []byte) (n int, err error) {
	var buffer []byte
	currentMessageSize := len(b)

	// Will simulate limit size to write buffer. For test purpose only
	if currentMessageSize > stream.limitWriter {
		buffer = make([]byte, stream.limitWriter)

	} else {
		buffer = make([]byte, currentMessageSize)
	}

	copy(buffer, b)

	return stream.buffer.Write(buffer)
}

func TestWriteShortMessage(t *testing.T) {
	stream := NewStream()

	bytesWriter, _ := WriteMessage("Hello!", stream)

	if bytesWriter != 6 {
		t.Log("Expected 6 bytes written, got ', bytesWriter")
		t.FailNow()
	}
}

func TestWriteLargeMessage(t *testing.T) {
	stream := NewStream()

	bytesWriter, _ := WriteMessage("message to write buffer", stream)

	if bytesWriter != 23 {
		t.FailNow()
	}
}
func TestReadShortMessage(t *testing.T) {
	stream := NewStream()
	WriteMessage("Hello!", stream) // Write message to read

	message, _ := ReadMessage(stream)

	if message != "Hello!" {
		t.FailNow()
	}
}

func TestReadLargeMessage(t *testing.T) {
	stream := NewStream()
	WriteMessage("message to read buffer", stream) // Write message to read

	message, _ := ReadMessage(stream)

	if message != "message to read buffer" {
		t.FailNow()
	}
}
