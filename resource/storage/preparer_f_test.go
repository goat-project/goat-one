package storage

import (
	"github.com/goat-project/goat-one/constants"
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/spf13/viper"
)

// the following tests test additive preparer functions

var _ = ginkgo.Describe("Preparer function test", func() {
	ginkgo.JustBeforeEach(func() {
	})

	ginkgo.Describe("getSiteName", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string", func() {
				gomega.Expect(getSite()).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when configuration is set correctly", func() {
			ginkgo.It("should return a correct string", func() {
				value := "test-storage-site"
				viper.SetDefault(constants.CfgSite, value)

				gomega.Expect(getSite().Value).To(gomega.Equal(value))
			})
		})
	})

	ginkgo.Describe("getStorageShare", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				gomega.Expect(getStorageShare(&resources.Image{})).To(gomega.BeNil())
				// panic?
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				gomega.Expect(getStorageShare(resources.CreateImageWithID(1))).To(gomega.BeNil())
			})
		})

		// TODO add storage XML and test correct getStorageShare()
	})

	ginkgo.Describe("getUID", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				gomega.Expect(getUID(nil)).To(gomega.BeNil())
				// panic?
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				gomega.Expect(getUID(resources.CreateImageWithID(1))).To(gomega.BeNil())
			})
		})

		// TODO add storage XML and test correct getUID()
	})

	ginkgo.Describe("getGID", func() {
		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				gomega.Expect(getGID(&resources.Image{})).To(gomega.BeNil())
				// panic?
			})
		})

		ginkgo.Context("when configuration is not set correctly", func() {
			ginkgo.It("should return an empty string value", func() {
				gomega.Expect(getGID(resources.CreateImageWithID(1))).To(gomega.BeNil())
			})
		})

		// TODO add storage XML and test correct getGID()
	})
})
