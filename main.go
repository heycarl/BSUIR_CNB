package main

import (
	"CNB/internal/colorTheme"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"go.bug.st/serial"
	"log"
	"strconv"
	"strings"
)

func scanPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	return ports
}

func serializeIntSlice(input []int) []string {
	var valuesText []string
	for i := range input {
		number := input[i]
		text := strconv.Itoa(number)
		valuesText = append(valuesText, text)
	}
	return valuesText
}

func sendStringToSerial(portName string, mode serial.Mode, msg string) {
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

func receiveStringFromSerial(portName string, mode serial.Mode, term rune) (string, error) {
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

func main() {
	a := app.New()
	ico, _ := fyne.LoadResourceFromPath(".\\resources\\icon.png")
	a.SetIcon(ico)
	a.Settings().SetTheme(colorTheme.GreyTheme())

	window := a.NewWindow("OCN LAB1")
	window.Resize(fyne.Size{Width: 400, Height: 250})

	var rxPortName string
	var txPortName string
	rxDropdown := widget.NewSelect(
		scanPorts(),
		func(v string) {
			rxPortName = v
			fmt.Println("RX port changed:", v)
		},
	)
	txDropdown := widget.NewSelect(
		scanPorts(),
		func(v string) {
			txPortName = v
			fmt.Println("TX port changed:", v)
		},
	)

	availableSpeeds := []int{9600, 115200, 50}
	var rxSpeed = availableSpeeds[0]
	var txSpeed = availableSpeeds[0]
	rxSpeedDropdown := widget.NewSelect(
		serializeIntSlice(availableSpeeds),
		func(v string) {
			rxSpeed, _ = strconv.Atoi(v)
			log.Println("RX Speed changed:", v)
		},
	)
	txSpeedDropdown := widget.NewSelect(
		serializeIntSlice(availableSpeeds),
		func(v string) {
			txSpeed, _ = strconv.Atoi(v)
			log.Println("TX Speed changed:", v)
		},
	)
	txDropdown.SetSelectedIndex(0)
	rxDropdown.SetSelectedIndex(1)
	txSpeedDropdown.SetSelectedIndex(0)
	rxSpeedDropdown.SetSelectedIndex(0)

	input := widget.NewEntry()
	var inputMessage string
	inputMessageBinding := binding.BindString(&inputMessage)
	input.Bind(inputMessageBinding)
	input.SetPlaceHolder("Type text to transfer")

	receivedText := widget.NewLabel("Received text will be here")
	var receivedMessage = "Received text will be here"
	receivedMessageBinding := binding.BindString(&receivedMessage)
	receivedText.Bind(receivedMessageBinding)

	buttonStart := widget.NewButton("Send", func() {
		msg, _ := inputMessageBinding.Get()
		sendStringToSerial(txPortName, serial.Mode{BaudRate: txSpeed}, msg+"\n")
		rxMsg, err := receiveStringFromSerial(rxPortName, serial.Mode{BaudRate: rxSpeed}, '\n')
		if err != nil {
			log.Fatal(err)
			return
		}
		err = receivedMessageBinding.Set(rxMsg)
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	l := container.NewVBox(
		container.NewHBox(widget.NewLabel("RX port:"), rxDropdown),
		container.NewHBox(widget.NewLabel("TX port:"), txDropdown),
		container.NewHBox(widget.NewLabel("RX speed:"), rxSpeedDropdown),
		container.NewHBox(widget.NewLabel("TX speed:"), txSpeedDropdown),
		input,
		receivedText,
		buttonStart,
	)

	window.SetContent(l)
	window.ShowAndRun()
}
