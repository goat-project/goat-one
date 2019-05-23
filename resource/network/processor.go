package network

import (
	"sync"

	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/resource"

	"github.com/remeh/sizedwaitgroup"

	log "github.com/sirupsen/logrus"
)

// Processor to process network data.
type Processor struct {
	reader reader.Reader
}

// CreateProcessor creates Processor to manage reading from OpenNebula.
func CreateProcessor(Reader *reader.Reader) *Processor {
	return &Processor{
		reader: *Reader,
	}
}

// Process provides listing of the networks with pagination.
func (p *Processor) Process(read chan resource.Resource, readDone chan bool, swg *sizedwaitgroup.SizedWaitGroup) {
	defer swg.Done()
	pageOffset := 1

processing:
	for {
		swg.Add()
		go p.List(read, readDone, swg, pageOffset)
		select {
		case <-readDone:
			break processing
		default:
		}

		pageOffset++
	}
}

// List calls method to list virtual networks by page offset.
func (p *Processor) List(read chan resource.Resource, readDone chan bool, swg *sizedwaitgroup.SizedWaitGroup,
	pageOffset int) {
	defer swg.Done()

	vnets, err := p.reader.ListAllVirtualNetworks(pageOffset)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error list networks")
	}

	if len(vnets) == 0 {
		readDone <- true
		return
	}

	for _, v := range vnets {
		read <- v
	}
}

// RetrieveInfo - only for VM relevant.
func (p *Processor) RetrieveInfo(fullInfo chan resource.Resource, wg *sync.WaitGroup, vnet resource.Resource) {
	defer wg.Done()

	fullInfo <- vnet
}
