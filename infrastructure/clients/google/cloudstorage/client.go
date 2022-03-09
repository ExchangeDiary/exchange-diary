package cloudstorage

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/spf13/viper"
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
	clientOnce        sync.Once
)

// VClient ...
type VClient struct {
	client *storage.Client
	ctx    context.Context
}

func init() {
	logger.Info("lazy init cloud storage")
	// LazyGlobal loading ...
	clientOnce.Do(func() {
		var client *storage.Client
		var err error
		ctx := context.Background()
		switch viper.GetString("PHASE") {
		case "prod":
			client, err = storage.NewClient(ctx)
		default:
			client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentials))
		}

		if err != nil {
			panic("Failed to load google storage  " + err.Error())
		}

		vodaStorageClient = &VClient{
			client: client,
			ctx:    ctx,
		}
	})

}

// GetVClient ...
func GetVClient() *VClient {
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
	// /download/storage/v1/b/voda_bucket/o/1%2Fimages%2Ftest2.png?jk=AFshE3UavexT0J0Ay69_MPDC0BeoHt4sFhybW7svkRp2mmcLxTRmat38Ta7FUd8y4pt4bWZ1qnoQ0WOEhKvRbBXNCbTbn8q-m7EMnnPAkHV-zDf976ZLsmw5BcnPh0g4-dcDKMxtGpkSxf1fGEq3dTm2dtQ1AMjQVjl0y4Hum1zbrX1Wz00YPuqanxCb2vfmL8Y_UpNEDj9KFD6iZkX-fB6swgQB_gHSRxNQEKnNXakRBFJFiy_SMScGswu5IcfzDq2-f3Zhw8Uzb8U-C2n7W7Rre6hoEPAq744pTH007xUHS9Bai1NgI2kgFYgpmdfFezTht-2jCyHYiknkuBrcFk2glUCBOQLFFqU7YhRjGpC7bPadPWxDEKFG74X67W1W_JaQbier86Mpu3Y1e1WEPiVzPNUqtnFIjIJoYUYvAzKmS_hQfLt4hjBMRo3GsbczL6wuhUPp-IDot3H51qVWsk-5tyPT9UGb5H8kySjbCISgxMhJhsOe8cpmElBY0L6Biq-L-UHMnmkFqYlY6IBNravtTh4ycESBUbWGJU4iXddbMtUdRCNjcJBN_xbTsofC6r4MEiYj4teBl6qfijeZcO84wJW3fpkqxWJ2MRes6U-7eeGKvqf5heZd4GYy3KY-oJf5fUYpiqxwJUhhz5DpXXcyNBAAJ-Ul-PA-YaLg9IvcY7eI8er4r56akxMged9LVqFSF6U7NpbGiMJUEjg8Z1f_QMlE0HjoWZV8vQaIFujU_rIfC-F6Oe3qz3keDnhYgffhHzitJqHTMamd6llr5i0v921SreKgej-LN07YFGFKtns72SyBMW7odB2EtUfHPO7UaM0Z_8G4C9Kcrl8JhYxTgpi71a1yH3NorkKrdqINnNWPY7Y9Dpxd51hYelTPVeMLA3T1wD4LkEJihj6fxG-Mu0x1xB-cL37a20yi12KoYBFg-u83bhiVZyZkjHigZqyYiOUO6fxaZP0&isca=1
	pieces := strings.SplitN(url.Path, "/", 8)
	bucketID := pieces[5]
	itemID := pieces[7]

	vb, err := vc.VBucket(bucketID)
	if err != nil {
		return nil, ErrNotFound
	}

	vi, err := vb.VItem(itemID)
	if err != nil {
		return nil, ErrNotFound
	}
	return vi, nil
}
