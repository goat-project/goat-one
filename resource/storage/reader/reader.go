package reader

import (
	"context"

	"github.com/goat-project/goat-one/resource"

	"github.com/onego-project/onego/services"

	"github.com/onego-project/onego"
)

// Reader structure for a Reader which read an array of images.
type Reader struct {
}

// ReadResources reads an array of images.
func (ir *Reader) ReadResources(ctx context.Context, client *onego.Client) ([]resource.Resource, error) {
	objs, err := client.ImageService.ListAll(ctx, services.OwnershipFilterAll)
	if err != nil {
		return nil, err
	}

	res := make([]resource.Resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}
