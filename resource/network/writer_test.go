package network_test

import (
	"fmt"
	"os"

	"cloud.google.com/go/rpcreplay"
	"github.com/goat-project/goat-one/resource/network"
	goat_grpc "github.com/goat-project/goat-proto-go"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

var recWriterDir = "records/writer/"

var _ = ginkgo.Describe("Network Writer tests", func() {
	var (
		recName string
		err     error
		rec     *rpcreplay.Recorder
		rep     *rpcreplay.Replayer
		conn    *grpc.ClientConn

		writer *network.Writer
	)

	ginkgo.JustBeforeEach(func() {
		recPath := recWriterDir + recName

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

		// create correct writer
		writer = network.CreateWriter(rate.NewLimiter(rate.Every(1), 1))
		writer.SetUp(conn)
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

	ginkgo.Describe("write", func() {
		ginkgo.Context("when record is correct", func() {
			ginkgo.BeforeEach(func() {
				recName = "recordOK"
			})

			ginkgo.It("should write record", func() {
				record := &goat_grpc.IpRecord{
					SiteName: "test-site-name",
				}

				gomega.Expect(writer.Write(record)).NotTo(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when record is nil", func() {
			ginkgo.BeforeEach(func() {
				recName = "recordNil"
			})

			ginkgo.It("should not write record", func() {
				gomega.Expect(func() { _ = writer.Write(nil) }).To(gomega.Panic())
			})
		})

		ginkgo.Context("when record is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = "recordEmpty"
			})

			ginkgo.It("should write record because server ignores empty records", func() {
				gomega.Expect(writer.Write(&goat_grpc.IpRecord{})).NotTo(gomega.HaveOccurred())
			})
		})
	})
})
