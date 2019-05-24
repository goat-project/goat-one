package storage

import (
	"sync"

	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/resource"

	"github.com/remeh/sizedwaitgroup"

	log "github.com/sirupsen/logrus"
)

// Processor to process storage data.
type Processor struct {
	reader reader.Reader
}

// CreateProcessor creates Processor to manage reading from OpenNebula.
func CreateProcessor(Reader *reader.Reader) *Processor {
	return &Processor{
		reader: *Reader,
	}
}

// Process provides listing of the storages once.
func (p *Processor) Process(read chan resource.Resource, readDone chan bool, swg *sizedwaitgroup.SizedWaitGroup) {
	defer swg.Done()

	swg.Add()
	go p.List(read, readDone, swg, 0)
}

// List calls method to list all images.
func (p *Processor) List(read chan resource.Resource, _ chan bool, swg *sizedwaitgroup.SizedWaitGroup, _ int) {
	defer swg.Done()

	images, err := p.reader.ListAllImages()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error list images")
	}

	for _, v := range images {
		read <- v
	}
}

// RetrieveInfo - only for VM relevant.
func (p *Processor) RetrieveInfo(fullInfo chan resource.Resource, wg *sync.WaitGroup, image resource.Resource) {
	defer wg.Done()

	fullInfo <- image
}
