package reader

import (
	"context"
	"net/http"
	"time"

	"github.com/goat-project/goat-one/resource"
	networkReader "github.com/goat-project/goat-one/resource/network/reader"
	storageReader "github.com/goat-project/goat-one/resource/storage/reader"
	virtualMachineReader "github.com/goat-project/goat-one/resource/virtualMachine/reader"

	"github.com/onego-project/onego/resources"

	"github.com/rafaeljesus/retry-go"

	"github.com/goat-project/goat-one/constants"
	"github.com/onego-project/onego/errors"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"

	"github.com/onego-project/onego"
	"golang.org/x/time/rate"
)

// Reader structure to list resources and retrieve info for specific resource from OpenNebula.
type Reader struct {
	client      *onego.Client
	rateLimiter *rate.Limiter
	timeout     time.Duration
}

type resourcesReaderI interface {
	ReadResources(context.Context, *onego.Client) ([]resource.Resource, error)
}

type resourceReaderI interface {
	ReadResource(context.Context, *onego.Client) (resource.Resource, error)
}

const attempts = 3
const sleepTime = time.Second * 1

// CreateReader creates reader with onego client, rate limiter and timeout.
func CreateReader(limiter *rate.Limiter) *Reader {
	// set up connection to OpenNebula
	oneClient := onego.CreateClient(viper.GetString(constants.CfgOpennebulaEndpoint),
		viper.GetString(constants.CfgOpennebulaSecret), &http.Client{})
	if oneClient == nil {
		log.WithFields(log.Fields{"error": errors.ErrNoClient}).Fatal("error create Reader")
	}

	log.WithFields(log.Fields{
		"page-size": resource.PageSize, "attempts": attempts, "sleepTime": sleepTime,
	}).Debug("Reader created with given settings for page size, number of iterations " +
		"for unsuccessful calls and sleep time between the calls")

	return &Reader{
		client:      oneClient,
		rateLimiter: limiter,
		timeout:     viper.GetDuration(constants.CfgOpennebulaTimeout),
	}
}

func (r *Reader) readResources(rri resourcesReaderI) ([]resource.Resource, error) {
	var res []resource.Resource
	var err error

	err = retry.Do(func() error {
		if err = r.rateLimiter.Wait(context.Background()); err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("error list resources")
		}

		ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
		defer cancel()

		res, err = rri.ReadResources(ctx, r.client)
		if err != nil {
			return err
		}

		return nil
	}, attempts, sleepTime)

	return res, err
}

func (r *Reader) readResource(rri resourceReaderI) (resource.Resource, error) {
	var res resource.Resource
	var err error

	err = retry.Do(func() error {
		if err = r.rateLimiter.Wait(context.Background()); err != nil {
			log.WithFields(log.Fields{"error": err}).Panic("error retrieve info")
		}

		ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
		defer cancel()

		res, err = rri.ReadResource(ctx, r.client)
		if err != nil {
			return err
		}

		return nil
	}, attempts, sleepTime)

	return res, err
}

// ListAllVirtualMachines lists all virtual machines by page offset.
func (r *Reader) ListAllVirtualMachines(pageOffset int) ([]*resources.VirtualMachine, error) {
	vmr := virtualMachineReader.VMsReader{
		PageOffset: pageOffset,
	}

	res, err := r.readResources(&vmr)
	if err != nil {
		return nil, err
	}

	vms := make([]*resources.VirtualMachine, len(res))
	for i, e := range res {
		vms[i] = e.(*resources.VirtualMachine)
	}

	return vms, err
}

// RetrieveVirtualMachineInfo returns virtual machines info by id.
func (r *Reader) RetrieveVirtualMachineInfo(id int) (*resources.VirtualMachine, error) {
	vmr := virtualMachineReader.VMReader{
		ID: id,
	}

	res, err := r.readResource(&vmr)
	if err != nil {
		return nil, err
	}

	return res.(*resources.VirtualMachine), err
}

// ListAllUsers lists all users.
func (r *Reader) ListAllUsers() ([]*resources.User, error) {
	or := resource.UserReader{}

	res, err := r.readResources(&or)
	if err != nil {
		return nil, err
	}

	objs := make([]*resources.User, len(res))
	for i, e := range res {
		objs[i] = e.(*resources.User)
	}

	return objs, err
}

// ListAllImages lists all images.
func (r *Reader) ListAllImages() ([]*resources.Image, error) {
	or := storageReader.Reader{}

	res, err := r.readResources(&or)
	if err != nil {
		return nil, err
	}

	objs := make([]*resources.Image, len(res))
	for i, e := range res {
		objs[i] = e.(*resources.Image)
	}

	return objs, err
}

// ListAllHosts lists all hosts.
func (r *Reader) ListAllHosts() ([]*resources.Host, error) {
	or := resource.HostReader{}

	res, err := r.readResources(&or)
	if err != nil {
		return nil, err
	}

	objs := make([]*resources.Host, len(res))
	for i, e := range res {
		objs[i] = e.(*resources.Host)
	}

	return objs, err
}

// ListAllVirtualNetworks lists all virtual networks by page offset.
func (r *Reader) ListAllVirtualNetworks(pageOffset int) ([]*resources.VirtualNetwork, error) {
	or := networkReader.Reader{
		PageOffset: pageOffset,
	}

	res, err := r.readResources(&or)
	if err != nil {
		return nil, err
	}

	objs := make([]*resources.VirtualNetwork, len(res))
	for i, e := range res {
		objs[i] = e.(*resources.VirtualNetwork)
	}

	return objs, err
}

// ListAllClusters lists all clusters.
func (r *Reader) ListAllClusters() ([]*resources.Cluster, error) {
	cr := resource.ClusterReader{}

	res, err := r.readResources(&cr)
	if err != nil {
		return nil, err
	}

	objs := make([]*resources.Cluster, len(res))
	for i, e := range res {
		objs[i] = e.(*resources.Cluster)
	}

	return objs, err
}
