package com

import (
	"CNB/internal/comPacket"
	"go.bug.st/serial"
	"log"
)

var ParityMap = map[string]serial.Parity{
	"No Parity":    serial.NoParity,
	"Odd Parity":   serial.OddParity,
	"Even Parity":  serial.EvenParity,
	"Mark Parity":  serial.MarkParity,
	"Space Parity": serial.SpaceParity,
}

type Port struct {
	Name   string
	Speed  int
	Parity serial.Parity
}

func (p Port) WriteBytes(payload []byte) {
	port, err := serial.Open(p.Name, &serial.Mode{BaudRate: p.Speed, Parity: p.Parity})
	if err != nil {
		log.Fatal(err)
		return
	}
	n, err := port.Write(payload)
	log.Printf("Written %d bytes\n", n)
	err = port.Close()
	if err != nil {
		return
	}
}

func (p Port) ReceiveBytes() ([]byte, error) {
	port, err := serial.Open(p.Name, &serial.Mode{BaudRate: p.Speed, Parity: p.Parity})
	if err != nil {
		log.Printf("Port %s open failed\n", p.Name)
		return nil, err
	}

	buff := make([]byte, 256)
	var n, readErr = port.Read(buff)
	if readErr != nil {
		return nil, err
	}
	err = port.Close()

	log.Printf("Recieved %d bytes\n", n)
	return buff[:n], nil
}

func (p Port) WritePacket(packet comPacket.Packet) {
	p.WriteBytes(packet.SerializePacket())
}

func (p Port) ReadPacket() (comPacket.Packet, error) {
	rawPacketData, err := p.ReceiveBytes()
	if err != nil {
		log.Println(err)
		return comPacket.Packet{}, err
	}
	packet, err := comPacket.DeserializePacket(rawPacketData)
	if err != nil {
		return comPacket.Packet{}, err
	}
	return packet, nil
}

func ScanPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No com ports found!")
	}
	return ports
}
