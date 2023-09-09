package main

import (
	"CNB/internal/colorTheme"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"go.bug.st/serial"
	"log"
	"strconv"
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

func main() {
	a := app.New()
	ico, _ := fyne.LoadResourceFromPath(".\\resources\\icon.png")
	a.SetIcon(ico)
	a.Settings().SetTheme(colorTheme.GreyTheme())

	window := a.NewWindow("OCN LAB1")
	window.Resize(fyne.Size{Width: 400, Height: 250})

	rxDropdown := widget.NewSelect(
		scanPorts(),
		func(value string) {
			fmt.Println("RX port changed:", value)
		},
	)
	txDropdown := widget.NewSelect(
		scanPorts(),
		func(value string) {
			fmt.Println("TX port changed:", value)
		},
	)

	availableSpeeds := []int{9600, 115200, 800}
	rxSpeedDropdown := widget.NewSelect(
		serializeIntSlice(availableSpeeds),
		func(value string) {
			fmt.Println("RX Speed changed:", value)
		},
	)
	txSpeedDropdown := widget.NewSelect(
		serializeIntSlice(availableSpeeds),
		func(value string) {
			fmt.Println("TX Speed changed:", value)
		},
	)
	txDropdown.SetSelectedIndex(0)
	rxDropdown.SetSelectedIndex(1)
	txSpeedDropdown.SetSelectedIndex(0)
	rxSpeedDropdown.SetSelectedIndex(0)

	input := widget.NewEntry()
	input.SetPlaceHolder("Type text to transfer")

	text := widget.NewLabel("Received text will be here")

	buttonStart := widget.NewButton("Start", func() {
		fmt.Println("Start button clicked")
	})

	l := container.NewVBox(
		container.NewHBox(widget.NewLabel("RX port:"), rxDropdown),
		container.NewHBox(widget.NewLabel("TX port:"), txDropdown),
		container.NewHBox(widget.NewLabel("RX speed:"), rxSpeedDropdown),
		container.NewHBox(widget.NewLabel("TX speed:"), txSpeedDropdown),
		input,
		text,
		buttonStart,
	)

	window.SetContent(l)
	window.ShowAndRun()
}

//package main
//
//import (
//	"go.bug.st/serial"
//	"log"
//	"strings"
//)
//
//func main() {
//	ports, err := serial.GetPortsList()
//	if err != nil {
//		log.Fatal(err)
//	}
//	if len(ports) == 0 {
//		log.Fatal("No serial ports found!")
//	}
//	for _, port := range ports {
//		log.Printf("Found port: %v\n", port)
//	}
//
//	rxMode := &serial.Mode{
//		BaudRate: 115200,
//	}
//	rxPort, err := serial.Open("COM3", rxMode)
//	txMode := &serial.Mode{
//		BaudRate: 9600,
//	}
//	txPort, err := serial.Open("COM4", txMode)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	n, err := txPort.Write([]byte("Hello, world!\n"))
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	log.Printf("Wrote %d bytes\n", n)
//
//	log.Println("Received message: ")
//	for {
//		buff := make([]byte, 256)
//		n, err := rxPort.Read(buff)
//		if err != nil {
//			log.Fatal(err)
//		}
//		if n == 0 {
//			log.Println("\nEOF")
//			break
//		}
//
//		log.Printf("%s", string(buff[:n]))
//
//		if strings.Contains(string(buff[:n]), "\n") {
//			break
//		}
//	}
//
//}
