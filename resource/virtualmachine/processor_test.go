package virtualmachine_test

import (
	"net/http"
	"sync"
	"time"

	"github.com/remeh/sizedwaitgroup"

	"github.com/onego-project/onego/resources"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"

	"github.com/goat-project/goat-one/resource"
	"github.com/goat-project/goat-one/resource/virtualmachine"

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
)

var recProcessorDir = "test/records/processor/"

var _ = ginkgo.Describe("Virtual machine Processor tests", func() {
	var (
		recName    string
		err        error
		clientHTTP *http.Client
		client     *onego.Client
		rec        *recorder.Recorder

		hook *test.Hook

		proc *virtualmachine.Processor
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
	})

	ginkgo.AfterEach(func() {
		err = rec.Stop()
		if err != nil {
			return // report error
		}
	})

	ginkgo.Describe("create processor", func() {
		ginkgo.Context("when reader is correct", func() {
			ginkgo.It("should create processor", func() {
				p := virtualmachine.CreateProcessor(read)

				gomega.Expect(p).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("when reader is not correct", func() {
			ginkgo.It("should not create processor", func() {
				p := virtualmachine.CreateProcessor(nil)

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
				proc = virtualmachine.CreateProcessor(read)
				channel = make(chan resource.Resource)
				readDone := make(chan bool)

				swg = sizedwaitgroup.New(3)
				swg.Add()
				go proc.Process(channel, readDone, &swg)

				for i := 0; i < 6; i++ {
					gomega.Expect((<-channel).ID()).To(gomega.Equal(i))
				}

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
				proc = virtualmachine.CreateProcessor(read)

				channel = make(chan resource.Resource)

				user := resources.CreateVirtualMachineWithID(0)

				var wg sync.WaitGroup
				wg.Add(1)
				go proc.RetrieveInfo(channel, &wg, user)

				vm := <-channel
				gomega.Expect(vm.ID()).To(gomega.Equal(0))

				close(done)
			}, 0.2)
		})
	})
})
