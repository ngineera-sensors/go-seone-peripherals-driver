package peripherals

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"periph.io/x/conn/v3/gpio"
)

var (
	SEONE_SN = ""
)

func init() {
	sn, err := os.ReadFile(filepath.Join("config", "serialnumber.txt"))
	if err != nil {
		log.Fatal(err)
	}
	if len(sn) != 0 {
		log.Printf("Setting SEONE_SN value: %s", string(sn))
	}
	snStr := string(sn)
	snStr = strings.TrimSpace(snStr)
	SEONE_SN = snStr
}

func GetVCSELHandler(pin gpio.PinIO, mu *sync.Mutex) mqtt.MessageHandler {

	var f = func(client mqtt.Client, msg mqtt.Message) {

		mu.Lock()
		defer mu.Unlock()

		v := GetVCSEL(pin)
		log.Printf("Getting VCSEL value: %d", v)

		respTopic := fmt.Sprintf("seone/%s%s", SEONE_SN, VCSEL_GET_CB_MQTT_TOPIC_PATH)
		respObj := MQTTResponse{
			Message: VCSELValue{
				Value: v,
			},
		}
		err := PublishJsonMsg(respTopic, respObj, client)
		if err != nil {
			log.Printf("Error occurred in GetVCSELHandler MQTT CB: %s", err.Error())
		}
	}
	return f
}

func SetVCSELHandler(pin gpio.PinIO, mu *sync.Mutex) mqtt.MessageHandler {
	var f = func(client mqtt.Client, msg mqtt.Message) {
		var err error

		payload := msg.Payload()
		var value VCSELValue
		err = json.Unmarshal(payload, &value)
		if err != nil {
			log.Printf("Error occurred in SetVCSELHandler MQTT CB while unmarshalling the JSON message: %s", err.Error())
		}

		log.Printf("Setting VCSEL value to %d", value.Value)

		mu.Lock()
		defer mu.Unlock()

		respTopic := fmt.Sprintf("seone/%s%s", SEONE_SN, VCSEL_SET_CB_MQTT_TOPIC_PATH)

		err = SetVCSEL(pin, value.Value)
		if err != nil {
			respObj := MQTTResponse{
				Error: err.Error(),
			}
			err = PublishJsonMsg(respTopic, respObj, client)
			if err != nil {
				log.Printf("Error occurred in GetVCSELHandler MQTT CB: %s", err.Error())
			}
			return
		}

		respObj := MQTTResponse{
			Message: VCSELValue{
				Value: GetVCSEL(pin),
			},
		}
		err = PublishJsonMsg(respTopic, respObj, client)
		if err != nil {
			log.Printf("Error occurred in GetVCSELHandler MQTT CB: %s", err.Error())
		}
	}
	return f
}

func MainLoop() {

}
