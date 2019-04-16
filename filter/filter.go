package filter

import (
	"sync"

	"github.com/goat-project/goat-one/resource"
)

// Filter to filter resource data.
type Filter struct {
	filterI filterI
}

type filterI interface {
	Filtering(res resource.Resource, filtered chan resource.Resource, wg *sync.WaitGroup)
}

// CreateFilter creates Filter.
func CreateFilter(filterI filterI) *Filter {
	return &Filter{
		filterI: filterI,
	}
}

// Filter reads resources from read channel, filter them according to configuration or command line flags
// and write them to filtered channel.
func (f *Filter) Filter(read, filtered chan resource.Resource) {
	var wg sync.WaitGroup

	for data := range read {
		wg.Add(1)
		go f.filterI.Filtering(data, filtered, &wg)
	}

	wg.Wait()
	close(filtered)
}
