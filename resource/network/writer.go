package network

import (
	"context"

	"github.com/goat-project/goat-one/writer"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/goat-project/goat-one/constants"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"

	pb "github.com/goat-project/goat-proto-go"

	log "github.com/sirupsen/logrus"
)

// Writer structure to write network data to Goat server.
type Writer struct {
	Stream      pb.AccountingService_ProcessIpsClient
	rateLimiter *rate.Limiter
}

// CreateWriter creates Writer for network data.
func CreateWriter(limiter *rate.Limiter) *Writer {
	return &Writer{
		rateLimiter: limiter,
	}
}

// SetUp creates gRPC client and sets up Stream to process networks to Writer.
func (w *Writer) SetUp(conn *grpc.ClientConn) {
	// create grpc client
	grpcClient := pb.NewAccountingServiceClient(conn)

	// create Stream to process VMs
	stream, err := grpcClient.ProcessIps(context.Background())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error create gRPC client stream")
	}

	w.Stream = stream
}

// Write writes network record to Goat server.
func (w *Writer) Write(record writer.Record) error {
	rec := record.(*pb.IpRecord)

	ipData := &pb.IpData{
		Data: &pb.IpData_Ip{
			Ip: rec,
		},
	}

	return w.Stream.Send(ipData)
}

// SendIdentifier sends identifier to Goat server.
func (w *Writer) SendIdentifier() error {
	ipDataIdentifier := pb.IpData_Identifier{Identifier: viper.GetString(constants.CfgIdentifier)}
	data := &pb.IpData{
		Data: &ipDataIdentifier,
	}

	return w.Stream.Send(data)
}

// Close gets to know to the goat server that a writing is finished and a response is expected.
func (w *Writer) Close() (*empty.Empty, error) {
	return w.Stream.CloseAndRecv()
}
