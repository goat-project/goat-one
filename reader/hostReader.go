package reader

import (
	"context"

	"github.com/onego-project/onego"
)

type hostReader struct {
}

func (hr *hostReader) readResources(ctx context.Context, client *onego.Client) ([]resource, error) {
	objs, err := client.HostService.List(ctx)

	res := make([]resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}
