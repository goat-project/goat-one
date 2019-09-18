package network_test

import (
	"sync"

	"github.com/goat-project/goat-one/resource/network"

	"github.com/onsi/gomega"

	"github.com/onego-project/onego/resources"

	"github.com/goat-project/goat-one/resource"
	"github.com/onsi/ginkgo"
)

var _ = ginkgo.Describe("Network Filter tests", func() {
	var (
		filter   *network.Filter
		net      resource.Resource
		filtered chan resource.Resource
		wg       sync.WaitGroup
	)

	ginkgo.JustBeforeEach(func() {
		filter = network.CreateFilter()
		wg.Add(1)
	})

	ginkgo.Describe("filter network", func() {
		ginkgo.Context("when channel is empty and resource correct", func() {
			ginkgo.BeforeEach(func() {
				net = createTestNetwork(1)
				filtered = make(chan resource.Resource)
			})

			ginkgo.It("should post network to the channel", func(done ginkgo.Done) {
				go filter.Filtering(net, filtered, &wg)

				gomega.Expect(<-filtered).To(gomega.Equal(net))

				close(done)
			}, 0.2)
		})

		ginkgo.Context("when channel is empty and resource is not correct", func() {
			ginkgo.BeforeEach(func() {
				filtered = make(chan resource.Resource)
			})

			ginkgo.It("should not post network to the channel", func(done ginkgo.Done) {
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

func createTestNetwork(userID int) *network.NetUser {
	vm1 := resources.CreateVirtualMachineWithID(1)
	vm2 := resources.CreateVirtualMachineWithID(2)
	avm := []*resources.VirtualMachine{vm1, vm2}

	return &network.NetUser{User: resources.CreateUserWithID(userID), ActiveVirtualMachines: avm}
}
