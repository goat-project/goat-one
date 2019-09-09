package network

import (
	"sync"

	"github.com/goat-project/goat-one/constants"

	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/resource"
	"github.com/onego-project/onego/resources"

	"github.com/remeh/sizedwaitgroup"

	log "github.com/sirupsen/logrus"
)

// Processor to process network data.
type Processor struct {
	reader reader.Reader
}

// NetUser represents "Resource" with information about user and his active virtual machines.
type NetUser struct {
	User                  *resources.User
	ActiveVirtualMachines []*resources.VirtualMachine
}

// CreateProcessor creates Processor to manage reading from OpenNebula.
func CreateProcessor(r *reader.Reader) *Processor {
	if r == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreateProcReaderNil)
		return nil
	}

	return &Processor{
		reader: *r,
	}
}

// Process provides listing of the users.
func (p *Processor) Process(read chan resource.Resource, readDone chan bool, swg *sizedwaitgroup.SizedWaitGroup) {
	defer swg.Done()

	users, err := p.reader.ListAllUsers()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error list users")
	}

	for _, user := range users {
		read <- user
	}
}

// RetrieveInfo about virtual machines specific for a given user.
func (p *Processor) RetrieveInfo(fullInfo chan resource.Resource, wg *sync.WaitGroup, user resource.Resource) {
	defer wg.Done()

	id, err := user.ID()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error get user id")
	}

	vms, err := p.reader.ListAllActiveVirtualMachinesForUser(id)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "userID": id}).Fatal("error retrieve virtual machines for user")
	}

	if len(vms) != 0 {
		fullInfo <- &NetUser{
			User:                  user.(*resources.User),
			ActiveVirtualMachines: vms,
		}
	}
}

// ID gets user ID - relevant method to implement "Resource".
func (vnu *NetUser) ID() (int, error) {
	return vnu.User.ID()
}

// Attribute gets user attribute given by path - relevant method to implement "Resource".
func (vnu *NetUser) Attribute(path string) (string, error) {
	return vnu.User.Attribute(path)
}
