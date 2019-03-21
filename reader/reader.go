package reader

import (
	"net/http"
	"time"

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
