package test

import (
	"context"
	"fmt"
	"hr-api/pkg/bus"
	"hr-api/pkg/keyvault"
	"testing"
)

func TestBus(t *testing.T) {
	ctx := context.Background()
	client, err := bus.NewBusClient(keyvault.ServiceBusNamespace)
	if err != nil {
		fmt.Println("create bus client err: ", err.Error())
		return
	}
	sender, err := client.NewQueueSender(keyvault.ServiceBusQueueName)
	if err != nil {
		fmt.Println("NewQueueSender err: ", err.Error())
		return
	}

	s := "2222222222222222222222"
	err = sender.Send(ctx, []byte(s))
	if err != nil {
		fmt.Println("sender.send err: ", err.Error())
		return
	}
	//	receiver, err := client.NewQueueReceiver(keyvault.ServiceBusQueueName)
	//	if err != nil {
	//		fmt.Println("NewQueueReceiver err: ", err.Error())
	//		return
	//	}
	//	err = receiver.ReceiveAndComplete(ctx, func(b []byte) error {
	//		fmt.Println("收到了消息", string(b))
	//		return nil
	//	})
	//	if err != nil {
	//		fmt.Println("ReceiveAndComplete err: ", err.Error())
	//		return
	//	}
	//}
}
