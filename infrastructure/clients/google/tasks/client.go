package tasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ExchangeDiary/exchange-diary/infrastructure"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"google.golang.org/api/option"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Kind represents the name of the location/storage type.
const Kind = "google"

var nilTime = time.Time{}

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

// TaskID returns google cloud task unique id
// projects/PROJECT_ID/locations/LOCATION_ID/queues/QUEUE_ID/tasks/TASK_ID
func (c *Client) TaskID(id string) string {
	return fmt.Sprintf("%s/tasks/%s", c.queuePath, id)
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
	})
}

// GetClient ...
func GetClient() *Client {
	return vodaStorageClient
}

// Close ...
func (c *Client) Close() {
	c.client.Close()
}

// BuildTask ...
func (c *Client) BuildTask(url string, body []byte, httpMethod taskspb.HttpMethod, scheduledAt time.Time) *taskspb.Task {
	task := &taskspb.Task{
		MessageType: &taskspb.Task_HttpRequest{
			HttpRequest: &taskspb.HttpRequest{
				HttpMethod: httpMethod,
				Url:        url,
				Body:       body,
			},
		},
	}

	// if ScheduleTime set nil, google cloud task run this task right away
	if scheduledAt != nilTime {
		task.ScheduleTime = timestamppb.New(scheduledAt)
	}

	return task
}

// RegisterTask register a task and adds it to a queue.
// https://pkg.go.dev/cloud.google.com/go/cloudtasks/apiv2#CallOptions
// https://github.com/GoogleCloudPlatform/golang-samples/blob/c5b5b4be9bb51fc05a8939b163374bc23084eb56/tasks/create_http_task.go
// https://ichi.pro/ko/gcp-cloud-tasksleul-sayonghaneun-ibenteu-giban-yeyag-jag-eob-254067840428949
// https://tkdguq05.github.io/2020/05/19/google-task/
// https://github.com/ArticsIS/Google-Cloud-Helpers/blob/master/services/taskqueue.py
// https://cloud.google.com/tasks/docs/tutorial-gcf
func (c *Client) RegisterTask(task *taskspb.Task) (*taskspb.Task, error) {
	req := &taskspb.CreateTaskRequest{
		Parent: c.queuePath,
		Task:   task,
	}

	registeredTask, err := c.client.CreateTask(c.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}
	return registeredTask, nil
}

// UpdateVTask ...
func (c *Client) UpdateVTask(id string) (*VTask, error) {
	return &VTask{}, nil
}

// DeleteVTask ...
func (c *Client) DeleteVTask(id string) error {
	return nil
}
