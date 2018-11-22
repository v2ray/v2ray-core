package common

// Closable is the interface for objects that can release its resources.
type Closable interface {
	// Close release all resources used by this object, including goroutines.
	Close() error
}

// Close closes the obj if it is a Closable.
func Close(obj interface{}) error {
	if c, ok := obj.(Closable); ok {
		return c.Close()
	}
	return nil
}

// Runnable is the interface for objects that can start to work and stop on demand.
type Runnable interface {
	// Start starts the runnable object. Upon the method returning nil, the object begins to function properly.
	Start() error

	Closable
}

// HasType is the interface for objects that knows its type.
type HasType interface {
	// Type returns the type of the object.
	// Usually it returns (*Type)(nil) of the object.
	Type() interface{}
}

type ChainedClosable []Closable

func NewChainedClosable(c ...Closable) ChainedClosable {
	return ChainedClosable(c)
}

func (cc ChainedClosable) Close() error {
	for _, c := range cc {
		c.Close()
	}
	return nil
}
