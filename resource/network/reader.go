package network

import (
	"context"

	"github.com/goat-project/goat-one/resource"

	"github.com/onego-project/onego"
	"github.com/onego-project/onego/services"
)

// Reader structure for a Reader which read an array of virtual networks.
type Reader struct {
	PageOffset int
}

// ReadResources reads an array of virtual networks.
func (vnr *Reader) ReadResources(ctx context.Context, client *onego.Client) ([]resource.Resource, error) {
	objs, err := client.VirtualNetworkService.List(ctx, vnr.PageOffset, resource.PageSize, services.OwnershipFilterAll)
	if err != nil {
		return nil, err
	}

	res := make([]resource.Resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}
