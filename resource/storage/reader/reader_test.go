package reader

import (
	"context"
	"net/http"

	"github.com/goat-project/goat-one/util"

	"github.com/goat-project/goat-one/constants"
	"github.com/goat-project/goat-one/resource"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/errors"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"github.com/dnaeon/go-vcr/recorder"

	"testing"
)

func TestResources(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Storage reader Suite")
}

var _ = ginkgo.Describe("Storage Reader tests", func() {
	var (
		clientHTTP *http.Client
		client     *onego.Client
		rec        *recorder.Recorder
		err        error
		reader     *Reader
		resources  []resource.Resource
	)

	ginkgo.AfterEach(func() {
		err = rec.Stop()
		if err != nil {
			return // report error
		}
	})

	ginkgo.Describe("read images from OpenNebula", func() {
		ginkgo.Context("when endpoint is wrong", func() {
			ginkgo.BeforeEach(func() {
				// no record for wrong endpoint
				rec, err = recorder.New("")
				if err != nil {
					return
				}

				// Create Onego client
				client = onego.CreateClient(constants.WrongOpenNebulaEndpoint, constants.Token, &http.Client{})
				if client == nil {
					err = errors.ErrNoClient
				}

				// Create Reader
				reader = &Reader{}
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				resources, err = reader.ReadResources(context.TODO(), client)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(resources).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when connection is correct, but username is wrong", func() {
			ginkgo.BeforeEach(func() {
				// Start recorder
				rec, err = recorder.New("records/readWithWrongUsername")
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
				client = onego.CreateClient(constants.OpenNebulaEndpoint, constants.WrongNameToken, clientHTTP)
				if client == nil {
					err = errors.ErrNoClient
				}

				// Create Reader
				reader = &Reader{}
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				resources, err = reader.ReadResources(context.TODO(), client)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(resources).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when connection is correct, but user password is wrong", func() {
			ginkgo.BeforeEach(func() {
				// Start recorder
				rec, err = recorder.New("records/readWithWrongPassword")
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
				client = onego.CreateClient(constants.OpenNebulaEndpoint, constants.WrongPswdToken, clientHTTP)
				if client == nil {
					err = errors.ErrNoClient
				}

				// Create Reader
				reader = &Reader{}
			})

			ginkgo.It("should return an error", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				resources, err = reader.ReadResources(context.TODO(), client)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(resources).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("read images from OpenNebula with correct client", func() {
		var recName string

		ginkgo.JustBeforeEach(func() {
			// Start recorder
			rec, err = recorder.New(recName)
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
		})

		ginkgo.Context("when connection is correct", func() {
			ginkgo.BeforeEach(func() {
				recName = "records/readImages"

				// Create Reader
				reader = &Reader{}
			})

			ginkgo.It("should return a list of Resources", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				resources, err = reader.ReadResources(context.TODO(), client)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resources).ShouldNot(gomega.BeNil())
				gomega.Expect(resources).ShouldNot(gomega.BeEmpty())
				gomega.Expect(resources).Should(gomega.HaveLen(constants.NumTestedNetworks))
			})
		})
	})
})
