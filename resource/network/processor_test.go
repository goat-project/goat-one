package network_test

import (
	"net/http"
	"sync"
	"time"

	"github.com/onego-project/onego/resources"

	"github.com/goat-project/goat-one/resource"
	"github.com/remeh/sizedwaitgroup"

	"github.com/spf13/viper"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/goat-project/goat-one/constants"
	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/resource/network"
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

var _ = ginkgo.Describe("Network Processor tests", func() {
	var (
		recName    string
		err        error
		clientHTTP *http.Client
		client     *onego.Client
		rec        *recorder.Recorder

		hook *test.Hook

		proc *network.Processor
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

		proc = network.CreateProcessor(read)

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
				p := network.CreateProcessor(read)

				gomega.Expect(p).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("when reader is not correct", func() {
			ginkgo.It("should not create processor", func() {
				p := network.CreateProcessor(nil)

				gomega.Expect(p).To(gomega.BeNil())

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

				x := <-channel
				y := <-channel
				gomega.Expect(x.ID()).To(gomega.Equal(0))
				gomega.Expect(y.ID()).To(gomega.Equal(1))

				close(done)
			}, 0.2)
		})
	})

	ginkgo.Describe("retrieveInfo", func() {
		ginkgo.Context("when channel is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = "retrieveOK"
			})

			ginkgo.It("should post resource to the channel", func(done ginkgo.Done) {
				var wg sync.WaitGroup
				wg.Add(1)

				user := resources.CreateUserWithID(0)

				go proc.RetrieveInfo(channel, &wg, user)

				nu := <-channel
				gomega.Expect(nu.ID()).To(gomega.Equal(0))
				gomega.Expect(nu.Attribute("ID")).To(gomega.Equal("0"))

				netUser := nu.(*network.NetUser)
				gomega.Expect(netUser.ActiveVirtualMachines).To(gomega.HaveLen(3))
				gomega.Expect(netUser.ActiveVirtualMachines[0].ID()).To(gomega.Equal(2))
				gomega.Expect(netUser.ActiveVirtualMachines[1].ID()).To(gomega.Equal(3))
				gomega.Expect(netUser.ActiveVirtualMachines[2].ID()).To(gomega.Equal(5))

				close(done)
			}, 0.2)
		})
	})
})
