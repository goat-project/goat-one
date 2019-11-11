package storage_test

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/onego-project/onego/resources"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/goat-project/goat-one/reader"
	"github.com/goat-project/goat-one/util"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/errors"

	"github.com/goat-project/goat-one/resource/storage"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/goat-project/goat-one/constants"
	"github.com/spf13/viper"

	"google.golang.org/grpc"

	"golang.org/x/time/rate"

	"github.com/onsi/gomega"

	"github.com/onsi/ginkgo"

	"cloud.google.com/go/rpcreplay"
)

var recPreparerDir = "records/preparer/"

var _ = ginkgo.Describe("Storage Preparer tests", func() {
	var (
		recName string
		err     error
		rec     *rpcreplay.Recorder
		rep     *rpcreplay.Replayer
		conn    *grpc.ClientConn

		clientHTTP *http.Client
		client     *onego.Client
		vcrRec     *recorder.Recorder

		read *reader.Reader
		prep *storage.Preparer
		wg   sync.WaitGroup
		hook *test.Hook
	)

	ginkgo.JustBeforeEach(func() {
		recPath := recPreparerDir + recName

		// Start gRPC recorder
		if _, err = os.Stat(recPath); os.IsNotExist(err) {
			rec, err = rpcreplay.NewRecorder(recPath, nil)
			if err != nil {
				fmt.Println("unable to create new recorder", err)
				return
			}
			conn, err = grpc.Dial("127.0.0.1:9623", append([]grpc.DialOption{grpc.WithInsecure()}, rec.DialOptions()...)...)
		} else {
			rep, err = rpcreplay.NewReplayer(recPath)
			if err != nil {
				fmt.Println("unable to create new replayer", err)
				return
			}
			conn, err = rep.Connection()
		}

		if err != nil {
			fmt.Println("unable to create connection", err)
			return
		}

		// Start XMLRPC recorder
		vcrRec, err = recorder.New(recPreparerDir + recName)
		if err != nil {
			return
		}

		// Set matcher
		vcrRec.SetMatcher(util.SetMatcher)

		// Create an HTTP client and inject transport
		clientHTTP = &http.Client{
			Transport: vcrRec, // Inject as transport!
		}

		// Create Onego client
		client = onego.CreateClient(constants.OpenNebulaEndpoint, constants.Token, clientHTTP)
		if client == nil {
			err = errors.ErrNoClient
		}

		hook = test.NewGlobal()

		viper.SetDefault(constants.CfgOpennebulaTimeout, constants.OpenNebulaTimeout)
		read = reader.CreateReader(client, rate.NewLimiter(rate.Every(time.Second/time.Duration(30)), 30))

		prep = storage.CreatePreparer(read, rate.NewLimiter(rate.Every(1), 1), conn)
		wg.Add(1)
	})

	ginkgo.AfterEach(func() {
		if rec != nil {
			err = rec.Close()
			if err != nil {
				return // report error
			}
		}

		if rep != nil {
			err = rep.Close()
			if err != nil {
				return // report error
			}
		}
	})

	ginkgo.Describe("create preparer", func() {
		ginkgo.Context("when limiter is correct", func() {
			ginkgo.BeforeEach(func() {
				recName = "createOK"
			})

			ginkgo.It("should create preparer", func() {
				gomega.Expect(conn).NotTo(gomega.BeNil())
				gomega.Expect(read).NotTo(gomega.BeNil())

				p := storage.CreatePreparer(read, rate.NewLimiter(rate.Every(1), 1), conn)

				gomega.Expect(p).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("when limiter is not correct", func() {
			ginkgo.BeforeEach(func() {
				recName = "limiterNil"
			})

			ginkgo.It("should not create preparer", func() {
				gomega.Expect(conn).NotTo(gomega.BeNil())
				gomega.Expect(read).NotTo(gomega.BeNil())

				p := storage.CreatePreparer(read, nil, conn)

				gomega.Expect(p).To(gomega.BeNil())

				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrCreatePrepLimiterNil))
			})
		})

		ginkgo.Context("when connection is not correct", func() {
			ginkgo.BeforeEach(func() {
				recName = "connectionNil"
			})

			ginkgo.It("should not create preparer", func() {
				gomega.Expect(read).NotTo(gomega.BeNil())

				p := storage.CreatePreparer(read, rate.NewLimiter(rate.Every(1), 1), nil)

				gomega.Expect(p).To(gomega.BeNil())

				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrCreatePrepConnNil))
			})
		})

		ginkgo.Context("when reader is not correct", func() {
			ginkgo.BeforeEach(func() {
				recName = "readerNil"
			})

			ginkgo.It("should not create preparer", func() {
				gomega.Expect(conn).NotTo(gomega.BeNil())

				p := storage.CreatePreparer(nil, rate.NewLimiter(rate.Every(1), 1), conn)

				gomega.Expect(p).To(gomega.BeNil())

				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrCreatePrepReaderNil))
			})
		})
	})

	ginkgo.Describe("initialize maps for preparer", func() {
		ginkgo.Context("when reader is OK", func() {
			ginkgo.BeforeEach(func() {
				recName = "initializeMap"
			})

			ginkgo.It("should add map with user template identity", func() {
				prep.InitializeMaps(&wg)

				// TODO map is not visible from this package,
				//  testing in the same package causes import cycle
				//  - add new preparer testing package
				//  - add map getter

				// TODO also test wrong reader
			})
		})
	})

	ginkgo.Describe("prepare storage record", func() {
		ginkgo.Context("when argument is nil", func() {
			ginkgo.BeforeEach(func() {
				recName = "preparationNil"
			})

			ginkgo.It("should not prepare record", func() {
				gomega.Expect(func() { prep.Preparation(nil, &wg) }).To(gomega.Panic())
			})
		})

		ginkgo.Context("when resource has no ID", func() {
			ginkgo.BeforeEach(func() {
				recName = "resourceNil"
			})

			ginkgo.It("should not prepare record", func() {
				prep.Preparation(&resources.Image{}, &wg)

				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrPrepNoImage))
			})
		})

		// The following tests test the usage of various small functions, the change of these functions
		// can effect the behavior of the following tests.

		// TODO add storage XMLs to finish the following tests

		ginkgo.Context("when parameters are correct", func() {
			ginkgo.BeforeEach(func() {
				recName = "preparationOK"
			})

			ginkgo.It("should prepare record", func() {
				image := resources.CreateImageWithID(1) // TODO create from XML

				prep.Preparation(image, &wg)

				// TODO check that record was sent
			})
		})

		// TODO test:
		//  - error send record
		//  - get REGTIME
		//  - get SIZE
	})

	ginkgo.Describe("send identifier", func() {
		ginkgo.Context("when preparer is set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = "sendID"
			})

			ginkgo.It("should send identifier", func() {
				viper.SetDefault(constants.CfgIdentifier, "test-ID")

				gomega.Expect(prep.SendIdentifier()).NotTo(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("finish", func() {
		ginkgo.Context("when preparer is set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = "finish"
			})

			ginkgo.It("should finish the connection", func() {
				viper.SetDefault(constants.CfgIdentifier, "test-ID")
				gomega.Expect(prep.SendIdentifier()).NotTo(gomega.HaveOccurred()) // before finish

				prep.Finish()

				// TODO check the connection was finished and closed
			})
		})
	})
})
