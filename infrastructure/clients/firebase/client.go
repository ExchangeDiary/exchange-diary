package firebase

// refs
// https://github.com/firebase/firebase-admin-go/blob/e60757f9b29711f19fa1f44ce9b5a6fae3baf3a5/snippets/messaging.go

import (
	"context"
	"fmt"
	"os"
	"sync"

	fb "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/ExchangeDiary/exchange-diary/domain/vo"
	"github.com/ExchangeDiary/exchange-diary/infrastructure"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"google.golang.org/api/option"
)

const (
	credentials = "firebase-credential.json"
	credKey     = "FIREBASE_CREDENTIALS_SECRET"
)

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
			// for local test cmd
			// export FIREBASE_CREDENTIALS_SECRET=`cat ./firebase-credential.json`
			app, err = fb.NewApp(context.Background(), nil, option.WithCredentialsJSON(getCredntial()))
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
func (c *Client) Push(deviceTokens []string, messageBody *vo.Alarm) (failedTokens []string, err error) {
	var batchResponse *messaging.BatchResponse
	messagePayload := messageBody.ConvertToMap()
	message := &messaging.MulticastMessage{
		Data:   messagePayload,
		Tokens: deviceTokens,
		Notification: &messaging.Notification{
			Title: messagePayload["roomName"] + " 다이어리방",
			Body:  messagePayload["title"],
		},
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

// refs: https://cloud.google.com/run/docs/configuring/secrets
// refs: https://cloud.google.com/run/docs/tutorials/identity-platform#secret-manager
func getCredntial() []byte {
	secret := os.Getenv(credKey)
	if len(secret) != 0 {
		return []byte(secret)
	}
	return []byte{}
}
