package bus

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type Receiver struct {
	raw *azservicebus.Receiver
}

func (c *Client) NewQueueReceiver(queue string) (*Receiver, error) {
	r, err := c.raw.NewReceiverForQueue(queue, &azservicebus.ReceiverOptions{
		ReceiveMode: azservicebus.ReceiveModePeekLock,
	})
	if err != nil {
		return nil, err
	}
	return &Receiver{raw: r}, nil
}

func (r *Receiver) ReceiveAndComplete(
	ctx context.Context,
	handler func([]byte) error,
) error {

	msgs, err := r.raw.ReceiveMessages(ctx, 1, nil)
	if err != nil {
		return err
	}

	for _, m := range msgs {
		if err := handler(m.Body); err != nil {
			r.raw.AbandonMessage(ctx, m, nil)
			continue
		}
		r.raw.CompleteMessage(ctx, m, nil)
	}

	return nil
}
