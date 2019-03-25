package resource

import (
	"context"

	"github.com/onego-project/onego"
)

// HostReader structure for a Reader which read an array of hosts.
type HostReader struct {
}

// ReadResources reads an array of hosts.
func (hr *HostReader) ReadResources(ctx context.Context, client *onego.Client) ([]Resource, error) {
	objs, err := client.HostService.List(ctx)

	res := make([]Resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}
