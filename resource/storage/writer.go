package storage

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

// Writer structure to write storage data to Goat server.
type Writer struct {
	Stream      pb.AccountingService_ProcessStoragesClient
	rateLimiter *rate.Limiter
}

// CreateWriter creates Writer for storage data.
func CreateWriter(limiter *rate.Limiter) *Writer {
	return &Writer{
		rateLimiter: limiter,
	}
}

// SetUp creates gRPC client and sets up Stream to process storages to Writer.
func (w *Writer) SetUp(conn *grpc.ClientConn) {
	// create gRPC client
	grpcClient := pb.NewAccountingServiceClient(conn)

	// create Stream to process VMs
	stream, err := grpcClient.ProcessStorages(context.Background())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Panic("error create gRPC client stream")
	}

	w.Stream = stream
}

// Write writes network record to Goat server.
func (w *Writer) Write(record writer.Record) error {
	rec := record.(*pb.StorageRecord)

	storageData := &pb.StorageData{
		Data: &pb.StorageData_Storage{
			Storage: rec,
		},
	}

	return w.Stream.Send(storageData)
}

// SendIdentifier sends identifier to Goat server.
func (w *Writer) SendIdentifier() error {
	storageDataIdentifier := pb.StorageData_Identifier{Identifier: viper.GetString(constants.CfgIdentifier)}
	data := &pb.StorageData{
		Data: &storageDataIdentifier,
	}

	return w.Stream.Send(data)
}

// Close gets to know to the goat server that a writing is finished and a response is expected.
func (w *Writer) Close() (*empty.Empty, error) {
	return w.Stream.CloseAndRecv()
}
