package virtualmachine

import (
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/goat-project/goat-one/initialize"

	"github.com/goat-project/goat-one/util"

	"github.com/goat-project/goat-one/resource"

	"github.com/goat-project/goat-one/writer"

	"golang.org/x/time/rate"

	"github.com/goat-project/goat-one/reader"

	"github.com/goat-project/goat-one/constants"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/spf13/viper"

	"github.com/onego-project/onego/errors"
	"github.com/onego-project/onego/resources"

	pb "github.com/goat-project/goat-proto-go"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

// Preparer to prepare virtual machine data to specific structure for writing to Goat server.
type Preparer struct {
	reader                                 reader.Reader
	Writer                                 writer.Writer
	userTemplateIdentity                   map[int]string
	imageTemplateCloudkeeperApplianceMpuri map[int]string
	hostTemplateBenchmarkType              map[int]string
	hostTemplateBenchmarkValue             map[int]string
}

// CreatePreparer creates Preparer for virtual machine records.
func CreatePreparer(reader *reader.Reader, limiter *rate.Limiter, conn *grpc.ClientConn) *Preparer {
	if reader == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreatePrepReaderNil)
		return nil
	}

	if limiter == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreatePrepLimiterNil)
		return nil
	}

	if conn == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreatePrepConnNil)
		return nil
	}

	return &Preparer{
		reader: *reader,
		Writer: *writer.CreateWriter(CreateWriter(limiter), conn),
	}
}

// InitializeMaps reads additional data for virtual machine record.
func (p *Preparer) InitializeMaps(wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(3)

	go func() {
		defer wg.Done()
		p.userTemplateIdentity = initialize.UserTemplateIdentity(p.reader)
	}()

	go func() {
		defer wg.Done()
		p.imageTemplateCloudkeeperApplianceMpuri = initialize.ImageTemplateCloudkeeperApplianceMpuri(p.reader)
	}()

	go func() {
		defer wg.Done()
		p.hostTemplateBenchmarkType, p.hostTemplateBenchmarkType = initialize.HostTemplateBenchmark(p.reader)
	}()
}

// Preparation prepares virtual machine data for writing and call method to write.
func (p *Preparer) Preparation(acc resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	vm := acc.(*resources.VirtualMachine)
	if vm == nil {
		log.WithFields(log.Fields{"error": errors.ErrNoVirtualMachine}).Error(constants.ErrPrepEmptyVM)
		return
	}

	id, err := vm.ID()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error(constants.ErrPrepNoVM)
		return
	}

	machineName, err := getMachineName(vm)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "id": id}).Error(constants.ErrPrepMachineName)
		return
	}

	globalUserName, err := getGlobalUserName(p, vm)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "id": id}).Error(constants.ErrPrepGlobalUserName)
		return
	}

	sTime, err := getStartTime(vm)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "id": id}).Error(constants.ErrPrepSTime)
		return
	}

	eTime := getEndTime(vm)
	wallDuration := getWallDuration(vm)

	vmRecord := pb.VmRecord{
		VmUuid:              uuid.New().String(),
		SiteName:            getSiteName(),
		CloudComputeService: getCloudComputeService(),
		MachineName:         machineName,
		LocalUserId:         getLocalUserID(vm),
		LocalGroupId:        getLocalGroupID(vm),
		GlobalUserName:      globalUserName,
		Fqan:                getFqan(vm),
		Status:              getStatus(vm),
		StartTime:           sTime,
		EndTime:             eTime,
		SuspendDuration:     getSuspendDuration(sTime, eTime, wallDuration),
		WallDuration:        wallDuration,
		CpuDuration:         wallDuration,
		CpuCount:            getCPUCount(vm),
		NetworkType:         getNetworkType(),
		NetworkInbound:      getNetworkInbound(vm),
		NetworkOutbound:     getNetworkOutbound(vm),
		PublicIpCount:       getPublicIPCount(vm),
		Memory:              getMemory(vm),
		Disk:                getDiskSizes(vm),
		BenchmarkType:       getBenchmarkType(p, vm),
		Benchmark:           getBenchmark(p, vm),
		StorageRecordId:     nil,
		ImageId:             getImageID(p, vm),
		CloudType:           getCloudType(),
	}

	if err := p.Writer.Write(&vmRecord); err != nil {
		log.WithFields(log.Fields{"error": err, "id": id}).Error(constants.ErrPrepWrite)
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

func getSiteName() string {
	siteName := viper.GetString(constants.CfgSiteName)
	if siteName == "" {
		log.WithFields(log.Fields{}).Error("no site name in configuration") // should never happen
	}

	return siteName
}

func getCloudComputeService() *wrappers.StringValue {
	return util.CheckValueErrStr(viper.GetString(constants.CfgCloudComputeService), nil)
}

func getMachineName(vm *resources.VirtualMachine) (string, error) {
	deployID, err := vm.DeployID()
	if err != nil {
		return "", err
	}

	return deployID, nil
}

func getLocalUserID(vm *resources.VirtualMachine) *wrappers.StringValue {
	if vm == nil {
		return nil
	}

	return util.CheckValueErrInt(vm.User())
}

func getLocalGroupID(vm *resources.VirtualMachine) *wrappers.StringValue {
	if vm == nil {
		return nil
	}

	return util.CheckValueErrInt(vm.Group())
}

// TODO fix to string (in proto) - global user name is required
func getGlobalUserName(p *Preparer, vm *resources.VirtualMachine) (*wrappers.StringValue, error) {
	userID, err := vm.User()
	if err == nil {
		gun := p.userTemplateIdentity[userID]
		if gun != "" {
			return &wrappers.StringValue{Value: gun}, nil
		}
	}

	return nil, err
}

func getFqan(vm *resources.VirtualMachine) *wrappers.StringValue {
	groupName, err := vm.Attribute("GNAME")
	if err == nil {
		return &wrappers.StringValue{Value: "/" + groupName + "/Role=NULL/Capability=NULL"}
	}

	return nil
}

func getStatus(vm *resources.VirtualMachine) *wrappers.StringValue {
	state, err := vm.State()
	if err == nil {
		return &wrappers.StringValue{Value: resources.VirtualMachineStateMap[state]}
	}

	return nil
}

func getStartTime(vm *resources.VirtualMachine) (*timestamp.Timestamp, error) {
	ts, err := util.CheckTime(vm.STime())
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func getEndTime(vm *resources.VirtualMachine) *timestamp.Timestamp {
	ts, err := util.CheckTime(vm.ETime())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error get end time")
		return nil
	}

	return ts
}

func getSuspendDuration(sTime, eTime *timestamp.Timestamp, wallDuration *duration.Duration) *duration.Duration {
	if eTime != nil && sTime != nil && wallDuration != nil {
		return &duration.Duration{Seconds: eTime.Seconds - sTime.Seconds - wallDuration.Seconds}
	}

	return nil
}

func getWallDuration(vm *resources.VirtualMachine) *duration.Duration {
	if vm.XMLData == nil {
		return nil
	}

	historyRecords, err := vm.HistoryRecords()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error get history records")
		return nil
	}

	currentTime := time.Now().Unix()

	var sum int64
	for _, record := range historyRecords {
		if record != nil {
			rsTime := record.RSTime
			if rsTime == nil {
				continue
			}

			reTime := record.RETime
			if reTime == nil {
				continue
			}

			reTimeUnix := reTime.Unix()
			if reTimeUnix == 0 {
				reTimeUnix = currentTime
			}

			sum += reTimeUnix - rsTime.Unix()
		}
	}

	return &duration.Duration{Seconds: sum}
}

