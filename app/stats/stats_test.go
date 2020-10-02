package stats_test

import (
	"context"
	"testing"
	"time"

	. "v2ray.com/core/app/stats"
	"v2ray.com/core/common"
	"v2ray.com/core/features/stats"
)

func TestInterface(t *testing.T) {
	_ = (stats.Manager)(new(Manager))
}

func TestStatsChannelRunnable(t *testing.T) {
	raw, err := common.CreateObject(context.Background(), &Config{})
	common.Must(err)

	m := raw.(stats.Manager)

	ch1, err := m.RegisterChannel("test.channel.1")
	c1 := ch1.(*Channel)
	common.Must(err)

	if c1.Running() {
		t.Fatalf("unexpected running channel: test.channel.%d", 1)
	}

	common.Must(m.Start())

	if !c1.Running() {
		t.Fatalf("unexpected non-running channel: test.channel.%d", 1)
	}

	ch2, err := m.RegisterChannel("test.channel.2")
	c2 := ch2.(*Channel)
	common.Must(err)

	if !c2.Running() {
		t.Fatalf("unexpected non-running channel: test.channel.%d", 2)
	}

	s1, err := c1.Subscribe()
	common.Must(err)
	common.Must(c1.Close())

	if c1.Running() {
		t.Fatalf("unexpected running channel: test.channel.%d", 1)
	}

	select { // Check all subscribers in closed channel are closed
	case _, ok := <-s1:
		if ok {
			t.Fatalf("unexpected non-closed subscriber in channel: test.channel.%d", 1)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("unexpected non-closed subscriber in channel: test.channel.%d", 1)
	}

	if len(c1.Subscribers()) != 0 { // Check subscribers in closed channel are emptied
		t.Fatalf("unexpected non-empty subscribers in channel: test.channel.%d", 1)
	}

	common.Must(m.Close())

	if c2.Running() {
		t.Fatalf("unexpected running channel: test.channel.%d", 2)
	}

	ch3, err := m.RegisterChannel("test.channel.3")
	c3 := ch3.(*Channel)
	common.Must(err)

	if c3.Running() {
		t.Fatalf("unexpected running channel: test.channel.%d", 3)
	}

	common.Must(c3.Start())
	common.Must(m.UnregisterChannel("test.channel.3"))

	if c3.Running() { // Test that unregistering will close the channel.
		t.Fatalf("unexpected running channel: test.channel.%d", 3)
	}
}
