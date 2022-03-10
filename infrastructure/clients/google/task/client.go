package task

import (
	"context"
	"fmt"
	"sync"

	"github.com/ExchangeDiary/exchange-diary/infrastructure"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"google.golang.org/api/option"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// Kind represents the name of the location/storage type.
const Kind = "google"

var (
	credentials       = "credentials.json"
	vodaStorageClient *Client
	clientOnce        sync.Once
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

		vodaStorageClient = &Client{
			client: client,
			ctx:    ctx,
			queuePath: fmt.Sprintf("projects/%s/locations/%s/queues/%s",
				infrastructure.Getenv("PROJECT_ID", "voda-342511"),
				infrastructure.Getenv("LOCATION_ID", "asia-northeast3"),
				infrastructure.Getenv("QUEUE_ID", "voda-alarm-queue")),
		}

		// resp := mockTask()
		// logger.Info(resp.String())
	})
}

// GetClient ...
func GetClient() *Client {
	return vodaStorageClient
}

// Close ...
func (tc *Client) Close() {
	tc.client.Close()
}

// CreateTask ...
// https://pkg.go.dev/cloud.google.com/go/cloudtasks/apiv2#CallOptions
// https://github.com/GoogleCloudPlatform/golang-samples/blob/c5b5b4be9bb51fc05a8939b163374bc23084eb56/tasks/create_http_task.go
// https://ichi.pro/ko/gcp-cloud-tasksleul-sayonghaneun-ibenteu-giban-yeyag-jag-eob-254067840428949
// https://tkdguq05.github.io/2020/05/19/google-task/
// https://github.com/ArticsIS/Google-Cloud-Helpers/blob/master/services/taskqueue.py
// https://cloud.google.com/tasks/docs/tutorial-gcf
func (tc *Client) CreateTask(url, message string, httpMethod taskspb.HttpMethod) (*taskspb.Task, error) {
	req := &taskspb.CreateTaskRequest{
		Parent: tc.queuePath,
		Task:   buildTask(url, httpMethod),
	}
	req.Task.GetHttpRequest().Body = []byte(message)
	createdTask, err := tc.client.CreateTask(tc.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}
	return createdTask, nil
}

func mockTask() *taskspb.Task {
	response, err := vodaStorageClient.CreateTask("http://api.duckduckgo.com/?q=minkj1992&format=json", "", taskspb.HttpMethod_GET)
	if err != nil {
		logger.Error(err.Error())
	}
	return response
}

func buildTask(url string, httpMethod taskspb.HttpMethod) *taskspb.Task {
	return &taskspb.Task{
		MessageType: &taskspb.Task_HttpRequest{
			HttpRequest: &taskspb.HttpRequest{
				HttpMethod: httpMethod,
				Url:        url,
			},
		},
	}
}

// UpdateVTask ...
func (tc *Client) UpdateVTask(id string) (*VTask, error) {
	return &VTask{}, nil
}

// DeleteVTask ...
func (tc *Client) DeleteVTask(id string) error {
	return nil
}
