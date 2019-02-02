package command_test

import (
	"context"
	"testing"

	"v2ray.com/core/app/stats"
	. "v2ray.com/core/app/stats/command"
	"v2ray.com/core/common"
	. "v2ray.com/ext/assert"
)

func TestGetStats(t *testing.T) {
	assert := With(t)

	m, err := stats.NewManager(context.Background(), &stats.Config{})
	common.Must(err)

	sc, err := m.RegisterCounter("test_counter")
	common.Must(err)

	sc.Set(1)

	s := NewStatsServer(m)

	testCases := []struct {
		name  string
		reset bool
		value int64
		err   bool
	}{
		{
			name: "counterNotExist",
			err:  true,
		},
		{
			name:  "test_counter",
			reset: true,
			value: 1,
		},
		{
			name:  "test_counter",
			value: 0,
		},
	}
	for _, tc := range testCases {
		resp, err := s.GetStats(context.Background(), &GetStatsRequest{
			Name:   tc.name,
			Reset_: tc.reset,
		})
		if tc.err {
			assert(err, IsNotNil)
		} else {
			common.Must(err)
			assert(resp.Stat.Name, Equals, tc.name)
			assert(resp.Stat.Value, Equals, tc.value)
		}
	}
}

func TestQueryStats(t *testing.T) {
	assert := With(t)

	m, err := stats.NewManager(context.Background(), &stats.Config{})
	common.Must(err)

	sc1, err := m.RegisterCounter("test_counter")
	common.Must(err)
	sc1.Set(1)

	sc2, err := m.RegisterCounter("test_counter_2")
	common.Must(err)
	sc2.Set(2)

	sc3, err := m.RegisterCounter("test_counter_3")
	common.Must(err)
	sc3.Set(3)

	s := NewStatsServer(m)
	resp, err := s.QueryStats(context.Background(), &QueryStatsRequest{
		Pattern: "counter_",
	})
	common.Must(err)
	assert(len(resp.Stat), Equals, 2)

	v2 := false
	v3 := false
	for _, sc := range resp.Stat {
		switch sc.Name {
		case "test_counter_2":
			assert(sc.Value, Equals, int64(2))
			v2 = true
		case "test_counter_3":
			assert(sc.Value, Equals, int64(3))
			v3 = true
		default:
			t.Error("unexpected stat name: ", sc.Name)
			t.Fail()
		}
	}
	assert(v2, IsTrue)
	assert(v3, IsTrue)
}
