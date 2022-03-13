package firebase

import (
	"context"
	"sync"

	fb "firebase.google.com/go/v4"
	"github.com/ExchangeDiary/exchange-diary/infrastructure"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"google.golang.org/api/option"
)

const credentials = "credentials.json"

var (
	firebaseClient *Client
	clientOnce     sync.Once
)

// Client ...
type Client struct {
	app *fb.App
	ctx context.Context
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

		firebaseClient = &Client{
			app: app,
			ctx: ctx,
		}
	})
}

// GetClient ...
func GetClient() *Client {
	return firebaseClient
}

// https://github.com/firebase/firebase-admin-go/blob/e60757f9b29711f19fa1f44ce9b5a6fae3baf3a5/snippets/messaging.go
