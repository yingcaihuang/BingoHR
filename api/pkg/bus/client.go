package bus

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"strings"
)

type Client struct {
	raw *azservicebus.Client
}

func NewFromConnStr(connStr string) (*Client, error) {
	c, err := azservicebus.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		return nil, err
	}
	return &Client{raw: c}, nil
}

func NewBusClient(namespace string) (*Client, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	// ✅ 自动补全 namespace
	fqdn := namespace
	if !strings.Contains(namespace, ".") {
		fqdn = namespace + ".servicebus.windows.net"
	}

	c, err := azservicebus.NewClient(fqdn, cred, &azservicebus.ClientOptions{
		// ✅ 防公司网络 / 5671 被封
		//EnableWebSockets: true,
	})
	if err != nil {
		return nil, err
	}

	return &Client{raw: c}, nil
}
