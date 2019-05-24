package client

import (
	"sync"

	"github.com/goat-project/goat-one/resource"

	"github.com/goat-project/goat-one/filter"
	"github.com/goat-project/goat-one/preparer"
	"github.com/goat-project/goat-one/processor"
)

// Client runs application.
type Client struct {
}

// Run reads, filters and writes Accountable.
func (c *Client) Run(processor processor.Interface, filter filter.Interface, preparer preparer.Interface) {
	var mapWg sync.WaitGroup
	mapWg.Add(1)

	go preparer.InitializeMaps(&mapWg)

	// initialize channels
	read := make(chan resource.Resource)
	filtered := make(chan resource.Resource)
	fullInfo := make(chan resource.Resource)

	// create done channel
	done := make(chan bool)
	defer close(done)

	go processor.ListResources(read)
	go filter.Filter(read, filtered)
	go processor.RetrieveInfoResource(filtered, fullInfo)
	go preparer.Prepare(fullInfo, done, &mapWg)

	<-done
}
