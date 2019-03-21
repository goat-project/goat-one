package reader

import (
	"context"

	"github.com/onego-project/onego/services"

	"github.com/onego-project/onego"
)

type imageReader struct {
}

func (ir *imageReader) readResources(ctx context.Context, client *onego.Client) ([]resource, error) {
	objs, err := client.ImageService.ListAll(ctx, services.OwnershipFilterAll)

	res := make([]resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}
