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

	s := "{\n  \"condition\": null,\n  \"conditionVersion\": null,\n  \"createdBy\": null,\n  \"createdOn\": \"2025-12-25T07:13:26.966520+00:00\",\n  \"delegatedManagedIdentityResourceId\": null,\n  \"description\": null,\n  \"id\": \"/subscriptions/2884693e-1b1f-4182-a931-38fce22157c4/resourceGroups/BingoHR/providers/Microsoft.ServiceBus/namespaces/hr-api-queue/providers/Microsoft.Authorization/roleAssignments/870845e3-0d0d-4762-8b08-34ee7de3efe6\",\n  \"name\": \"870845e3-0d0d-4762-8b08-34ee7de3efe6\",\n  \"principalId\": \"9219bfd9-35e8-40ce-aa87-63dedf3b3008\",\n  \"principalType\": \"User\",\n  \"resourceGroup\": \"BingoHR\",\n  \"roleDefinitionId\": \"/subscriptions/2884693e-1b1f-4182-a931-38fce22157c4/providers/Microsoft.Authorization/roleDefinitions/090c5cfd-751d-490a-894a-3ce6f1109419\",\n  \"scope\": \"/subscriptions/2884693e-1b1f-4182-a931-38fce22157c4/resourceGroups/BingoHR/providers/Microsoft.ServiceBus/namespaces/hr-api-queue\",\n  \"type\": \"Microsoft.Authorization/roleAssignments\",\n  \"updatedBy\": \"9219bfd9-35e8-40ce-aa87-63dedf3b3008\",\n  \"updatedOn\": \"2025-12-25T07:13:28.004561+00:00\"\n}"
	err = sender.Send(ctx, []byte(s))
	if err != nil {
		fmt.Println("sender.send err: ", err.Error())
		return
	}
	receiver, err := client.NewQueueReceiver(keyvault.ServiceBusQueueName)
	if err != nil {
		fmt.Println("NewQueueReceiver err: ", err.Error())
		return
	}
	err = receiver.ReceiveAndComplete(ctx, func(b []byte) error {
		fmt.Println("收到了消息", string(b))
		return nil
	})
	if err != nil {
		fmt.Println("ReceiveAndComplete err: ", err.Error())
		return
	}
}
