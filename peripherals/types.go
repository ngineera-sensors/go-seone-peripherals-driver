package peripherals

type VCSELValue struct {
	Value byte
}

type MQTTResponse struct {
	Message interface{}
	Error   string
}
