package stats_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "v2ray.com/core/app/stats"
	"v2ray.com/core/common"
	"v2ray.com/core/features/stats"
)

func TestInterface(t *testing.T) {
	_ = (stats.Manager)(new(Manager))
}

func TestStatsCounter(t *testing.T) {
	raw, err := common.CreateObject(context.Background(), &Config{})
	common.Must(err)

	m := raw.(stats.Manager)
	c, err := m.RegisterCounter("test.counter")
	common.Must(err)

	if v := c.Add(1); v != 1 {
		t.Fatal("unpexcted Add(1) return: ", v, ", wanted ", 1)
	}

	if v := c.Set(0); v != 1 {
		t.Fatal("unexpected Set(0) return: ", v, ", wanted ", 1)
	}

	if v := c.Value(); v != 0 {
		t.Fatal("unexpected Value() return: ", v, ", wanted ", 0)
	}
}

func TestStatsChannel(t *testing.T) {
	raw, err := common.CreateObject(context.Background(), &Config{})
	common.Must(err)

	m := raw.(stats.Manager)
	c, err := m.RegisterChannel("test.channel")
	common.Must(err)

	source := c.Channel()
	a := c.Subscribe()
	b := c.Subscribe()
	defer c.Unsubscribe(a)
	defer c.Unsubscribe(b)

	stopCh := make(chan struct{})
	errCh := make(chan string)

	go func() {
		source <- 1
		source <- 2
		source <- "3"
		source <- []int{4}
		source <- nil // Dummy messsage with no subscriber receiving
		select {
		case source <- nil: // Source should be blocked here, for last message was not cleared
			errCh <- fmt.Sprint("unexpected non-blocked source")
		default:
			close(stopCh)
		}
	}()

	go func() {
		if v, ok := (<-a).(int); !ok || v != 1 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 1)
		}
		if v, ok := (<-a).(int); !ok || v != 2 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 2)
		}
		if v, ok := (<-a).(string); !ok || v != "3" {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", "3")
		}
		if v, ok := (<-a).([]int); !ok || v[0] != 4 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", []int{4})
		}
	}()

	go func() {
		if v, ok := (<-b).(int); !ok || v != 1 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 1)
		}
		if v, ok := (<-b).(int); !ok || v != 2 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 2)
		}
		if v, ok := (<-b).(string); !ok || v != "3" {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", "3")
		}
		if v, ok := (<-b).([]int); !ok || v[0] != 4 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", []int{4})
		}
	}()

	select {
	case <-time.After(2 * time.Second):
		t.Fatal("Test timeout after 2s")
	case e := <-errCh:
		t.Fatal(e)
	case <-stopCh:
	}
}

func TestStatsChannelUnsubcribe(t *testing.T) {
	raw, err := common.CreateObject(context.Background(), &Config{})
	common.Must(err)

	m := raw.(stats.Manager)
	c, err := m.RegisterChannel("test.channel")
	common.Must(err)

	source := c.Channel()
	a := c.Subscribe()
	b := c.Subscribe()
	defer c.Unsubscribe(a)

	pauseCh := make(chan struct{})
	stopCh := make(chan struct{})
	errCh := make(chan string)

	{
		var aSet, bSet bool
		for _, s := range c.Subscribers() {
			if s == a {
				aSet = true
			}
			if s == b {
				bSet = true
			}
		}
		if !(aSet && bSet) {
			t.Fatal("unexpected subscribers: ", c.Subscribers())
		}
	}

	go func() {
		source <- 1
		<-pauseCh // Wait for `b` goroutine to resume sending message
		source <- 2
	}()

	go func() {
		if v, ok := (<-a).(int); !ok || v != 1 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 1)
		}
		if v, ok := (<-a).(int); !ok || v != 2 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 2)
		}
	}()

	go func() {
		if v, ok := (<-b).(int); !ok || v != 1 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 1)
		}
		// Unsubscribe `b` while `source`'s messaging is paused
		c.Unsubscribe(b)
		{ // Test `b` is not in subscribers
			var aSet, bSet bool
			for _, s := range c.Subscribers() {
				if s == a {
					aSet = true
				}
				if s == b {
					bSet = true
				}
			}
			if !(aSet && !bSet) {
				errCh <- fmt.Sprint("unexpected subscribers: ", c.Subscribers())
			}
		}
		// Resume `source`'s progress
		close(pauseCh)
		// Test `b` is neither closed nor able to receive any data
		select {
		case v, ok := <-b:
			if ok {
				errCh <- fmt.Sprint("unexpected data received: ", v)
			} else {
				errCh <- fmt.Sprint("unexpected closed channel: ", b)
			}
		default:
		}
		close(stopCh)
	}()

	select {
	case <-time.After(2 * time.Second):
		t.Fatal("Test timeout after 2s")
	case e := <-errCh:
		t.Fatal(e)
	case <-stopCh:
	}
}

