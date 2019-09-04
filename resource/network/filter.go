package network

import (
	"sync"

	"github.com/goat-project/goat-one/resource"
)

// Filter to filter network data.
type Filter struct {
}

// CreateFilter creates Filter.
func CreateFilter() *Filter {
	return &Filter{}
}

// Filtering - only for VM relevant.
func (f *Filter) Filtering(network resource.Resource, filtered chan resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	if network == nil {
		return
	}

	filtered <- network
}
