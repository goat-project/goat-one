package network

import (
	"github.com/goat-project/goat-one/constants"
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"
)

// the following tests test additive preparer functions

var _ = ginkgo.Describe("Preparer function test", func() {
	var hook *test.Hook

	ginkgo.JustBeforeEach(func() {
		hook = test.NewGlobal()
	})

	ginkgo.Describe("getSiteName", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getSiteName()).To(gomega.BeEmpty())

				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrNoSiteName))
			})
		})

		ginkgo.Context("when configuration is set correctly", func() {
			ginkgo.It("should return a correct string", func() {
				value := "test-network-site-name"
				viper.SetDefault(constants.CfgNetworkSiteName, value)

				gomega.Expect(getSiteName()).To(gomega.Equal(value))
			})
		})
	})

	ginkgo.Describe("getCloudComputeService", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				gomega.Expect(getCloudComputeService()).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is set correctly", func() {
			ginkgo.It("should return a correct string", func() {
				value := "test-network-cloud-compute-service"
				viper.SetDefault(constants.CfgNetworkCloudComputeService, value)

				gomega.Expect(getCloudComputeService().Value).To(gomega.Equal(value))
			})
		})
	})

	ginkgo.Describe("getCloudType", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getCloudType()).To(gomega.BeEmpty())

				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrNoCloudType))
			})
		})

		ginkgo.Context("when configuration is set correctly", func() {
			ginkgo.It("should return a correct string", func() {
				value := "test-network-cloud-type"
				viper.SetDefault(constants.CfgNetworkCloudType, value)

				gomega.Expect(getCloudType()).To(gomega.Equal(value))
			})
		})
	})

	ginkgo.Describe("getFqan", func() {
		ginkgo.Context("when net user is nil", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getFqan(NetUser{})).To(gomega.BeEmpty())
				// panic?
				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrPrepNoNetUser))
			})
		})

		ginkgo.Context("when user has no group", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getFqan(NetUser{User: resources.CreateUserWithID(1)})).To(gomega.BeEmpty())

				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrNoGroupName))
			})
		})

		// TODO add User XML
		//ginkgo.Context("when configuration is set correctly", func() {
		//	ginkgo.It("should return a correct string", func() {
		//		value := "/" + groupName + "/Role=NULL/Capability=NULL"
		//		viper.SetDefault(constants.CfgNetworkCloudType, value)
		//
		//		gomega.Expect(getFqan(NetUser{User: resources.CreateUserWithID(1)})).To(gomega.Equal(value))
		//	})
		//})
	})

	// TODO test countIPs(user NetUser) (uint32, uint32)
	// add user and vm XMLs

	// TODO test createIPRecord(netUser NetUser, ipType string, ipCount uint32) (*pb.IpRecord, error)
	// add user and vm XMLs
})
