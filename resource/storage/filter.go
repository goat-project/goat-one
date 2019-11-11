package storage

import (
	"sync"

	"github.com/goat-project/goat-one/resource"
)

// Filter to filter storage data.
type Filter struct {
}

// CreateFilter creates Filter.
func CreateFilter() *Filter {
	return &Filter{}
}

// Filtering - only for VM relevant.
func (f *Filter) Filtering(storage resource.Resource, filtered chan resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	if storage == nil {
		return
	}

	filtered <- storage
}
