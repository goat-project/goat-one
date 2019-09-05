package network_test

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/goat-project/goat-one/constants"
	"github.com/spf13/viper"

	"github.com/onego-project/onego/resources"

	"google.golang.org/grpc"

	"github.com/goat-project/goat-one/resource/network"

	"golang.org/x/time/rate"

	"github.com/onsi/gomega"

	"github.com/onsi/ginkgo"

	"cloud.google.com/go/rpcreplay"
)

var recDir = "records/preparer/"

var _ = ginkgo.Describe("Network Preparer tests", func() {
	var (
		recName string
		err     error
		rec     *rpcreplay.Recorder
		rep     *rpcreplay.Replayer
		conn    *grpc.ClientConn

		prep *network.Preparer
		wg   sync.WaitGroup
		hook *test.Hook
	)

	ginkgo.JustBeforeEach(func() {
		recPath := recDir + recName

		// Start recorder
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

		hook = test.NewGlobal()

		prep = network.CreatePreparer(rate.NewLimiter(rate.Every(1), 1), conn)
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

				p := network.CreatePreparer(rate.NewLimiter(rate.Every(1), 1), conn)

				gomega.Expect(p).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("when limiter is not correct", func() {
			ginkgo.BeforeEach(func() {
				recName = "limiterNil"
			})

			ginkgo.It("should not create preparer", func() {
				gomega.Expect(conn).NotTo(gomega.BeNil())

				p := network.CreatePreparer(nil, conn)

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
				p := network.CreatePreparer(rate.NewLimiter(rate.Every(1), 1), nil)

				gomega.Expect(p).To(gomega.BeNil())

				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrCreatePrepConnNil))
			})
		})
	})

	ginkgo.Describe("initialize maps for preparer", func() {
		ginkgo.Context("not relevant for network", func() {
			ginkgo.BeforeEach(func() {
				recName = "initializeMap"
			})

			ginkgo.It("should do nothing", func() {
				prep.InitializeMaps(&wg)
			})
		})
	})

	ginkgo.Describe("prepare network", func() {
		ginkgo.Context("when argument is nil", func() {
			ginkgo.BeforeEach(func() {
				recName = "preparationNil"
			})

			ginkgo.It("should not prepare record", func() {
				gomega.Expect(func() { prep.Preparation(nil, &wg) }).To(gomega.Panic())
			})
		})

		ginkgo.Context("when net user has no ID", func() {
			ginkgo.BeforeEach(func() {
				recName = "netUserNil"
			})

			ginkgo.It("should not prepare record", func() {
				prep.Preparation(&network.NetUser{}, &wg)

				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrPrepEmptyNetUser))
			})
		})

		// The following tests test the usage of various small functions, the change of these functions
		// can effect the behavior of the following tests.

		// TODO add XML networks to finish the following tests

		ginkgo.Context("when net user has no VM", func() {
			ginkgo.BeforeEach(func() {
				recName = "netUserNoVM"
			})

			ginkgo.It("should not prepare record", func() {
				netUser := &network.NetUser{User: resources.CreateUserWithID(1)}

				prep.Preparation(netUser, &wg)

				// TODO check that no record was sent
			})
		})

		ginkgo.Context("when net user has one VM", func() {
			ginkgo.BeforeEach(func() {
				recName = "netUserOneVM"
			})

			ginkgo.It("should prepare record", func() {
				vms := []*resources.VirtualMachine{resources.CreateVirtualMachineWithID(5)} // TODO create from XML
				netUser := &network.NetUser{User: resources.CreateUserWithID(1), ActiveVirtualMachines: vms}

				prep.Preparation(netUser, &wg)

				// TODO check that record which contains one IP was sent
			})
		})

		ginkgo.Context("when net user has more VMs", func() {
			ginkgo.BeforeEach(func() {
				recName = "netUserMoreVMs"
			})

			ginkgo.It("should prepare record", func() {
				vms := []*resources.VirtualMachine{resources.CreateVirtualMachineWithID(5),
					resources.CreateVirtualMachineWithID(6), resources.CreateVirtualMachineWithID(7),
					resources.CreateVirtualMachineWithID(8)} // TODO create from XML
				netUser := &network.NetUser{User: resources.CreateUserWithID(1), ActiveVirtualMachines: vms}

				prep.Preparation(netUser, &wg)

				// TODO check that record which contains 4 IPs was sent
			})
		})
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
