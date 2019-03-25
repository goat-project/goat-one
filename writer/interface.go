package writer

// Interface to write records to Goat server.
type Interface interface {
	Write(Record) error
	SendIdentifier() error
}

// Record represents data for writing.
type Record interface {
	Reset()
	String() string
	ProtoMessage()
}