func TestStatsChannelTimeout(t *testing.T) {
	raw, err := common.CreateObject(context.Background(), &Config{})
	common.Must(err)

	m := raw.(stats.Manager)
	c, err := m.RegisterChannel("test.channel")
	common.Must(err)

	source := c.Channel()
	a := c.Subscribe()
	b := c.Subscribe()
	defer c.Unsubscribe(a)
	defer c.Unsubscribe(b)

	stopCh := make(chan struct{})
	errCh := make(chan string)

	go func() {
		source <- 1
		source <- 2
	}()

	go func() {
		if v, ok := (<-a).(int); !ok || v != 1 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 1)
		}
		if v, ok := (<-a).(int); !ok || v != 2 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 2)
		}
		{ // Test `b` is still in subscribers yet (because `a` receives 2 first)
			var aSet, bSet bool
			for _, s := range c.Subscribers() {
				if s == a {
					aSet = true
				}
				if s == b {
					bSet = true
				}
			}
			if !(aSet && bSet) {
				errCh <- fmt.Sprint("unexpected subscribers: ", c.Subscribers())
			}
		}
	}()

	go func() {
		if v, ok := (<-b).(int); !ok || v != 1 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 1)
		}
		// Block `b` channel for a time longer than `source`'s timeout
		<-time.After(150 * time.Millisecond)
		{ // Test `b` has been unsubscribed by source
			var aSet, bSet bool
			for _, s := range c.Subscribers() {
				if s == a {
					aSet = true
				}
				if s == b {
					bSet = true
				}
			}
			if !(aSet && !bSet) {
				errCh <- fmt.Sprint("unexpected subscribers: ", c.Subscribers())
			}
		}
		select { // Test `b` has been closed by source
		case v, ok := <-b:
			if ok {
				errCh <- fmt.Sprint("unexpected data received: ", v)
			}
		default:
		}
		close(stopCh)
	}()

	select {
	case <-time.After(2 * time.Second):
		t.Fatal("Test timeout after 2s")
	case e := <-errCh:
		t.Fatal(e)
	case <-stopCh:
	}
}

func TestStatsChannelConcurrency(t *testing.T) {
	raw, err := common.CreateObject(context.Background(), &Config{})
	common.Must(err)

	m := raw.(stats.Manager)
	c, err := m.RegisterChannel("test.channel")
	common.Must(err)

	source := c.Channel()
	a := c.Subscribe()
	b := c.Subscribe()
	defer c.Unsubscribe(a)

	stopCh := make(chan struct{})
	errCh := make(chan string)

	go func() {
		source <- 1
		source <- 2
	}()

	go func() {
		if v, ok := (<-a).(int); !ok || v != 1 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 1)
		}
		if v, ok := (<-a).(int); !ok || v != 2 {
			errCh <- fmt.Sprint("unexpected receiving: ", v, ", wanted ", 2)
		}
	}()

	go func() {
		// Block `b` for a time shorter than `source`'s timeout
		// So as to ensure source channel is trying to send message to `b`.
		<-time.After(25 * time.Millisecond)
		// This causes concurrency scenario: unsubscribe `b` while trying to send message to it
		c.Unsubscribe(b)
		// Test `b` is not closed and can still receive data 1:
		// Because unsubscribe won't affect the ongoing process of sending message.
		select {
		case v, ok := <-b:
			if v1, ok1 := v.(int); !(ok && ok1 && v1 == 1) {
				errCh <- fmt.Sprint("unexpected failure in receiving data: ", 1)
			}
		default:
			errCh <- fmt.Sprint("unexpected block from receiving data: ", 1)
		}
		// Test `b` is not closed but cannot receive data 2:
		// Becuase in a new round of messaging, `b` has been unsubscribed.
		select {
		case v, ok := <-b:
			if ok {
				errCh <- fmt.Sprint("unexpected receving: ", v)
			} else {
				errCh <- fmt.Sprint("unexpected closing of channel")
			}
		default:
		}
		close(stopCh)
	}()

	select {
	case <-time.After(2 * time.Second):
		t.Fatal("Test timeout after 2s")
	case e := <-errCh:
		t.Fatal(e)
	case <-stopCh:
	}
}
