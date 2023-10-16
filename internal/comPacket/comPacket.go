package comPacket

import (
	"CNB/internal/hammingCode"
	"errors"
	"log"
)

const PacketHeader = 'z' + 0
const byteStuffingEsc = 'z' + 33

type Packet struct {
	Header             byte
	DestinationAddress byte
	SourceAddress      byte
	Data               []byte
	Fcs                byte
}

func CreatePacket(payload []byte) Packet {
	packet := Packet{
		Header:             PacketHeader,
		DestinationAddress: 0,
		SourceAddress:      0,
		Data:               payload,
		Fcs:                1,
	}
	return packet
}

func (p Packet) SerializePacket() []byte {
	var packetBytes []byte
	packetBytes = append(packetBytes, p.Header, p.DestinationAddress, p.SourceAddress)
	frameData := hammingCode.Encode(p.Data)
	frameData = byteStuffing(frameData)
	packetBytes = append(packetBytes, frameData...)
	packetBytes = append(packetBytes, p.Fcs)
	return packetBytes
}

func DeserializePacket(raw []byte) (Packet, error) {
	if raw[0] != PacketHeader {
		return Packet{}, errors.New("incorrect header byte")
	}
	payloadBytes := deByteStuffing(raw[3 : len(raw)-1])
	payloadBytes, _ = hammingCode.Decode(payloadBytes)
	packet := Packet{
		Header:             raw[0],
		DestinationAddress: raw[1],
		SourceAddress:      raw[2],
		Data:               payloadBytes,
		Fcs:                raw[len(raw)-1]}
	return packet, nil
}

func byteStuffing(data []byte) []byte {
	var byteStuffed []byte
	for _, b := range data {
		if b == PacketHeader || b == byteStuffingEsc {
			byteStuffed = append(byteStuffed, byteStuffingEsc)
		}
		byteStuffed = append(byteStuffed, b)
	}
	if len(data) != len(byteStuffed) {
		log.Println("Byte stuffing:\n", data, " -> ", byteStuffed)
	}
	return byteStuffed
}

func deByteStuffing(data []byte) []byte {
	var deByteStuffed []byte
	escaped := false // flag to track if the previous byte was the escape byte

	for _, b := range data {
		if escaped {
			if b == PacketHeader || b == byteStuffingEsc {
				deByteStuffed = append(deByteStuffed, b)
			}
			escaped = false
		} else {
			if b == byteStuffingEsc {
				escaped = true
			} else {
				deByteStuffed = append(deByteStuffed, b)
			}
		}
	}
	if len(data) != len(deByteStuffed) {
		log.Println("Byte de-stuffing:\n", data, " -> ", deByteStuffed)
	}
	return deByteStuffed
}
