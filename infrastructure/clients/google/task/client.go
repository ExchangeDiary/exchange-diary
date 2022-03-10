package task

import (
	"context"
	"fmt"
	"sync"

	"github.com/ExchangeDiary/exchange-diary/infrastructure"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"google.golang.org/api/option"
)

// Kind represents the name of the location/storage type.
const Kind = "google"

var (
	credentials = "credentials.json"
	vodaClient  *Client
	clientOnce  sync.Once
)

// Client for google cloud tasks
type Client struct {
	client    *cloudtasks.Client
	ctx       context.Context
	queuePath string
}

// https://github.com/GoogleCloudPlatform/golang-samples/blob/c5b5b4be9bb51fc05a8939b163374bc23084eb56/tasks/create_http_task.go
func init() {
	logger.Info("lazy init google cloud task client")
	// LazyGlobal loading ...
	clientOnce.Do(func() {
		var client *cloudtasks.Client
		var err error
		ctx := context.Background()
		switch infrastructure.Getenv("PHASE", "dev") {
		case "prod":
			client, err = cloudtasks.NewClient(ctx)
		default:
			client, err = cloudtasks.NewClient(ctx, option.WithCredentialsFile(credentials))
		}

		if err != nil {
			panic("Failed to load google cloud tasks  " + err.Error())
		}

		vodaClient = &Client{
			client: client,
			ctx:    ctx,
			queuePath: fmt.Sprintf("projects/%s/locations/%s/queues/%s",
				infrastructure.Getenv("PROJECT_ID", "voda-342511"),
				infrastructure.Getenv("LOCATION_ID", "asia-northeast3"),
				infrastructure.Getenv("QUEUE_ID", "voda-alarm-queue")),
		}
	})

}

// GetClient ...
func GetClient() *Client {
	return vodaClient
}

// Close ...
func (tc *Client) Close() {
	tc.client.Close()
}

// VTask ...
func (tc *Client) VTask(id string) (*VTask, error) {
	// task, err := tc.client.GetTask()
	return &VTask{}, nil
}

// CreateVTask ...
func (tc *Client) CreateVTask(id string) (*VTask, error) {
	return &VTask{}, nil
}

// UpdateVTask ...
func (tc *Client) UpdateVTask(id string) (*VTask, error) {
	return &VTask{}, nil
}

// DeleteVTask ...
func (tc *Client) DeleteVTask(id string) error {
	return nil
}