func getCPUCount(vm *resources.VirtualMachine) uint32 {
	vcpu, err := vm.VCPU()
	if err == nil {
		return uint32(vcpu)
	}

	return 0
}

func getNetworkType() *wrappers.StringValue {
	return nil
}

func getNetworkInbound(vm *resources.VirtualMachine) *wrappers.UInt64Value {
	return util.CheckErrUint64(vm.Attribute("MONITORING/NETTX"))
}

func getNetworkOutbound(vm *resources.VirtualMachine) *wrappers.UInt64Value {
	return util.CheckErrUint64(vm.Attribute("MONITORING/NETRX"))
}

func getPublicIPCount(vm *resources.VirtualMachine) *wrappers.UInt64Value {
	if vm.XMLData == nil {
		return nil
	}

	nics, err := vm.NICs()
	if err != nil {
		return nil
	}

	var count uint64
	for _, nic := range nics {
		if util.IsPublicIPv4(nic.IP) || nic.IP6Global != nil {
			count++
		}
	}

	return &wrappers.UInt64Value{Value: count}
}

func getMemory(vm *resources.VirtualMachine) *wrappers.UInt64Value {
	return util.CheckErrUint64(vm.Attribute("TEMPLATE/MEMORY"))
}

func getDiskSizes(vm *resources.VirtualMachine) *wrappers.UInt64Value {
	if vm.XMLData == nil {
		return nil
	}

	disks, err := vm.Disks()
	if err != nil {
		return nil
	}

	var sum uint64

	for _, disk := range disks {
		sum += uint64(disk.Size)
	}

	return &wrappers.UInt64Value{Value: sum}
}

func getBenchmarkType(p *Preparer, vm *resources.VirtualMachine) *wrappers.StringValue {
	if vm == nil || p == nil {
		return nil
	}

	historyRecords, err := vm.HistoryRecords()
	if err == nil && len(historyRecords) > 0 {
		tbt := p.hostTemplateBenchmarkType[*historyRecords[0].HID]
		if tbt != "" {
			return &wrappers.StringValue{Value: tbt}
		}
	}

	return nil
}

func getBenchmark(p *Preparer, vm *resources.VirtualMachine) *wrappers.FloatValue {
	if vm == nil || p == nil {
		return nil
	}

	historyRecords, err := vm.HistoryRecords()
	if err == nil && len(historyRecords) > 0 {
		tbv := p.hostTemplateBenchmarkValue[*historyRecords[0].HID]
		if tbv != "" {
			f, err := strconv.ParseFloat(tbv, 32)
			if err == nil {
				return &wrappers.FloatValue{Value: float32(f)}
			}
		}
	}

	return nil
}

func getImageID(p *Preparer, vm *resources.VirtualMachine) *wrappers.StringValue {
	if vm == nil || p == nil {
		return nil
	}

	disks, err := vm.Disks()
	if err == nil && len(disks) != 0 && disks[0] != nil {
		iid := p.imageTemplateCloudkeeperApplianceMpuri[disks[0].ImageID]
		if iid != "" {
			return &wrappers.StringValue{Value: iid}
		}
	}

	return nil
}

func getCloudType() *wrappers.StringValue {
	ct := viper.GetString(constants.CfgCloudType)
	if ct == "" {
		log.WithFields(log.Fields{}).Error(constants.ErrNoCloudType) // should never happen
	}

	return &wrappers.StringValue{Value: ct}
}
