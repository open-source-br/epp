package core

import (
	"bytes"
	"testing"
)

func TestWriteShortMessage(t *testing.T) {
	stream := &MockStream{buffer: bytes.Buffer{}, limitReader: 20, limitWriter: 20}

	bytesWriter, _ := WriteMessage("Hello!", stream)

	if bytesWriter != 6 {
		t.Log("Expected 6 bytes written, got ', bytesWriter")
		t.FailNow()
	}
}

func TestWriteLargeMessage(t *testing.T) {
	stream := &MockStream{buffer: bytes.Buffer{}, limitReader: 20, limitWriter: 20}

	bytesWriter, _ := WriteMessage("message to write buffer", stream)

	if bytesWriter != 23 {
		t.FailNow()
	}
}
func TestReadShortMessage(t *testing.T) {
	stream := &MockStream{buffer: bytes.Buffer{}, limitReader: 20, limitWriter: 20}
	WriteMessage("Hello!", stream) // echo

	message, _ := ReadMessage(stream)

	if message != "Hello!" {
		t.FailNow()
	}
}

func TestReadLargeMessage(t *testing.T) {
	stream := &MockStream{buffer: bytes.Buffer{}, limitReader: 20, limitWriter: 20}
	WriteMessage("message to read buffer", stream) // echo

	message, _ := ReadMessage(stream)

	if message != "message to read buffer" {
		t.FailNow()
	}
}

func TestReadFromEmptyStream(t *testing.T) {
	stream := &MockStream{buffer: bytes.Buffer{}, limitReader: 20, limitWriter: 20}

	_, err := ReadMessage(stream)

	if err == nil {
		t.Fatalf("Expected error when reading from empty stream, got nil")
	}
}

func TestReadShortMessageWithShortLimitReader(t *testing.T) {
	stream := &MockStream{buffer: bytes.Buffer{}, limitReader: 3, limitWriter: 20}
	WriteMessage("Hello!", stream) // echo

	message, _ := ReadMessage(stream)

	if message != "Hello!" {
		t.FailNow()
	}
}
