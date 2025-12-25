package bus

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type Sender struct {
	raw *azservicebus.Sender
}

func (c *Client) NewQueueSender(queue string) (*Sender, error) {
	s, err := c.raw.NewSender(queue, nil)
	if err != nil {
		return nil, err
	}
	return &Sender{raw: s}, nil
}

func (s *Sender) Send(ctx context.Context, body []byte) error {
	msg := &azservicebus.Message{Body: body}
	return s.raw.SendMessage(ctx, msg, nil)
}

func (s *Sender) SendScheduled(ctx context.Context, body []byte, at time.Time) error {
	msg := &azservicebus.Message{
		Body:                 body,
		ScheduledEnqueueTime: &at,
	}
	return s.raw.SendMessage(ctx, msg, nil)
}
