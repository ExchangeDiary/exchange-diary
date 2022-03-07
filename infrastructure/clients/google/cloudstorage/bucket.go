package cloudstorage

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

// VBucket is Voda Bucket which is wrapper of google cloud storage "Bucket"
type VBucket struct {
	// Name is needed to retrieve items.
	name string

	// Client is responsible for performing the requests.
	client *storage.Client

	// ctx is used on google storage API calls
	ctx context.Context
}

// ID returns a string value which represents the name of the container.
// TODO: UUID
func (vb *VBucket) ID() string {
	return vb.name
}

// Name returns a string value which represents the name of the container.
func (vb *VBucket) Name() string {
	return vb.name
}

// Bucket returns the google bucket attributes
func (vb *VBucket) Bucket() *storage.BucketHandle {
	return vb.client.Bucket(vb.name)
}

// VItem returns VodaItem (*VItem)
func (vb *VBucket) VItem(id string) (*VItem, error) {
	item, err := vb.Bucket().Object(id).Attrs(vb.ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return vb.convertToVItem(item)
}

// RemoveItem will delete a google storage Object
func (vb *VBucket) RemoveItem(id string) error {
	return vb.Bucket().Object(id).Delete(vb.ctx)
}

// Put sends a request to upload content to the container. The arguments
// received are the name of the item, a reader representing the
// content, and the size of the file.
func (vb *VBucket) Put(name string, r io.Reader, size int64, metadata map[string]interface{}) (*VItem, error) {
	obj := vb.Bucket().Object(name)

	mdPrepped, err := prepMetadata(metadata)
	if err != nil {
		return nil, err
	}

	w := obj.NewWriter(vb.ctx)
	if _, err := io.Copy(w, r); err != nil {
		return nil, err
	}
	w.Close()

	attr, err := obj.Update(vb.ctx, storage.ObjectAttrsToUpdate{Metadata: mdPrepped})
	if err != nil {
		return nil, err
	}

	return vb.convertToVItem(attr)
}

func (vb *VBucket) convertToVItem(attr *storage.ObjectAttrs) (*VItem, error) {
	u, err := prepUrl(attr.MediaLink)
	if err != nil {
		return nil, err
	}

	mdParsed, err := parseMetadata(attr.Metadata)
	if err != nil {
		return nil, err
	}

	return &VItem{
		name:         attr.Name,
		container:    vb,
		client:       vb.client,
		size:         attr.Size,
		etag:         attr.Etag,
		hash:         string(attr.MD5),
		lastModified: attr.Updated,
		url:          u,
		metadata:     mdParsed,
		object:       attr,
		ctx:          vb.ctx,
	}, nil
}

func parseMetadata(metadataParsed map[string]string) (map[string]interface{}, error) {
	metadataParsedMap := make(map[string]interface{}, len(metadataParsed))
	for key, value := range metadataParsed {
		metadataParsedMap[key] = value
	}
	return metadataParsedMap, nil
}

func prepMetadata(metadataParsed map[string]interface{}) (map[string]string, error) {
	returnMap := make(map[string]string, len(metadataParsed))
	for key, value := range metadataParsed {
		str, ok := value.(string)
		if !ok {
			return nil, errors.Errorf(`value of key '%s' in metadata must be of type string`, key)
		}
		returnMap[key] = str
	}
	return returnMap, nil
}
