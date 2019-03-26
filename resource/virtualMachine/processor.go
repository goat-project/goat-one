package virtualMachine

import (
	"sync"

	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/resource"

	"github.com/remeh/sizedwaitgroup"

	log "github.com/sirupsen/logrus"
)

// Processor to process virtual machine data.
type Processor struct {
	reader reader.Reader
}

// CreateProcessor creates processor with reader.
func CreateProcessor(Reader *reader.Reader) *Processor {
	return &Processor{
		reader: *Reader,
	}
}

// List calls method to list virtual machines by page offset.
func (p *Processor) List(read chan resource.Resource, readDone chan bool, swg *sizedwaitgroup.SizedWaitGroup,
	pageOffset int) {
	defer swg.Done()

	vms, err := p.reader.ListAllVirtualMachines(pageOffset)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "page-offset": pageOffset}).Fatal("error list virtual machines")
	}

	if len(vms) == 0 {
		readDone <- true
		return
	}

	for _, v := range vms {
		read <- v
	}
}

// RetrieveInfo calls method to retrieve virtual machine info.
func (p *Processor) RetrieveInfo(fullInfo chan resource.Resource, wg *sync.WaitGroup, vm resource.Resource) {
	defer wg.Done()

	id, err := vm.ID()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error get virtual machine id")
	}

	v, err := p.reader.RetrieveVirtualMachineInfo(id)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error retrieve virtual machine info")
	}

	fullInfo <- v
}
