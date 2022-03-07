package cloudstorage

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Kind represents the name of the location/storage type.
const Kind = "google"

var (
	// ErrNotFound ...
	ErrNotFound = errors.New("google storage not found")
)

var (
	credentials       = "credentials.json"
	vodaStorageClient *VClient
)

// VClient ...
type VClient struct {
	client *storage.Client
	ctx    context.Context
}

func init() {
	ctx, client, err := newGoogleStorageClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	vodaStorageClient = &VClient{client: client, ctx: ctx}
}

// GetClient returns vodaStorageClient refs
func GetClient() *VClient {
	return vodaStorageClient
}

// Close ...
func (vc *VClient) Close() {
	vc.client.Close()
}

// VBucket retrieves a google cloud bucket based on its name which must be exact.
func (vc *VClient) VBucket(id string) (*VBucket, error) {
	attrs, err := vc.client.Bucket(id).Attrs(vc.ctx)
	if err != nil {
		if err == storage.ErrBucketNotExist {
			return nil, ErrNotFound
		}
		return nil, err
	}

	b := &VBucket{
		name:   attrs.Name,
		client: vc.client,
		ctx:    vc.ctx,
	}
	return b, nil
}

// ItemByURL get url and retuns real item itself (VItem)
func (vc *VClient) ItemByURL(url *url.URL) (*VItem, error) {
	if url.Scheme != Kind {
		return nil, errors.New("not valid google storage URL")
	}
	// /download/storage/v1/b/stowtesttoudhratik/o/a_first%2Fthe%20item
	pieces := strings.SplitN(url.Path, "/", 8)

	vb, err := vc.VBucket(pieces[5])
	if err != nil {
		return nil, ErrNotFound
	}

	vi, err := vb.VItem(pieces[7])
	if err != nil {
		return nil, ErrNotFound
	}
	return vi, nil
}

// Attempts to create a session based on the information given.
func newGoogleStorageClient() (context.Context, *storage.Client, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentials))
	if err != nil {
		return nil, nil, err
	}
	return ctx, client, nil
}
