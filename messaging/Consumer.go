package messaging

type Consumer interface {
	Receive() interface{}
}
