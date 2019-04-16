package filter

import "github.com/goat-project/goat-one/resource"

// Interface to filter resources.
type Interface interface {
	Filter(chan resource.Resource, chan resource.Resource)
}
