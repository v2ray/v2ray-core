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

func TestStatsChannel(t *testing.T) {
	// At most 2 subscribers could be registered
	c := NewChannel(&ChannelConfig{SubscriberLimit: 2, Blocking: true})

	a, err := stats.SubscribeRunnableChannel(c)
	common.Must(err)
	if !c.Running() {
		t.Fatal("unexpected failure in running channel after first subscription")
	}

	b, err := c.Subscribe()
	common.Must(err)

	// Test that third subscriber is forbidden
	_, err = c.Subscribe()
	if err == nil {
		t.Fatal("unexpected successful subscription")
	}
	t.Log("expected error: ", err)

	stopCh := make(chan struct{})
	errCh := make(chan string)

	go func() {
		c.Publish(context.Background(), 1)
		c.Publish(context.Background(), 2)
		c.Publish(context.Background(), "3")
		c.Publish(context.Background(), []int{4})
		stopCh <- struct{}{}
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
		stopCh <- struct{}{}
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
		stopCh <- struct{}{}
	}()

	timeout := time.After(2 * time.Second)
	for i := 0; i < 3; i++ {
		select {
		case <-timeout:
			t.Fatal("Test timeout after 2s")
		case e := <-errCh:
			t.Fatal(e)
		case <-stopCh:
		}
	}

	// Test the unsubscription of channel
	common.Must(c.Unsubscribe(b))

	// Test the last subscriber will close channel with `UnsubscribeClosableChannel`
	common.Must(stats.UnsubscribeClosableChannel(c, a))
	if c.Running() {
		t.Fatal("unexpected running channel after unsubscribing the last subscriber")
	}
}

func TestStatsChannelUnsubcribe(t *testing.T) {
	c := NewChannel(&ChannelConfig{Blocking: true})
	common.Must(c.Start())
	defer c.Close()

	a, err := c.Subscribe()
	common.Must(err)
	defer c.Unsubscribe(a)

	b, err := c.Subscribe()
	common.Must(err)

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

	go func() { // Blocking publish
		c.Publish(context.Background(), 1)
		<-pauseCh // Wait for `b` goroutine to resume sending message
		c.Publish(context.Background(), 2)
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
		// Unsubscribe `b` while publishing is paused
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
		// Resume publishing progress
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

func TestStatsChannelBlocking(t *testing.T) {
	// Do not use buffer so as to create blocking scenario
	c := NewChannel(&ChannelConfig{BufferSize: 0, Blocking: true})
	common.Must(c.Start())
	defer c.Close()

	a, err := c.Subscribe()
	common.Must(err)
	defer c.Unsubscribe(a)

	pauseCh := make(chan struct{})
	stopCh := make(chan struct{})
	errCh := make(chan string)

	ctx, cancel := context.WithCancel(context.Background())

	// Test blocking channel publishing
	go func() {
		// Dummy messsage with no subscriber receiving, will block broadcasting goroutine
		c.Publish(context.Background(), nil)

		<-pauseCh

		// Publishing should be blocked here, for last message was not cleared and buffer was full
		c.Publish(context.Background(), nil)

		pauseCh <- struct{}{}

		// Publishing should still be blocked here
		c.Publish(ctx, nil)

		// Check publishing is done because context is canceled
		select {
		case <-ctx.Done():
			if ctx.Err() != context.Canceled {
				errCh <- fmt.Sprint("unexpected error: ", ctx.Err())
			}
		default:
			errCh <- "unexpected non-blocked publishing"
		}
		close(stopCh)
	}()

	go func() {
		pauseCh <- struct{}{}

		select {
		case <-pauseCh:
			errCh <- "unexpected non-blocked publishing"
		case <-time.After(100 * time.Millisecond):
		}

		// Receive first published message
		<-a

		select {
		case <-pauseCh:
		case <-time.After(100 * time.Millisecond):
			errCh <- "unexpected blocking publishing"
		}

		// Manually cancel the context to end publishing
		cancel()
	}()

	select {
	case <-time.After(2 * time.Second):
		t.Fatal("Test timeout after 2s")
	case e := <-errCh:
		t.Fatal(e)
	case <-stopCh:
	}
}

func TestStatsChannelNonBlocking(t *testing.T) {
	// Do not use buffer so as to create blocking scenario
	c := NewChannel(&ChannelConfig{BufferSize: 0, Blocking: false})
	common.Must(c.Start())
	defer c.Close()

	a, err := c.Subscribe()
	common.Must(err)
	defer c.Unsubscribe(a)

	pauseCh := make(chan struct{})
	stopCh := make(chan struct{})
	errCh := make(chan string)

	ctx, cancel := context.WithCancel(context.Background())

	// Test blocking channel publishing
	go func() {
		c.Publish(context.Background(), nil)
		c.Publish(context.Background(), nil)
		pauseCh <- struct{}{}
		<-pauseCh
		c.Publish(ctx, nil)
		c.Publish(ctx, nil)
		// Check publishing is done because context is canceled
		select {
		case <-ctx.Done():
			if ctx.Err() != context.Canceled {
				errCh <- fmt.Sprint("unexpected error: ", ctx.Err())
			}
		case <-time.After(100 * time.Millisecond):
			errCh <- "unexpected non-cancelled publishing"
		}
	}()

	go func() {
		// Check publishing won't block even if there is no subscriber receiving message
		select {
		case <-pauseCh:
		case <-time.After(100 * time.Millisecond):
			errCh <- "unexpected blocking publishing"
		}

		// Receive first and second published message
		<-a
		<-a

		pauseCh <- struct{}{}

		// Manually cancel the context to end publishing
		cancel()

		// Check third and forth published message is cancelled and cannot receive
		<-time.After(100 * time.Millisecond)
		select {
		case <-a:
			errCh <- "unexpected non-cancelled publishing"
		default:
		}
		select {
		case <-a:
			errCh <- "unexpected non-cancelled publishing"
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
	// Do not use buffer so as to create blocking scenario
	c := NewChannel(&ChannelConfig{BufferSize: 0, Blocking: true})
	common.Must(c.Start())
	defer c.Close()

	a, err := c.Subscribe()
	common.Must(err)
	defer c.Unsubscribe(a)

	b, err := c.Subscribe()
	common.Must(err)
	defer c.Unsubscribe(b)

	stopCh := make(chan struct{})
	errCh := make(chan string)

	go func() { // Blocking publish
		c.Publish(context.Background(), 1)
		c.Publish(context.Background(), 2)
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
		// Block `b` for a time so as to ensure source channel is trying to send message to `b`.
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
		// Because in a new round of messaging, `b` has been unsubscribed.
		select {
		case v, ok := <-b:
			if ok {
				errCh <- fmt.Sprint("unexpected receiving: ", v)
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
