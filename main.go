package main

import (
	"CNB/internal/colorTheme"
	"CNB/internal/com"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"go.bug.st/serial"
	"log"
	"strconv"
)

func serializeIntSlice(input []int) []string {
	var valuesText []string
	for i := range input {
		number := input[i]
		text := strconv.Itoa(number)
		valuesText = append(valuesText, text)
	}
	return valuesText
}

func getParities(p map[string]serial.Parity) []string {
	keys := make([]string, 0, len(p))
	for s := range p {
		keys = append(keys, s)
	}
	return keys
}

func main() {
	a := app.New()
	ico, _ := fyne.LoadResourceFromPath(".\\resources\\icon.png")
	a.SetIcon(ico)
	a.Settings().SetTheme(colorTheme.GreyTheme())

	mainWindow := a.NewWindow("OCN LAB1")
	mainWindow.Resize(fyne.Size{Width: 400, Height: 250})

	var rxPortName string
	var txPortName string
	rxDropdown := widget.NewSelect(
		com.ScanPorts(),
		func(v string) {
			rxPortName = v
			fmt.Println("RX port changed:", v)
		},
	)
	txDropdown := widget.NewSelect(
		com.ScanPorts(),
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

	var parity serial.Parity
	paritySelect := widget.NewSelect(
		getParities(com.ParityMap),
		func(v string) {
			parity = com.ParityMap[v]
			log.Println("Parity:", v)
		},
	)
	paritySelect.SetSelectedIndex(0)

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
		com.SendStringToSerial(txPortName, serial.Mode{BaudRate: txSpeed, Parity: parity}, msg+"\n")
		rxMsg, err := com.ReceiveStringFromSerial(rxPortName, serial.Mode{BaudRate: rxSpeed, Parity: parity}, '\n')
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
		paritySelect,
		input,
		receivedText,
		buttonStart,
	)

	mainWindow.SetContent(l)
	mainWindow.ShowAndRun()
}
