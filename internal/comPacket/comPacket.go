package comPacket

const PacketHeader = 'z' + 4
const byteStuffingEsc = 'z' + 33

type Packet struct {
	Header             byte
	DestinationAddress byte
	SourceAddress      byte
	Data               []byte
	Fcs                byte
}

func CreatePacket(payload []byte) []byte {
	packet := Packet{
		Header:             PacketHeader,
		DestinationAddress: 0,
		SourceAddress:      0,
		Data:               payload,
		Fcs:                0,
	}
	return packet.serializePacket()
}

func (p Packet) serializePacket() []byte {
	var packetBytes []byte
	packetBytes = append(packetBytes, p.Header, p.DestinationAddress, p.SourceAddress)
	packetBytes = append(packetBytes, byteStuffing(p.Data)...)
	packetBytes = append(packetBytes, p.Fcs)
	return packetBytes
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
