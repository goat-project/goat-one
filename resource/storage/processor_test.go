package storage_test

import (
	"net/http"
	"sync"
	"time"

	"github.com/goat-project/goat-one/resource/storage"

	"github.com/onego-project/onego/resources"

	"github.com/goat-project/goat-one/resource"
	"github.com/remeh/sizedwaitgroup"

	"github.com/spf13/viper"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/goat-project/goat-one/constants"
	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/util"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/errors"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"golang.org/x/time/rate"
)

var recProcessorDir = "records/processor/"

var _ = ginkgo.Describe("Storage Processor tests", func() {
	var (
		recName    string
		err        error
		clientHTTP *http.Client
		client     *onego.Client
		rec        *recorder.Recorder

		hook *test.Hook

		proc *storage.Processor
		read *reader.Reader

		channel chan resource.Resource

		swg sizedwaitgroup.SizedWaitGroup
	)

	ginkgo.JustBeforeEach(func() {
		// Start recorder
		rec, err = recorder.New(recProcessorDir + recName)
		if err != nil {
			return
		}

		// Set matcher
		rec.SetMatcher(util.SetMatcher)

		// Create an HTTP client and inject transport
		clientHTTP = &http.Client{
			Transport: rec, // Inject as transport!
		}

		// Create Onego client
		client = onego.CreateClient(constants.OpenNebulaEndpoint, constants.Token, clientHTTP)
		if client == nil {
			err = errors.ErrNoClient
		}

		hook = test.NewGlobal()

		viper.SetDefault(constants.CfgOpennebulaTimeout, constants.OpenNebulaTimeout)

		read = reader.CreateReader(client, rate.NewLimiter(rate.Every(time.Second/time.Duration(30)), 30))

		proc = storage.CreateProcessor(read)

		channel = make(chan resource.Resource)
	})

	ginkgo.AfterEach(func() {
		err = rec.Stop()
		if err != nil {
			return // report error
		}
	})

	ginkgo.Describe("create processor", func() {
		ginkgo.Context("when read is correct", func() {
			ginkgo.It("should create processor", func() {
				gomega.Expect(storage.CreateProcessor(read)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("when reader is not correct", func() {
			ginkgo.It("should not create processor", func() {
				gomega.Expect(storage.CreateProcessor(nil)).To(gomega.BeNil())

				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrCreateProcReaderNil))
			})
		})
	})

	ginkgo.Describe("process", func() {
		ginkgo.Context("when channel is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = "channelOK"
			})

			ginkgo.It("should post resource to the channel", func(done ginkgo.Done) {
				go proc.Process(channel, nil, &swg)

				gomega.Expect((<-channel).ID()).To(gomega.Equal(0))

				close(done)
			}, 0.2)
		})
	})

	ginkgo.Describe("retrieveInfo", func() {
		ginkgo.Context("when channel is empty", func() {
			ginkgo.It("should post resource to the channel", func(done ginkgo.Done) {
				var wg sync.WaitGroup
				wg.Add(1)

				go proc.RetrieveInfo(channel, &wg, resources.CreateImageWithID(0))

				gomega.Expect((<-channel).ID()).To(gomega.Equal(0))

				close(done)
			}, 0.2)
		})
	})
})
