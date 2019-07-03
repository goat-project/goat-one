package preparer

import (
	"sync"

	"github.com/goat-project/goat-one/resource"
	log "github.com/sirupsen/logrus"
)

// Preparer to prepare data to specific structure for writing to Goat server.
type Preparer struct {
	prep preparerI
}

type preparerI interface {
	Preparation(resource.Resource, *sync.WaitGroup)
	InitializeMaps(*sync.WaitGroup)
	SendIdentifier() error
	Finish()
}

// CreatePreparer creates Preparer for accountable records.
func CreatePreparer(prep preparerI) *Preparer {
	return &Preparer{
		prep: prep,
	}
}

// Prepare gets networks from channel and call method to prepare network record and send.
func (p *Preparer) Prepare(fullInfo chan resource.Resource, done chan bool, mapWg *sync.WaitGroup) {
	mapWg.Wait()

	var wg sync.WaitGroup

	identifierSend := false

	for data := range fullInfo {
		if !identifierSend {
			err := p.prep.SendIdentifier()
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Fatal("error send identifier")
			}
			identifierSend = true
		}
		wg.Add(1)
		go p.prep.Preparation(data, &wg)
	}

	wg.Wait()

	// If the identifier was not sent, there is no resource to prepare and send,
	// a gRPC connection was not open and no finishing and closing of a connection are needed.
	if identifierSend {
		p.prep.Finish()
	}

	done <- true
}

// InitializeMaps reads additional data for record.
func (p *Preparer) InitializeMaps(wg *sync.WaitGroup) {
	p.prep.InitializeMaps(wg)
}
