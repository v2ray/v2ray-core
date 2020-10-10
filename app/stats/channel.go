// +build !confonly

package stats

import (
	"context"
	"sync"

	"v2ray.com/core/common"
)

// Channel is an implementation of stats.Channel.
type Channel struct {
	channel     chan channelMessage
	subscribers []chan interface{}

	// Synchronization components
	access sync.RWMutex
	closed chan struct{}

	// Channel options
	blocking   bool // Set blocking state if channel buffer reaches limit
	bufferSize int  // Set to 0 as no buffering
	subsLimit  int  // Set to 0 as no subscriber limit
}

// NewChannel creates an instance of Statistics Channel.
func NewChannel(config *ChannelConfig) *Channel {
	return &Channel{
		channel:    make(chan channelMessage, config.BufferSize),
		subsLimit:  int(config.SubscriberLimit),
		bufferSize: int(config.BufferSize),
		blocking:   config.Blocking,
	}
}

// Subscribers implements stats.Channel.
func (c *Channel) Subscribers() []chan interface{} {
	c.access.RLock()
	defer c.access.RUnlock()
	return c.subscribers
}

// Subscribe implements stats.Channel.
func (c *Channel) Subscribe() (chan interface{}, error) {
	c.access.Lock()
	defer c.access.Unlock()
	if c.subsLimit > 0 && len(c.subscribers) >= c.subsLimit {
		return nil, newError("Number of subscribers has reached limit")
	}
	subscriber := make(chan interface{}, c.bufferSize)
	c.subscribers = append(c.subscribers, subscriber)
	return subscriber, nil
}

// Unsubscribe implements stats.Channel.
func (c *Channel) Unsubscribe(subscriber chan interface{}) error {
	c.access.Lock()
	defer c.access.Unlock()
	for i, s := range c.subscribers {
		if s == subscriber {
			// Copy to new memory block to prevent modifying original data
			subscribers := make([]chan interface{}, len(c.subscribers)-1)
			copy(subscribers[:i], c.subscribers[:i])
			copy(subscribers[i:], c.subscribers[i+1:])
			c.subscribers = subscribers
		}
	}
	return nil
}

// Publish implements stats.Channel.
func (c *Channel) Publish(ctx context.Context, msg interface{}) {
	select { // Early exit if channel closed
	case <-c.closed:
		return
	default:
		pub := channelMessage{context: ctx, message: msg}
		if c.blocking {
			pub.publish(c.channel)
		} else {
			pub.publishNonBlocking(c.channel)
		}
	}
}

// Running returns whether the channel is running.
func (c *Channel) Running() bool {
	select {
	case <-c.closed: // Channel closed
	default: // Channel running or not initialized
		if c.closed != nil { // Channel initialized
			return true
		}
	}
	return false
}

// Start implements common.Runnable.
func (c *Channel) Start() error {
	c.access.Lock()
	defer c.access.Unlock()
	if !c.Running() {
		c.closed = make(chan struct{}) // Reset close signal
		go func() {
			for {
				select {
				case pub := <-c.channel: // Published message received
					for _, sub := range c.Subscribers() { // Concurrency-safe subscribers retrievement
						if c.blocking {
							pub.broadcast(sub)
						} else {
							pub.broadcastNonBlocking(sub)
						}
					}
				case <-c.closed: // Channel closed
					for _, sub := range c.Subscribers() { // Remove all subscribers
						common.Must(c.Unsubscribe(sub))
						close(sub)
					}
					return
				}
			}
		}()
	}
	return nil
}

// Close implements common.Closable.
func (c *Channel) Close() error {
	c.access.Lock()
	defer c.access.Unlock()
	if c.Running() {
		close(c.closed) // Send closed signal
	}
	return nil
}

// channelMessage is the published message with guaranteed delivery.
// message is discarded only when the context is early cancelled.
type channelMessage struct {
	context context.Context
	message interface{}
}

func (c channelMessage) publish(publisher chan channelMessage) {
	select {
	case publisher <- c:
	case <-c.context.Done():
	}
}

func (c channelMessage) publishNonBlocking(publisher chan channelMessage) {
	select {
	case publisher <- c:
	default: // Create another goroutine to keep sending message
		go c.publish(publisher)
	}
}

func (c channelMessage) broadcast(subscriber chan interface{}) {
	select {
	case subscriber <- c.message:
	case <-c.context.Done():
	}
}

func (c channelMessage) broadcastNonBlocking(subscriber chan interface{}) {
	select {
	case subscriber <- c.message:
	default: // Create another goroutine to keep sending message
		go c.broadcast(subscriber)
	}
}
