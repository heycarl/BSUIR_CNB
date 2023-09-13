package comPacket

import (
	"bytes"
	"encoding/binary"
)

type ComPacket struct {
	header   byte
	checksum [4]byte
	payload  []byte
}

func CreatePacket(payload []byte, headerByte byte) ComPacket {
	var packet ComPacket
	packet.header = headerByte
	packet.payload = ByteStuffing(payload, headerByte)
	return packet
}

func ByteStuffing(data []byte, headerByte byte) []byte {
	output := bytes.Buffer{}
	for _, b := range data {
		if b == headerByte {
			output.WriteByte(headerByte)
			output.WriteByte(0x01)
		} else {
			output.WriteByte(b)
		}
	}
	return output.Bytes()
}

func getPayloadSize(packet []byte) int {
	payloadSizeBytes := packet[:3]
	payloadSize := int(binary.BigEndian.Uint32(payloadSizeBytes))
	return payloadSize
}

func DeByteStuffing(data []byte, headerByte byte) []byte {
	payload := bytes.Buffer{}
	i := 0
	for i < len(data) {
		if data[i] == headerByte {
			i++
			if i < len(data) && data[i] == 0x01 {
				payload.WriteByte(headerByte)
			} else {
				panic("Invalid byte stuffing")
			}
		} else {
			payload.WriteByte(data[i])
		}
		i++
	}
	return payload.Bytes()
}
