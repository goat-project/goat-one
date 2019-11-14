package virtualmachine

import (
	"github.com/beevik/etree"
	"github.com/goat-project/goat-one/constants"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"
)

// the following tests test additive preparer functions

var _ = ginkgo.Describe("Preparer function test", func() {
	var (
		hook *test.Hook
		doc  *etree.Document
	)

	ginkgo.JustBeforeEach(func() {
		hook = test.NewGlobal()

		doc = etree.NewDocument()
		gomega.Expect(doc.ReadFromFile("test/xml/vm.xml")).NotTo(gomega.HaveOccurred())

		viper.SetDefault(constants.CfgSiteName, "")
		viper.SetDefault(constants.CfgCloudComputeService, "")
		viper.SetDefault(constants.CfgCloudType, "")
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
				value := "test-vm-site"
				viper.SetDefault(constants.CfgSiteName, value)

				gomega.Expect(getSiteName()).To(gomega.Equal(value))
			})
		})
	})

	ginkgo.Describe("getCloudComputeService", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getCloudComputeService()).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is set correctly", func() {
			ginkgo.It("should return a correct string", func() {
				value := "test-CloudComputeService"
				viper.SetDefault(constants.CfgCloudComputeService, value)

				gomega.Expect(getCloudComputeService().GetValue()).To(gomega.Equal(value))
			})
		})
	})

	ginkgo.Describe("getMachineName", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				name, err := getMachineName(&resources.VirtualMachine{})

				gomega.Expect(name).To(gomega.BeEmpty())
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				name, err := getMachineName(resources.CreateVirtualMachineWithID(1))

				gomega.Expect(name).To(gomega.BeEmpty())
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when configuration is set correctly", func() {
			ginkgo.It("should return a string value", func() {
				name, err := getMachineName(resources.CreateVirtualMachineFromXML(doc.Root()))

				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(name).To(gomega.Equal("one-57502"))
			})
		})
	})

	ginkgo.Describe("getLocalUserID", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				gomega.Expect(getLocalUserID(nil)).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				gomega.Expect(getLocalUserID(resources.CreateVirtualMachineWithID(1))).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getLocalUserID(resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal("46"))
			})
		})
	})

	ginkgo.Describe("getLocalGroupID", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				gomega.Expect(getLocalGroupID(nil)).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				gomega.Expect(getLocalGroupID(resources.CreateVirtualMachineWithID(1))).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getLocalGroupID(resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal("113"))
			})
		})
	})

	ginkgo.Describe("getGlobalUserName", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				value, err := getGlobalUserName(nil, &resources.VirtualMachine{})

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(value.GetValue()).To(gomega.BeEmpty())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				uti := map[int]string{
					46: "world",
				}
				preparer := &Preparer{userTemplateIdentity: uti}

				value, err := getGlobalUserName(preparer, resources.CreateVirtualMachineWithID(1))

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(value.GetValue()).To(gomega.BeEmpty())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				uti := map[int]string{
					46: "world",
				}
				preparer := &Preparer{userTemplateIdentity: uti}

				value, err := getGlobalUserName(preparer, resources.CreateVirtualMachineFromXML(doc.Root()))

				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(value.GetValue()).To(gomega.Equal("world"))
			})
		})
	})

	ginkgo.Describe("getFqan", func() {
		ginkgo.Context("when net user is nil", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getFqan(&resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getFqan(resources.CreateVirtualMachineWithID(1))).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getFqan(resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal("/cloud-devel/Role=NULL/Capability=NULL"))
			})
		})
	})

	ginkgo.Describe("getStatus", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getStatus(&resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getStatus(resources.CreateVirtualMachineWithID(1))).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getStatus(resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal("ACTIVE"))
			})
		})
	})

	ginkgo.Describe("getStartTime", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				t, err := getStartTime(&resources.VirtualMachine{})

				gomega.Expect(t).To(gomega.BeNil())
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				t, err := getStartTime(resources.CreateVirtualMachineWithID(1))

				gomega.Expect(t).To(gomega.BeNil())
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when configuration is set correctly", func() {
			ginkgo.It("should return a string value", func() {
				t, err := getStartTime(resources.CreateVirtualMachineFromXML(doc.Root()))

				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(t.GetSeconds()).To(gomega.Equal(int64(1540931164)))
			})
		})
	})

	ginkgo.Describe("getEndTime", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getEndTime(&resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getEndTime(resources.CreateVirtualMachineWithID(1))).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getEndTime(resources.CreateVirtualMachineFromXML(doc.Root())).GetSeconds()).To(
					gomega.Equal(int64(0)))
			})
		})
	})

	ginkgo.Describe("getSuspendDuration", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getSuspendDuration(nil, nil, nil)).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getSuspendDuration(&timestamp.Timestamp{Seconds: 1573641110}, &timestamp.Timestamp{Seconds: 1573643810}, &duration.Duration{Seconds: 1000}).GetSeconds()).To(
					gomega.Equal(int64(1700)))
			})
		})
	})

	ginkgo.Describe("getWallDuration", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getWallDuration(&resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(
					getWallDuration(resources.CreateVirtualMachineWithID(1)).GetSeconds()).To(
					gomega.Equal(int64(0)))
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getWallDuration(resources.CreateVirtualMachineFromXML(doc.Root())).GetSeconds()).To(
					gomega.Equal(int64(7707605)))
			})
		})
	})

	ginkgo.Describe("getCPUCount", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getCPUCount(&resources.VirtualMachine{})).To(gomega.Equal(uint32(0)))
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(
					getCPUCount(resources.CreateVirtualMachineWithID(1))).To(gomega.Equal(uint32(0)))
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(getCPUCount(resources.CreateVirtualMachineFromXML(doc.Root()))).To(gomega.Equal(uint32(1)))
			})
		})
	})

	ginkgo.Describe("getNetworkType", func() {
		ginkgo.Context("when configuration is set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getNetworkType()).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("getNetworkInbound", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getNetworkInbound(&resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getNetworkInbound(resources.CreateVirtualMachineWithID(1))).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getNetworkInbound(resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal(uint64(48708945)))
			})
		})
	})

	ginkgo.Describe("getNetworkOutbound", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getNetworkOutbound(&resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getNetworkOutbound(resources.CreateVirtualMachineWithID(1))).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getNetworkOutbound(resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal(uint64(12983215634)))
			})
		})
	})

	ginkgo.Describe("getPublicIPCount", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getPublicIPCount(&resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getPublicIPCount(resources.CreateVirtualMachineWithID(1)).GetValue()).To(gomega.Equal(uint64(0)))
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getPublicIPCount(resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal(uint64(1)))
			})
		})
	})

	ginkgo.Describe("getMemory", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getMemory(&resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getMemory(resources.CreateVirtualMachineWithID(1))).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getMemory(resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal(uint64(2048)))
			})
		})
	})

	ginkgo.Describe("getDiskSizes", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getDiskSizes(&resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getDiskSizes(resources.CreateVirtualMachineWithID(1)).GetValue()).To(gomega.Equal(uint64(0)))
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				gomega.Expect(
					getDiskSizes(resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal(uint64(13312)))
			})
		})
	})

	ginkgo.Describe("getBenchmarkType", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getBenchmarkType(nil, &resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				htbt := map[int]string{
					932: "hello",
				}
				preparer := &Preparer{hostTemplateBenchmarkType: htbt}

				gomega.Expect(getBenchmarkType(preparer, resources.CreateVirtualMachineWithID(1)).GetValue()).To(gomega.BeEmpty())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				htbt := map[int]string{
					932: "hello",
				}
				preparer := &Preparer{hostTemplateBenchmarkType: htbt}

				gomega.Expect(
					getBenchmarkType(preparer, resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal("hello"))
			})
		})
	})

	ginkgo.Describe("getBenchmark", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getBenchmark(nil, &resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				htbv := map[int]string{
					932: "100",
				}
				preparer := &Preparer{hostTemplateBenchmarkValue: htbv}

				gomega.Expect(getBenchmark(preparer, resources.CreateVirtualMachineWithID(1)).GetValue()).To(gomega.Equal(float32(0)))
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				htbv := map[int]string{
					932: "100",
				}
				preparer := &Preparer{hostTemplateBenchmarkValue: htbv}

				gomega.Expect(
					getBenchmark(preparer, resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal(float32(100)))
			})
		})
	})

	ginkgo.Describe("getImageID", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getImageID(nil, &resources.VirtualMachine{})).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				itcam := map[int]string{
					7161: "world",
				}
				preparer := &Preparer{imageTemplateCloudkeeperApplianceMpuri: itcam}

				gomega.Expect(getImageID(preparer, resources.CreateVirtualMachineWithID(1)).GetValue()).To(gomega.BeEmpty())
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return a string value", func() {
				itcam := map[int]string{
					7161: "world",
				}
				preparer := &Preparer{imageTemplateCloudkeeperApplianceMpuri: itcam}

				gomega.Expect(
					getImageID(preparer, resources.CreateVirtualMachineFromXML(doc.Root())).GetValue()).To(
					gomega.Equal("world"))
			})
		})
	})

	ginkgo.Describe("getCloudType", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getCloudType().GetValue()).To(gomega.BeEmpty())

				gomega.Expect(hook.LastEntry().Level).To(gomega.Equal(logrus.ErrorLevel))
				gomega.Expect(hook.LastEntry().Message).To(gomega.Equal(constants.ErrNoCloudType))
			})
		})

		ginkgo.Context("when configuration is set correctly", func() {
			ginkgo.It("should return a correct string", func() {
				value := "test-cloud-type"
				viper.SetDefault(constants.CfgCloudType, value)

				gomega.Expect(getCloudType().GetValue()).To(gomega.Equal(value))
			})
		})
	})
})
