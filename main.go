package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3"
)

func ExamplePinOut_pWM() {
	// Make sure periph is initialized.
	// TODO: Use host.Init(). It is not used in this example to prevent circular
	// go package import.
	if _, err := driverreg.Init(); err != nil {
		log.Fatal(err)
	}

	// Use gpioreg GPIO pin registry to find a GPIO pin by name.
	pinName := "13"
	p := gpioreg.ByName(pinName)
	if p == nil {
		log.Fatalf("Failed to find pin %s", pinName)
	}

	defer p.Halt()

	err := p.Out(gpio.High)
	if err != nil {
		log.Fatal(err)
	}

	var halt = make(chan os.Signal)
	signal.Notify(halt, syscall.SIGTERM)
	signal.Notify(halt, syscall.SIGINT)

	go func() {
		select {
		case <-halt:
			if err := p.Halt(); err != nil {
				log.Println(err)
			}
			os.Exit(1)
		}
	}()

	t := time.NewTicker(1 * time.Second)

	// for l := gpio.Low; ; l = !l {
	// 	log.Println(l)
	// 	err := p.Out(l)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	<-t.C
	// }

	for l := 0; l < 100; l += 10 {
		duty, err := gpio.ParseDuty(fmt.Sprintf("%d%%", l))
		if err != nil {
			log.Fatal(err)
		}
		log.Println(l, duty)
		if err := p.PWM(
			duty,
			50*physic.KiloHertz,
		); err != nil {
			log.Fatal(err)
		}
		<-t.C
	}
}

func main() {
	state, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Host init status: ", state.Loaded)

	ExamplePinOut_pWM()
}
