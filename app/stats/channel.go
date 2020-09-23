// +build !confonly

package stats

import (
	"sync"
	"time"

	"v2ray.com/core/common"
)

// Channel is an implementation of stats.Channel.
type Channel struct {
	channel     chan interface{}
	subscribers []chan interface{}

	// Synchronization components
	access sync.RWMutex
	closed chan struct{}

	// Channel options
	subscriberLimit   int           // Set to 0 as no subscriber limit
	channelBufferSize int           // Set to 0 as no buffering
	broadcastTimeout  time.Duration // Set to 0 as non-blocking immediate timeout
}

// NewChannel creates an instance of Statistics Channel.
func NewChannel(config *ChannelConfig) *Channel {
	return &Channel{
		channel:           make(chan interface{}, config.BufferSize),
		subscriberLimit:   int(config.SubscriberLimit),
		channelBufferSize: int(config.BufferSize),
		broadcastTimeout:  time.Duration(config.BroadcastTimeout+1) * time.Millisecond,
	}
}

// Channel returns the underlying go channel.
func (c *Channel) Channel() chan interface{} {
	c.access.RLock()
	defer c.access.RUnlock()
	return c.channel
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
	if c.subscriberLimit > 0 && len(c.subscribers) >= c.subscriberLimit {
		return nil, newError("Number of subscribers has reached limit")
	}
	subscriber := make(chan interface{}, c.channelBufferSize)
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
func (c *Channel) Publish(message interface{}) {
	select { // Early exit if channel closed
	case <-c.closed:
		return
	default:
	}
	select { // Drop message if not successfully sent
	case c.channel <- message:
	default:
		return
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
				case message := <-c.channel: // Broadcast message
					for _, sub := range c.Subscribers() { // Concurrency-safe subscribers retreivement
						select {
						case sub <- message: // Successfully sent message
						case <-time.After(c.broadcastTimeout): // Remove timeout subscriber
							common.Must(c.Unsubscribe(sub))
							close(sub) // Actively close subscriber as notification
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
