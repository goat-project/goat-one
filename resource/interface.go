package resource

// Resource interface represents resource from OpenNebula.
type Resource interface {
	ID() (int, error)
	Attribute(path string) (string, error)
}

// PageSize is a number of read resources in one call. It is used for resources with pagination.
const PageSize = 100
