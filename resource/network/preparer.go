package network

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/goat-project/goat-one/util"
	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/goat-project/goat-one/constants"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/spf13/viper"

	"github.com/goat-project/goat-one/resource"

	"github.com/goat-project/goat-one/writer"

	"golang.org/x/time/rate"

	pb "github.com/goat-project/goat-proto-go"
	log "github.com/sirupsen/logrus"
)

// Preparer to prepare network data to specific structure for writing to Goat server.
type Preparer struct {
	Writer writer.Writer
}

// CreatePreparer creates Preparer for network records.
func CreatePreparer(limiter *rate.Limiter, conn *grpc.ClientConn) *Preparer {
	if limiter == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreatePrepLimiterNil)
		return nil
	}

	if conn == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrCreatePrepConnNil)
		return nil
	}

	return &Preparer{
		Writer: *writer.CreateWriter(CreateWriter(limiter), conn),
	}
}

// InitializeMaps - only for VM relevant.
func (p *Preparer) InitializeMaps(wg *sync.WaitGroup) {
	defer wg.Done()
}

// Preparation prepares network data for writing and call method to write.
func (p *Preparer) Preparation(acc resource.Resource, wg *sync.WaitGroup) {
	defer wg.Done()

	netUser := acc.(*NetUser)
	if netUser == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrPrepEmptyNetUser)
		return
	}

	id, err := netUser.ID()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error(constants.ErrPrepNoNetUser)
		return
	}

	countIPv4, countIPv6 := countIPs(*netUser)

	fmt.Println(id, len(netUser.ActiveVirtualMachines), countIPv4, countIPv6)

	if countIPv4 != 0 {
		ipv4Record, err := createIPRecord(*netUser, "IPv4", countIPv4)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "user-id": id}).Error(constants.ErrPrepIPv4)
			return
		}

		if err := p.Writer.Write(ipv4Record); err != nil {
			log.WithFields(log.Fields{"error": err}).Error(constants.ErrPrepWrite)
		}
	}

	if countIPv6 != 0 {
		ipv6Record, err := createIPRecord(*netUser, "IPv6", countIPv6)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "user-id": id}).Error(constants.ErrPrepIPv6)
			return
		}

		if err := p.Writer.Write(ipv6Record); err != nil {
			log.WithFields(log.Fields{"error": err}).Error(constants.ErrPrepWrite)
		}
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
	siteName := viper.GetString(constants.CfgNetworkSiteName)
	if siteName == "" {
		log.WithFields(log.Fields{}).Error(constants.ErrNoSiteName) // should never happen
	}

	return siteName
}

func getCloudComputeService() *wrappers.StringValue {
	return util.CheckValueErrStr(viper.GetString(constants.CfgNetworkCloudComputeService), nil)
}

func getCloudType() string {
	ct := viper.GetString(constants.CfgNetworkCloudType)
	if ct == "" {
		log.WithFields(log.Fields{}).Error(constants.ErrNoCloudType) // should never happen
	}

	return ct
}

func getFqan(netUser NetUser) string {
	if netUser.User == nil {
		log.WithFields(log.Fields{}).Error(constants.ErrPrepNoNetUser)
		return ""
	}

	groupName, err := netUser.User.Attribute("GNAME")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error(constants.ErrNoGroupName)
		return ""
	}

	return "/" + groupName + "/Role=NULL/Capability=NULL"
}

func countIPs(user NetUser) (uint32, uint32) {
	var countIPv4 uint32
	var countIPv6 uint32

	for _, vm := range user.ActiveVirtualMachines {
		nics, err := vm.NICs()
		if err != nil {
			continue
		}

		for _, nic := range nics {
			if util.IsPublicIPv4(nic.IP) {
				countIPv4++
			} else if nic.IP6Global != nil {
				countIPv6++
			}
		}
	}

	return countIPv4, countIPv6
}

func createIPRecord(netUser NetUser, ipType string, ipCount uint32) (*pb.IpRecord, error) {
	id, err := netUser.ID()
	if err != nil {
		return nil, err
	}

	gid, err := netUser.User.MainGroup()
	if err != nil {
		return nil, err
	}

	globalUserName, err := netUser.User.Attribute(constants.TemplateIdentity)
	if err != nil {
		globalUserName = strconv.Itoa(id)
	}

	return &pb.IpRecord{
		MeasurementTime:     &timestamp.Timestamp{Seconds: time.Now().Unix()},
		SiteName:            getSiteName(),
		CloudComputeService: getCloudComputeService(),
		CloudType:           getCloudType(),
		LocalUser:           strconv.Itoa(id),
		LocalGroup:          strconv.Itoa(gid),
		GlobalUserName:      globalUserName,
		Fqan:                getFqan(netUser),
		IpType:              ipType,
		IpCount:             ipCount,
	}, nil
}
