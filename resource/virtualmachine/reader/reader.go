package reader

import (
	"context"

	"github.com/goat-project/goat-one/resource"

	"github.com/onego-project/onego"
	"github.com/onego-project/onego/services"
)

// VMsReader structure for a Reader which read an array of virtual machines.
type VMsReader struct {
	PageOffset int
}

// VMReader structure for a Reader which read virtual machine by id.
type VMReader struct {
	ID int
}

// ReadResources reads an array of virtual machines.
func (vmr *VMsReader) ReadResources(ctx context.Context, client *onego.Client) ([]resource.Resource, error) {
	objs, err := client.VirtualMachineService.List(ctx, vmr.PageOffset, resource.PageSize, services.OwnershipFilterAll,
		services.AnyStateIncludingDone)
	if err != nil {
		return nil, err
	}

	res := make([]resource.Resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}

// ReadResource reads a virtual machine.
func (vmr *VMReader) ReadResource(ctx context.Context, client *onego.Client) (resource.Resource, error) {
	return client.VirtualMachineService.RetrieveInfo(ctx, vmr.ID)
}
