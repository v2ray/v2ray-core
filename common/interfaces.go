package common

import "v2ray.com/core/common/errors"

// Closable is the interface for objects that can release its resources.
//
// v2ray:api:beta
type Closable interface {
	// Close release all resources used by this object, including goroutines.
	Close() error
}

// Interruptible is an interface for objects that can be stopped before its completion.
//
// v2ray:api:beta
type Interruptible interface {
	Interrupt()
}

// Close closes the obj if it is a Closable.
//
// v2ray:api:beta
func Close(obj interface{}) error {
	if c, ok := obj.(Closable); ok {
		return c.Close()
	}
	return nil
}

// Interrupt calls Interrupt() if object implements Interruptible interface, or Close() if the object implements Closable interface.
//
// v2ray:api:beta
func Interrupt(obj interface{}) error {
	if c, ok := obj.(Interruptible); ok {
		c.Interrupt()
		return nil
	}
	return Close(obj)
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

// ChainedClosable is a Closable that consists of multiple Closable objects.
type ChainedClosable []Closable

// Close implements Closable.
func (cc ChainedClosable) Close() error {
	var errs []error
	for _, c := range cc {
		if err := c.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Combine(errs...)
}
