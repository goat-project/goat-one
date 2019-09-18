package writer

import (
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Writer structure to write data to Goat server.
type Writer struct {
	writerI  writerI
	grpcConn *grpc.ClientConn
}

type writerI interface {
	SetUp(*grpc.ClientConn)
	Write(Record) error
	SendIdentifier() error
	Close() (*empty.Empty, error)
}

// CreateWriter creates writer with writer interface and gRPC connection.
func CreateWriter(w writerI, conn *grpc.ClientConn) *Writer {
	w.SetUp(conn)

	return &Writer{
		writerI:  w,
		grpcConn: conn,
	}
}

// Write writes to Goat server.
func (w *Writer) Write(rec Record) error {
	return w.writerI.Write(rec)
}

// SendIdentifier sends identifier to Goat server.
func (w *Writer) SendIdentifier() error {
	return w.writerI.SendIdentifier()
}

// Finish gets to know to the Goat server that a writing is finished and a response is expected.
// Then, it closes the gRPC connection.
func (w *Writer) Finish() {
	// close sending stream
	_, err := w.writerI.Close()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error close and receive")
	}

	// close connection
	err = w.grpcConn.Close()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error close gRPC connection")
	}
}
