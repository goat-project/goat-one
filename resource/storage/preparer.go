package storage

import (
	"sync"
	"time"

	"github.com/goat-project/goat-one/initialization"

	"github.com/goat-project/goat-one/constants"
	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/util"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/spf13/viper"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/goat-project/goat-one/resource"
	"github.com/goat-project/goat-one/writer"

	"golang.org/x/time/rate"

	"github.com/onego-project/onego/errors"
	"github.com/onego-project/onego/resources"

	pb "github.com/goat-project/goat-proto-go"

	log "github.com/sirupsen/logrus"

	"github.com/beevik/guid"
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

	id, err := storage.ID()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error prepare storage record")
		return
	}

	startTime, err := getStartTime(storage)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "id": id}).Error("error get REGTIME, unable to prepare storage")
		return
	}

	size, err := getResourceCapacityUsed(storage)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "id": id}).Error("error get SIZE, unable to prepare storage")
		return
	}

	now := time.Now().Unix()

	storageRecord := pb.StorageRecord{
		RecordID:      guid.New().String(),
		CreateTime:    &timestamp.Timestamp{Seconds: now},
		StorageSystem: viper.GetString(constants.CfgOpennebulaEndpoint),
		Site:          getSite(),
		StorageShare:  getStorageShare(storage),
		StorageMedia:  &wrappers.StringValue{Value: "disk"},
		// StorageClass: nil,
		FileCount: &wrappers.StringValue{Value: "1"},
		// DirectoryPath: nil,
		LocalUser:    getUID(storage),
		LocalGroup:   getGID(storage),
		UserIdentity: getUserIdentity(p, storage),
		Group:        getGroup(storage),
		// GroupAttribute: nil,
		// GroupAttributeType: nil,
		StartTime:                 startTime,
		EndTime:                   &timestamp.Timestamp{Seconds: now},
		ResourceCapacityUsed:      size,
		LogicalCapacityUsed:       &wrappers.UInt64Value{Value: size},
		ResourceCapacityAllocated: &wrappers.UInt64Value{Value: size},
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

func getSite() *wrappers.StringValue {
	return util.CheckValueErrStr(viper.GetString(constants.CfgSite), nil)
}

func getStorageShare(storage *resources.Image) *wrappers.StringValue {
	return util.CheckValueErrStr(storage.Attribute("DATASTORE"))
}

func getUID(storage *resources.Image) *wrappers.StringValue {
	return util.CheckValueErrInt(storage.User())
}

func getGID(storage *resources.Image) *wrappers.StringValue {
	return util.CheckValueErrInt(storage.Group())
}

func getUserIdentity(p *Preparer, storage *resources.Image) *wrappers.StringValue {
	uid, err := storage.User()
	if err == nil {
		ui := p.userTemplateIdentity[uid]
		if ui != "" {
			return &wrappers.StringValue{Value: ui}
		}
	}

	return nil
}

func getGroup(storage *resources.Image) *wrappers.StringValue {
	groupName, err := storage.Attribute("GNAME")
	if err == nil {
		return &wrappers.StringValue{Value: "/" + groupName + "/Role=NULL/Capability=NULL"}
	}

	return nil
}

func getStartTime(storage *resources.Image) (*timestamp.Timestamp, error) {
	rs, err := util.CheckTime(storage.RegistrationTime())
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func getResourceCapacityUsed(storage *resources.Image) (uint64, error) {
	size, err := storage.Size()
	if err != nil {
		return 0, err
	}

	return uint64(size * 1024), nil
}
