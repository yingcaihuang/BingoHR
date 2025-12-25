package bus

import (
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func (c *Client) NewSessionReceiver(queue string) (*azservicebus.Receiver, error) {
	return c.raw.NewReceiverForQueue(queue, nil)
}
