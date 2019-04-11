package network

import (
	"sync"

	"github.com/goat-project/goat-one/resource"

	"github.com/goat-project/goat-one/writer"

	"golang.org/x/time/rate"

	"github.com/onego-project/onego/errors"
	"github.com/onego-project/onego/resources"

	pb "github.com/goat-project/goat-proto-go"

	log "github.com/sirupsen/logrus"
)

// Preparer to prepare network data to specific structure for writing to Goat server.
type Preparer struct {
	Writer writer.Writer
}

// CreatePreparer creates Preparer for network records.
func CreatePreparer(limiter *rate.Limiter) *Preparer {
	return &Preparer{
		Writer: *writer.CreateWriter(CreateWriter(limiter)),
	}
}

// InitializeMaps - only for VM relevant.
func (p *Preparer) InitializeMaps(wg *sync.WaitGroup) {
	defer wg.Done()
}

// Preparation prepares network data for writing and call method to write.
func (p *Preparer) Preparation(acc resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := acc.(*resources.VirtualNetwork)
	if ip == nil {
		log.WithFields(log.Fields{"error": errors.ErrNoVirtualNetwork}).Error("error prepare empty network")
		return
	}

	ipRecord := pb.IpRecord{
		// TODO: add attributes and values
	}

	if err := p.Writer.Write(&ipRecord); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error write network record")
	}
}

// SendIdentifier sends identifier to Goat server.
func (p *Preparer) SendIdentifier() error {
	return p.Writer.SendIdentifier()
}

// Finish gets to know to the Goat server that a writing is finished and a response is expected.
// Then, it closes the gRPC connection.
func (p *Preparer) Finish() {
	p.Writer.Finish()
}
