package reader

import (
	"context"

	"github.com/onego-project/onego"
	"github.com/onego-project/onego/services"
)

type vnetReader struct {
	pageOffset int
}

func (vnr *vnetReader) readResources(ctx context.Context, client *onego.Client) ([]resource, error) {
	objs, err := client.VirtualNetworkService.List(ctx, vnr.pageOffset, pageSize, services.OwnershipFilterAll)

	res := make([]resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}
