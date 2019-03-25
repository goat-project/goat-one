package reader

import (
	"context"

	"github.com/onego-project/onego"
)

type userReader struct {
}

func (ur *userReader) readResources(ctx context.Context, client *onego.Client) ([]resource, error) {
	objs, err := client.UserService.List(ctx)

	res := make([]resource, len(objs))
	for i, e := range objs {
		res[i] = e
	}

	return res, err
}
