package network_handler

import (
	"fmt"
	"global"

	"github.com/pebbe/zmq4"
)

var err error

const (
	filter = "CUBE"
)

func InitZeroMQ() error {
	global.Publisher, err = zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		return fmt.Errorf("Creating zmq4 publisher failed: %v", err)
	}
	global.Subscriber, err = zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		return fmt.Errorf("Creating zmq4 subscriber failed: %v", err)
	}
	err = global.Publisher.Bind(fmt.Sprintf("tcp://*:%v", global.Conf.PublishPort))
	if err != nil {
		return fmt.Errorf("Binding zmq4 publisher to tcp://*:%v failed: %v", global.Conf.PublishPort, err)
	}
	err = global.Subscriber.Connect(fmt.Sprintf("tcp://localhost:%v", global.Conf.SubscribePort))
	if err != nil {
		return fmt.Errorf("Connecting zmq4 subscriber to tcp://localhost:%v failed: %v", global.Conf.SubscribePort, err)
	}
	err = global.Subscriber.SetSubscribe(filter)
	if err != nil {
		return fmt.Errorf("Setting subscribtion (filter: %v) failed: %v", filter, err)
	}
	return nil
}
