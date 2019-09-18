package storage

import (
	"sync"

	"github.com/goat-project/goat-one/resource"
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Storage Filter tests", func() {
	var (
		filter   *Filter
		res      resource.Resource
		filtered chan resource.Resource
		wg       sync.WaitGroup
	)

	ginkgo.JustBeforeEach(func() {
		filter = CreateFilter()
		wg.Add(1)
	})

	ginkgo.Describe("filter storage", func() {
		ginkgo.Context("when channel is empty and resource correct", func() {
			ginkgo.BeforeEach(func() {
				res = resources.CreateImageWithID(1)
				filtered = make(chan resource.Resource)
			})

			ginkgo.It("should post storage to the channel", func(done ginkgo.Done) {
				go filter.Filtering(res, filtered, &wg)

				gomega.Expect(<-filtered).To(gomega.Equal(res))

				close(done)
			}, 0.2)
		})

		ginkgo.Context("when channel is empty and resource is not correct", func() {
			ginkgo.BeforeEach(func() {
				filtered = make(chan resource.Resource)
			})

			ginkgo.It("should not post storage to the channel", func(done ginkgo.Done) {
				go filter.Filtering(nil, filtered, &wg)

				gomega.Expect(filtered).To(gomega.BeEmpty())

				close(done)
			}, 0.2)
		})

		// TODO add test with full channel
		// we expect that the Filter waits until the channel is empty
		// we need some test with timeout
		//
		// possibly we should test also null channel or null wait group, but that situations should never happen
	})
})
