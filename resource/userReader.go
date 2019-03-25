package resource

import (
	"context"

	"github.com/onego-project/onego"
)

// UserReader structure for a Reader which read an array of users.
type UserReader struct {
}

// ReadResources reads an array of users.
func (ur *UserReader) ReadResources(ctx context.Context, client *onego.Client) ([]Resource, error) {
	objs, err := client.UserService.List(ctx)

	res := make([]Resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}
