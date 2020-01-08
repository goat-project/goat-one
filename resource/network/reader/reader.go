package reader

import (
	"context"

	"github.com/onego-project/onego/resources"

	"github.com/goat-project/goat-one/resource"

	"github.com/onego-project/onego"
	"github.com/onego-project/onego/services"
)

//// Reader structure for a Reader which read an array of virtual networks.
//type Reader struct {
//	PageOffset int
//}

// VMReader structure for a Reader which read an array of virtual machines specific for a user given by id.
type VMReader struct {
	User *resources.User
}

//// ReadResources reads an array of virtual networks.
//func (vnr *Reader) ReadResources(ctx context.Context, client *onego.Client) ([]resource.Resource, error) {
//	objs, err := client.VirtualNetworkService.List(ctx, vnr.PageOffset, resource.PageSize, services.OwnershipFilterAll)
//	if err != nil {
//		return nil, err
//	}
//
//	res := make([]resource.Resource, len(objs))
//	for i, e := range objs {
//		res[i] = e
//	}
//
//	return res, err
//}

// ReadResourcesForUser reads an array of virtual machines specific for a user given by id.
func (vmr *VMReader) ReadResourcesForUser(ctx context.Context, client *onego.Client) ([]resource.Resource, error) {
	objs, err := client.VirtualMachineService.ListAllForUser(ctx, *vmr.User, services.Active)
	if err != nil {
		return nil, err
	}

	res := make([]resource.Resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}
