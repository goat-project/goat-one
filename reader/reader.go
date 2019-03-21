package reader

import (
	"context"
	"net/http"
	"time"

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
	readResources(context.Context, *onego.Client) ([]resource, error)
}

type resourceReaderI interface {
	readResource(context.Context, *onego.Client) (resource, error)
}

type resource interface {
	ID() (int, error)
}

const pageSize = 100

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
		"page-size": pageSize, "attempts": attempts, "sleepTime": sleepTime,
	}).Debug("Reader created with given settings for page size, number of iterations " +
		"for unsuccessful calls and sleep time between the calls")

	return &Reader{
		client:      oneClient,
		rateLimiter: limiter,
		timeout:     viper.GetDuration(constants.CfgOpennebulaTimeout),
	}
}

func (r *Reader) readResources(rri resourcesReaderI) ([]resource, error) {
	var res []resource
	var err error

	err = retry.Do(func() error {
		if err = r.rateLimiter.Wait(context.Background()); err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("error list resources")
		}

		ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
		defer cancel()

		res, err = rri.readResources(ctx, r.client)

		if err != nil {
			return err
		}

		return nil
	}, attempts, sleepTime)

	return res, err
}

func (r *Reader) readResource(rri resourceReaderI) (resource, error) {
	var res resource
	var err error

	err = retry.Do(func() error {
		if err = r.rateLimiter.Wait(context.Background()); err != nil {
			log.WithFields(log.Fields{"error": err}).Panic("error retrieve info")
		}

		ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
		defer cancel()

		res, err = rri.readResource(ctx, r.client)

		if err != nil {
			return err
		}

		return nil
	}, attempts, sleepTime)

	return res, err
}

// ListAllVirtualMachines lists all virtual machines by page offset.
func (r *Reader) ListAllVirtualMachines(pageOffset int) ([]*resources.VirtualMachine, error) {
	vmr := vmsReader{
		pageOffset: pageOffset,
	}

	res, err := r.readResources(&vmr)

	vms := make([]*resources.VirtualMachine, len(res))
	for i, e := range res {
		vms[i] = e.(*resources.VirtualMachine)
	}

	return vms, err
}

// RetrieveVirtualMachineInfo returns virtual machines info by id.
func (r *Reader) RetrieveVirtualMachineInfo(id int) (*resources.VirtualMachine, error) {
	vmr := vmReader{id: id}
	res, err := r.readResource(&vmr)

	return res.(*resources.VirtualMachine), err
}

// ListAllUsers lists all users.
func (r *Reader) ListAllUsers() ([]*resources.User, error) {
	or := userReader{}
	res, err := r.readResources(&or)

	objs := make([]*resources.User, len(res))
	for i, e := range res {
		objs[i] = e.(*resources.User)
	}

	return objs, err
}

// ListAllImages lists all images.
func (r *Reader) ListAllImages() ([]*resources.Image, error) {
	or := imageReader{}
	res, err := r.readResources(&or)

	objs := make([]*resources.Image, len(res))
	for i, e := range res {
		objs[i] = e.(*resources.Image)
	}

	return objs, err
}

// ListAllHosts lists all hosts.
func (r *Reader) ListAllHosts() ([]*resources.Host, error) {
	or := hostReader{}
	res, err := r.readResources(&or)

	objs := make([]*resources.Host, len(res))
	for i, e := range res {
		objs[i] = e.(*resources.Host)
	}

	return objs, err
}

// ListAllVirtualNetworks lists all virtual networks by page offset.
func (r *Reader) ListAllVirtualNetworks(pageOffset int) ([]*resources.VirtualNetwork, error) {
	or := vnetReader{
		pageOffset: pageOffset,
	}
	res, err := r.readResources(&or)

	objs := make([]*resources.VirtualNetwork, len(res))
	for i, e := range res {
		objs[i] = e.(*resources.VirtualNetwork)
	}

	return objs, err
}
