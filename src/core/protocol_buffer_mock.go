package core

import (
	"bytes"
)

// Stream is a mock for test purpose only
type MockStream struct {
	buffer      bytes.Buffer
	limitReader int
	limitWriter int
}

func (stream *MockStream) Read(b []byte) (n int, err error) {
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

func (stream *MockStream) Write(b []byte) (n int, err error) {
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
