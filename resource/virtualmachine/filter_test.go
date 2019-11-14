package virtualmachine

import (
	"sync"
	"time"

	"github.com/beevik/etree"

	"github.com/onego-project/onego/resources"

	"github.com/goat-project/goat-one/constants"
	"github.com/goat-project/goat-one/resource"
	"github.com/spf13/viper"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Virtual machine Filter tests", func() {
	var (
		filter   *Filter
		res      resource.Resource
		filtered chan resource.Resource
		wg       sync.WaitGroup

		doc *etree.Document
	)

	ginkgo.JustBeforeEach(func() {
		doc = etree.NewDocument()
		gomega.Expect(doc.ReadFromFile("test/xml/vm.xml")).NotTo(gomega.HaveOccurred())

		viper.SetDefault(constants.CfgRecordsFrom, time.Time{})
		viper.SetDefault(constants.CfgRecordsTo, time.Time{})
		viper.SetDefault(constants.CfgRecordsForPeriod, time.Time{})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when no values are set", func() {
			ginkgo.It("should create filter with no restrictions", func() {
				filter := CreateFilter()

				gomega.Expect(filter.recordsFrom).To(gomega.Equal(time.Time{}))
				gomega.Expect(filter.recordsTo).To(gomega.And(
					gomega.BeTemporally("<", time.Now()),
					gomega.BeTemporally(">", time.Now().Add(-time.Minute))))
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when time from is set", func() {
			ginkgo.It("should create filter", func() {
				dateFrom := time.Now().Add(-48 * time.Hour)
				viper.SetDefault(constants.CfgRecordsFrom, dateFrom)

				filter := CreateFilter()

				gomega.Expect(filter.recordsFrom).To(gomega.Equal(dateFrom))
				gomega.Expect(filter.recordsTo).To(gomega.And(
					gomega.BeTemporally("<", time.Now()),
					gomega.BeTemporally(">", time.Now().Add(-time.Minute))))
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when time from and to are set", func() {
			ginkgo.It("should create filter", func() {
				dateFrom := time.Now().Add(-48 * time.Hour)
				dateTo := time.Now().Add(-24 * time.Hour)

				viper.SetDefault(constants.CfgRecordsFrom, dateFrom)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)

				filter := CreateFilter()

				gomega.Expect(filter.recordsFrom).To(gomega.Equal(dateFrom))
				gomega.Expect(filter.recordsTo).To(gomega.Equal(dateTo))
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when time to is set", func() {
			ginkgo.It("should create filter", func() {
				dateTo := time.Now().Add(-48 * time.Hour)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)

				filter := CreateFilter()

				gomega.Expect(filter.recordsFrom).To(gomega.Equal(time.Time{}))
				gomega.Expect(filter.recordsTo).To(gomega.Equal(dateTo))
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when time from and to and period are set", func() {
			ginkgo.It("should not create filter", func() {
				dateFrom := time.Now().Add(-48 * time.Hour)
				dateTo := time.Now().Add(-24 * time.Hour)
				period := "1y"

				viper.SetDefault(constants.CfgRecordsFrom, dateFrom)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)
				viper.SetDefault(constants.CfgRecordsForPeriod, period)

				// TODO test Fatal error
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when period is set", func() {
			ginkgo.It("should create filter", func() {
				period := "1y"
				viper.SetDefault(constants.CfgRecordsForPeriod, period)

				filter := CreateFilter()

				expectation := time.Now().Add(-365 * 24 * time.Hour)

				gomega.Expect(filter.recordsFrom).To(gomega.And(
					gomega.BeTemporally("<", expectation),
					gomega.BeTemporally(">", expectation.Add(-time.Minute))))

				gomega.Expect(filter.recordsTo).To(gomega.And(
					gomega.BeTemporally("<", time.Now()),
					gomega.BeTemporally(">", time.Now().Add(-time.Minute))))
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when period and time to are set", func() {
			ginkgo.It("should create filter", func() {
				dateTo := time.Now().Add(-24 * time.Hour)
				period := "1y"

				viper.SetDefault(constants.CfgRecordsForPeriod, period)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)

				// TODO test Fatal error
			})
		})
	})

	ginkgo.Describe("create filter", func() {
		ginkgo.Context("when period and time from are set", func() {
			ginkgo.It("should create filter", func() {
				dateFrom := time.Now().Add(-48 * time.Hour)
				period := "1y"

				viper.SetDefault(constants.CfgRecordsForPeriod, period)
				viper.SetDefault(constants.CfgRecordsFrom, dateFrom)

				// TODO test Fatal error
			})
		})
	})

	ginkgo.Describe("filter virtual machine", func() {
		ginkgo.Context("when channel is empty and resource correct", func() {
			ginkgo.It("should not post vm to the channel", func(done ginkgo.Done) {
				res = resources.CreateVirtualMachineWithID(1)
				filtered = make(chan resource.Resource)

				filter = CreateFilter()

				wg.Add(1)
				go filter.Filtering(res, filtered, &wg)

				gomega.Expect(filtered).To(gomega.BeEmpty())

				close(done)
			}, 0.2)
		})

		ginkgo.Context("when channel is empty and resource is not correct", func() {
			ginkgo.It("should not post vm to the channel", func(done ginkgo.Done) {
				filtered = make(chan resource.Resource)

				filter = CreateFilter()

				wg.Add(1)
				go filter.Filtering(nil, filtered, &wg)

				gomega.Expect(filtered).To(gomega.BeEmpty())

				close(done)
			}, 0.2)
		})

		// TODO add test with full channel

		ginkgo.Context("when channel is empty and resource time is in range", func() {
			ginkgo.It("should post vm to the channel", func(done ginkgo.Done) {
				dateTo := time.Now().Add(-24 * time.Hour)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)

				filter = CreateFilter()

				res = resources.CreateVirtualMachineFromXML(doc.Root())
				filtered = make(chan resource.Resource)

				wg.Add(1)
				go filter.Filtering(res, filtered, &wg)

				gomega.Expect(<-filtered).To(gomega.Equal(res))

				close(done)
			}, 0.2)
		})

		ginkgo.Context("when channel is empty and resource time is out of range", func() {
			ginkgo.It("should not post vm to the channel", func(done ginkgo.Done) {
				dateTo := time.Now().Add(-2 * 356 * 24 * time.Hour)
				viper.SetDefault(constants.CfgRecordsTo, dateTo)

				filter = CreateFilter()

				res = resources.CreateVirtualMachineFromXML(doc.Root())
				filtered = make(chan resource.Resource)

				wg.Add(1)
				go filter.Filtering(res, filtered, &wg)

				gomega.Expect(filtered).To(gomega.BeEmpty())

				close(done)
			}, 0.2)
		})
	})
})
