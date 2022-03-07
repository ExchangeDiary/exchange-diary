package cloudstorage

import (
	"context"
	"io"
	"net/url"
	"time"

	"cloud.google.com/go/storage"
)

// VItem is a Voda Item, which wraps Google Cloud Storage Bucket's file itself.
type VItem struct {
	container    *VBucket        // Container information is required by a few methods.
	client       *storage.Client // A client is needed to make requests.
	name         string
	hash         string
	etag         string
	size         int64
	url          *url.URL
	lastModified time.Time
	metadata     map[string]interface{}
	object       *storage.ObjectAttrs
	ctx          context.Context
}

// ID returns a string value that represents the name of a file.
func (vi *VItem) ID() string {
	return vi.name
}

// Name returns a string value that represents the name of the file.
func (vi *VItem) Name() string {
	return vi.name
}

// Size returns the size of an item in bytes.
func (vi *VItem) Size() (int64, error) {
	return vi.size, nil
}

// URL returns a url which follows the predefined format
func (vi *VItem) URL() *url.URL {
	return vi.url
}

// Open returns an io.ReadCloser to the object. Useful for downloading/streaming the object.
func (vi *VItem) Open() (io.ReadCloser, error) {
	obj := vi.container.Bucket().Object(vi.name)
	return obj.NewReader(vi.ctx)
}

// OpenRange returns an io.Reader to the object for a specific byte range
func (vi *VItem) OpenRange(start, end uint64) (io.ReadCloser, error) {
	obj := vi.container.Bucket().Object(vi.name)
	return obj.NewRangeReader(vi.ctx, int64(start), int64(end-start)+1)
}

// LastMod returns the last modified date of the item.
func (vi *VItem) LastMod() (time.Time, error) {
	return vi.lastModified, nil
}

// Metadata returns a nil map and no error.
func (vi *VItem) Metadata() (map[string]interface{}, error) {
	return vi.metadata, nil
}

// ETag returns the ETag value
func (vi *VItem) ETag() (string, error) {
	return vi.etag, nil
}

// StorageObject returns the Google Storage Object
func (vi *VItem) StorageObject() *storage.ObjectAttrs {
	return vi.object
}

// prepURL takes a MediaLink string and returns a url
func prepURL(str string) (*url.URL, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	u.Scheme = Kind

	// Discard the query string
	u.RawQuery = ""
	return u, nil
}
