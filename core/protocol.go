package core

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

/*

   Data Unit Format - https://datatracker.ietf.org/doc/html/rfc5734#section-4

   The EPP data unit contains two fields: a 32-bit header that describes
   the total length of the data unit, and the EPP XML instance.  The
   length of the EPP XML instance is determined by subtracting four
   octets from the total length of the data unit.  A receiver must
   successfully read that many octets to retrieve the complete EPP XML
   instance before processing the EPP message.

   EPP Data Unit Format (one tick mark represents one bit position):

       0                   1                   2                   3
       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                           Total Length                        |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                         EPP XML Instance                      |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+//-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

   Total Length (32 bits): The total length of the EPP data unit
   measured in octets in network (big endian) byte order.  The octets
   contained in this field MUST be included in the total length
   calculation.

   EPP XML Instance (variable length): The EPP XML instance carried in
   the data unit.
*/

const TotalLength int32 = 4

func writeBufferSize(size int, conn net.Conn) (n int, err error) {
	buffer := make([]byte, TotalLength)
	binary.BigEndian.PutUint32(buffer, uint32(size))
	sizeWritten, err := conn.Write(buffer)

	if err != nil {
		err := fmt.Errorf("error on write size buffer:  %s", err)
		return -1, err
	}

	return sizeWritten, nil
}

func readBufferSize(conn net.Conn) (size int, err error) {
	bufferSize := make([]byte, TotalLength)

	if _, err := conn.Read(bufferSize); err != nil {
		err := fmt.Errorf("error on read size buffer:  %s", err)
		return -1, err
	}

	return int(binary.BigEndian.Uint32(bufferSize)), nil
}

func writePendingBytes(messageSize int, messageSizeWritten int, messageBuffer []byte, conn net.Conn) (n int, err error) {
	var messageSizeCurrent = messageSizeWritten

	for {
		if messageSize != messageSizeCurrent {
			buffer := messageBuffer[messageSizeWritten:messageSize]

			currentWriteBytesSize, err := conn.Write(buffer)

			if err != nil {
				err := fmt.Errorf("error on write message buffer:  %s", err)
				return -1, err
			}

			messageSizeCurrent += currentWriteBytesSize

		} else {
			return messageSizeCurrent, nil
		}
	}

}

func readPendingBytes(messageSize int, messageSizeRead int, messageBuffer []byte, conn net.Conn) (buffer []byte, err error) {
	var messageSizeCurrent = messageSizeRead

	for {
		if messageSize != messageSizeCurrent {
			buffer := make([]byte, messageSize-messageSizeCurrent)

			currentReadBytesSize, err := conn.Read(buffer)

			messageSizeCurrent += currentReadBytesSize

			if err != nil {
				err := fmt.Errorf("error on read message buffer:  %s", err)
				return nil, err
			}

			messageBufferFull := append(messageBuffer[0:messageSizeRead], buffer...)

			return messageBufferFull, nil

		} else {
			return messageBuffer, nil
		}
	}
}

func WriteMessage(message string, conn net.Conn) (n int, err error) {
	messageBuffer := []byte(message)
	messageSize := len(messageBuffer)
	writeBufferSize(messageSize, conn)
	messageSizeWritten, err := conn.Write(messageBuffer)

	if err != nil {
		err := fmt.Errorf("error on write message buffer:  %s", err)
		return -1, err
	}

	if messageSize != messageSizeWritten {
		log.Println("not all bytes were written in first tentative, trying to write pending ")

		byteWritten, err := writePendingBytes(messageSize, messageSizeWritten, messageBuffer, conn)

		if err != nil {
			err := fmt.Errorf("error on write pending bytes to server:  %s", err)
			return byteWritten, err
		}

		return byteWritten, nil
	}

	return messageSizeWritten, nil
}

func ReadMessage(conn net.Conn) (message string, err error) {
	messageSize, _ := readBufferSize(conn)
	messageBuffer := make([]byte, messageSize)
	messageSizeRead, err := conn.Read(messageBuffer)

	if err != nil {
		err := fmt.Errorf("error on read message buffer:  %s", err)
		return "", err
	}

	if messageSize != messageSizeRead {
		log.Println("not all bytes were read in first tentative, trying to read pending")

		buffer, err := readPendingBytes(messageSize, messageSizeRead, messageBuffer, conn)

		if err != nil {
			err := fmt.Errorf("error on read pending bytes from server:  %s", err)
			return "", err
		}

		return string(buffer), nil
	}

	return string(messageBuffer), nil
}
