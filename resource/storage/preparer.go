package storage

import (
	"sync"

	"github.com/goat-project/goat-one/initialization"

	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/resource"
	"github.com/goat-project/goat-one/writer"

	"golang.org/x/time/rate"

	"github.com/onego-project/onego/errors"
	"github.com/onego-project/onego/resources"

	pb "github.com/goat-project/goat-proto-go"

	log "github.com/sirupsen/logrus"
)

// Preparer to prepare storage data to specific structure for writing to Goat server.
type Preparer struct {
	reader               reader.Reader
	Writer               writer.Writer
	userTemplateIdentity map[int]string
}

// CreatePreparer creates Preparer for storage records.
func CreatePreparer(reader *reader.Reader, limiter *rate.Limiter) *Preparer {
	return &Preparer{
		reader: *reader,
		Writer: *writer.CreateWriter(CreateWriter(limiter)),
	}
}

// InitializeMaps reads additional data for storage record.
func (p *Preparer) InitializeMaps(wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)
	go p.initializeUserTemplateIdentity(wg)
}

func (p *Preparer) initializeUserTemplateIdentity(wg *sync.WaitGroup) {
	defer wg.Done()

	p.userTemplateIdentity = initialization.InitializeUserTemplateIdentity(p.reader)
}

// Preparation prepares storage data for writing and call method to write.
func (p *Preparer) Preparation(acc resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	storage := acc.(*resources.Image)
	if storage == nil {
		log.WithFields(log.Fields{"error": errors.ErrNoImage}).Error("error prepare empty storage")
		return
	}

	storageRecord := pb.StorageRecord{
		// TODO: add attributes and values
	}

	if err := p.Writer.Write(&storageRecord); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error send storage record")
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
