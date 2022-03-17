package firebase

// refs
// https://github.com/firebase/firebase-admin-go/blob/e60757f9b29711f19fa1f44ce9b5a6fae3baf3a5/snippets/messaging.go

import (
	"context"
	"fmt"
	"sync"

	fb "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/ExchangeDiary/exchange-diary/domain/vo"
	"github.com/ExchangeDiary/exchange-diary/infrastructure"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"google.golang.org/api/option"
)

const credentials = "firebase-credential.json"

var (
	firebaseClient *Client
	clientOnce     sync.Once
)

// Client ...
type Client struct {
	client *messaging.Client
	ctx    context.Context
}

func init() {
	logger.Info("init firebase alarm client")
	clientOnce.Do(func() {
		var app *fb.App
		var err error
		ctx := context.Background()

		switch infrastructure.Getenv("PHASE", "dev") {
		case "prod":
			app, err = fb.NewApp(context.Background(), nil)
		default:
			app, err = fb.NewApp(ctx, nil, option.WithCredentialsFile(credentials))
		}

		if err != nil {
			panic("Failed to load firebase cloud  " + err.Error())
		}

		msgClient, err := app.Messaging(ctx)
		if err != nil {
			panic("Failed to load firebase cloud messaging client  " + err.Error())
		}

		firebaseClient = &Client{
			client: msgClient,
			ctx:    ctx,
		}
	})
}

// GetClient ...
func GetClient() *Client {
	return firebaseClient
}

// Push ...
func (c *Client) Push(deviceTokens []string, messageBody *vo.AlarmBody) (failedTokens []string, err error) {
	var batchResponse *messaging.BatchResponse

	message := &messaging.MulticastMessage{
		Data:   messageBody.ConvertToMap(),
		Tokens: deviceTokens,
	}

	batchResponse, err = c.client.SendMulticast(c.ctx, message)
	if err != nil {
		return
	}

	if batchResponse.FailureCount > 0 {
		for idx, resp := range batchResponse.Responses {
			if !resp.Success {
				failedTokens = append(failedTokens, deviceTokens[idx])
			}
		}
		logger.Info(fmt.Sprintf("List of tokens that caused failures: %v\n", failedTokens))
		return
	}
	return nil, nil
}
