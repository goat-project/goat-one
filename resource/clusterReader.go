package resource

import (
	"context"

	"github.com/onego-project/onego"
)

// ClusterReader structure for a Reader which read an array of clusters.
type ClusterReader struct {
}

// ReadResources reads an array of clusters.
func (cr *ClusterReader) ReadResources(ctx context.Context, client *onego.Client) ([]Resource, error) {
	objs, err := client.ClusterService.List(ctx)

	res := make([]Resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}
