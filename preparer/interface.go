package preparer

import (
	"sync"

	"github.com/goat-project/goat-one/resource"
)

// Interface to prepare data to specific structure for writing to Goat server.
type Interface interface {
	InitializeMaps(*sync.WaitGroup)
	Prepare(chan resource.Resource, chan bool, *sync.WaitGroup)
}
