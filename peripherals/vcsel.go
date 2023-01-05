package peripherals

import (
	"fmt"

	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

const (
	VCSEL_GPIO_BCM_PIN_NUMBER = "18"

	VCSEL_SET_MQTT_TOPIC_PATH    = "/periph/vcsel/set"
	VCSEL_SET_CB_MQTT_TOPIC_PATH = "/periph/vcsel/set/cb"

	VCSEL_GET_MQTT_TOPIC_PATH    = "/periph/vcsel/get"
	VCSEL_GET_CB_MQTT_TOPIC_PATH = "/periph/vcsel/get/cb"
)

var (
	VCSEL_OUTPUT_VALUE byte = 0
)

func GetVCSEL(pin gpio.PinIO) byte {
	return VCSEL_OUTPUT_VALUE
}

func SetVCSEL(pin gpio.PinIO, value byte) error {
	var level gpio.Level

	switch value {
	case 0:
		level = gpio.Low
		break
	case 1:
		level = gpio.High
		break
	default:
		err := fmt.Errorf("invalid value for vcsel setting")
		return err
	}
	err := pin.Out(level)

	VCSEL_OUTPUT_VALUE = value

	return err
}

func NewVCSELPin() (gpio.PinIO, error) {
	var err error
	var pin gpio.PinIO
	if _, err = driverreg.Init(); err != nil {
		return pin, err
	}

	pin = gpioreg.ByName(VCSEL_GPIO_BCM_PIN_NUMBER)
	if pin == nil {
		err = fmt.Errorf("failed to find pin %s", VCSEL_GPIO_BCM_PIN_NUMBER)
		return pin, err
	}
	return pin, err
}
