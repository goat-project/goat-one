package reader

import (
	"context"

	"github.com/onego-project/onego"
)

type clusterReader struct {
}

func (cr *clusterReader) readResources(ctx context.Context, client *onego.Client) ([]resource, error) {
	objs, err := client.ClusterService.List(ctx)

	res := make([]resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}
