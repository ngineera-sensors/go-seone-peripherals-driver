package main

import (
	"fmt"
	"go-seone-peripherals-driver/peripherals"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"periph.io/x/host/v3"
)

func main() {
	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}

	vcselPin, err := peripherals.NewVCSELPin()
	if err != nil {
		log.Printf("Error occurred while initializing the VCSEL pin: %s. Exiting.", err)
		return
	}
	defer vcselPin.Halt()

	var halt = make(chan os.Signal)
	signal.Notify(halt, syscall.SIGTERM)
	signal.Notify(halt, syscall.SIGINT)

	go func() {
		select {
		case <-halt:
			if err := vcselPin.Halt(); err != nil {
				log.Println(err)
			}
			os.Exit(1)
		}
	}()

	mu := sync.Mutex{}

	client, err := peripherals.NewMQTTClient()
	if err != nil {
		log.Printf("Error occurred while initializing MQTT client: %s. Exiting.", err)
		return
	}
	getVcselSubscriptionTopic := fmt.Sprintf("seone/%s%s", peripherals.SEONE_SN, peripherals.VCSEL_GET_MQTT_TOPIC_PATH)
	log.Printf("Subscribing to VCSEL Getter: %s", getVcselSubscriptionTopic)
	client.Subscribe(getVcselSubscriptionTopic, 2, peripherals.GetVCSELHandler(vcselPin, &mu))

	setVcselSubscriptionTopic := fmt.Sprintf("seone/%s%s", peripherals.SEONE_SN, peripherals.VCSEL_SET_MQTT_TOPIC_PATH)
	log.Printf("Subscribing to VCSEL Setter: %s", setVcselSubscriptionTopic)
	client.Subscribe(setVcselSubscriptionTopic, 2, peripherals.SetVCSELHandler(vcselPin, &mu))

	for {
		time.Sleep(1 * time.Second)
	}
}
