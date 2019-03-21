package reader

import (
	"context"

	"github.com/onego-project/onego"
	"github.com/onego-project/onego/services"
)

type vmsReader struct {
	pageOffset int
}

type vmReader struct {
	id int
}

func (vmr *vmsReader) readResources(ctx context.Context, client *onego.Client) ([]resource, error) {
	objs, err := client.VirtualMachineService.List(ctx, vmr.pageOffset, pageSize, services.OwnershipFilterAll,
		services.AnyStateIncludingDone)

	res := make([]resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}

func (vmr *vmReader) readResource(ctx context.Context, client *onego.Client) (resource, error) {
	return client.VirtualMachineService.RetrieveInfo(ctx, vmr.id)
}
