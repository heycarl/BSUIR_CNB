package comPacket

import "errors"

const PacketHeader = 'z' + 4
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
		Fcs:                0,
	}
	return packet
}

func (p Packet) SerializePacket() []byte {
	var packetBytes []byte
	packetBytes = append(packetBytes, p.Header, p.DestinationAddress, p.SourceAddress)
	packetBytes = append(packetBytes, byteStuffing(p.Data)...)
	packetBytes = append(packetBytes, p.Fcs)
	return packetBytes
}

func DeserializePacket(raw []byte) (Packet, error) {
	if raw[0] != PacketHeader {
		return Packet{}, errors.New("incorrect header byte")
	}
	packet := Packet{
		Header:             raw[0],
		DestinationAddress: raw[1],
		SourceAddress:      raw[2],
		Data:               deByteStuffing(raw[2 : len(raw)-1]),
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
	return deByteStuffed
}
