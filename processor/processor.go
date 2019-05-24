package processor

import (
	"sync"

	"github.com/goat-project/goat-one/resource"
	"github.com/onego-project/onego/errors"
	"github.com/remeh/sizedwaitgroup"

	log "github.com/sirupsen/logrus"
)

// Processor to process resource data.
type Processor struct {
	proc processorI
}

type processorI interface {
	Process(chan resource.Resource, chan bool, *sizedwaitgroup.SizedWaitGroup)
	List(chan resource.Resource, chan bool, *sizedwaitgroup.SizedWaitGroup, int)
	RetrieveInfo(chan resource.Resource, *sync.WaitGroup, resource.Resource)
}

const wgSize = 10

// CreateProcessor creates Processor to manage reading from OpenNebula.
func CreateProcessor(proc processorI) *Processor {
	return &Processor{
		proc: proc,
	}
}

// ListResources calls method to list resource from OpenNebula.
func (p *Processor) ListResources(read chan resource.Resource) {
	swg := sizedwaitgroup.New(wgSize + 1)
	readDone := make(chan bool, wgSize)

	swg.Add()
	go p.proc.Process(read, readDone, &swg)

	swg.Wait()
	close(read)
	close(readDone)
}

// RetrieveInfoResource range over filtered resource and calls method to retrieve resource info.
func (p *Processor) RetrieveInfoResource(filtered, fullInfo chan resource.Resource) {
	var wg sync.WaitGroup

	for accountable := range filtered {
		if accountable == nil {
			log.WithFields(log.Fields{"error": errors.ErrNoVirtualMachine}).Fatal("error retrieve resource info")
		}

		wg.Add(1)
		go p.proc.RetrieveInfo(fullInfo, &wg, accountable)
	}

	wg.Wait()
	close(fullInfo)
}
