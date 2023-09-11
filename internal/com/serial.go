package com

import (
	"go.bug.st/serial"
	"log"
	"strings"
)

var ParityMap = map[string]serial.Parity{
	"No Parity":    serial.NoParity,
	"Odd Parity":   serial.OddParity,
	"Even Parity":  serial.EvenParity,
	"Mark Parity":  serial.MarkParity,
	"Space Parity": serial.SpaceParity,
}

func SendStringToSerial(portName string, mode serial.Mode, msg string) {
	port, err := serial.Open(portName, &mode)
	if err != nil {
		log.Fatal(err)
		return
	}
	n, err := port.Write([]byte(msg))
	log.Printf("Written %d bytes\n", n)
	err = port.Close()
	if err != nil {
		return
	}
}

func ReceiveStringFromSerial(portName string, mode serial.Mode, term rune) (string, error) {
	port, rxErr := serial.Open(portName, &mode)
	if rxErr != nil {
		log.Printf("Port %s open failed\n", portName)
		return "", rxErr
	}

	buff := make([]byte, 256)
	var rxBytes = 0
	var err error
	for {
		rxBytes, err = port.Read(buff)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		if rxBytes == 0 {
			log.Println("\nEOF")
			break
		}
		if strings.Contains(string(buff[:rxBytes]), string(term)) {
			break
		}
	}
	err = port.Close()
	if err != nil {
		return "", err
	}
	log.Printf("Recieved %d bytes\n", rxBytes)
	return string(buff), nil
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
