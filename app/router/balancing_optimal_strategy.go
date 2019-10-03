package router

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/task"
	"v2ray.com/core/features/outbound"
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/pipe"
)

// OptimalStrategy pick outbound by net speed
type OptimalStrategy struct {
	timeout  time.Duration
	interval time.Duration
	url      *url.URL
	count    uint32
	score    float64
	obm      outbound.Manager
	tag      string
	tags     []string
	periodic *task.Periodic
}

// NewOptimalStrategy create new strategy
func NewOptimalStrategy(config *OptimalStrategyConfig) *OptimalStrategy {
	s := &OptimalStrategy{}
	if config.Timeout == 0 {
		s.timeout = time.Second * 5
	} else {
		s.timeout = time.Second * time.Duration(config.Timeout)
	}
	if config.Interval == 0 {
		s.interval = time.Second * 60 * 10
	} else {
		s.interval = time.Second * time.Duration(config.Interval)
	}
	if config.URL == "" {
		s.url, _ = url.Parse("https://www.google.com")
	} else {
		var err error
		s.url, err = url.Parse(config.URL)
		if err != nil {
			panic(err)
		}
		if s.url.Scheme != "http" && s.url.Scheme != "https" {
			panic("Only http/https url support")
		}
	}
	if config.Count == 0 {
		s.count = 3
	} else {
		s.count = config.Count
	}
	s.score = 0

	return s
}

// PickOutbound implement BalancingStrategy interface
func (s *OptimalStrategy) PickOutbound(obm outbound.Manager, tags []string) string {
	if len(tags) == 0 {
		panic("0 tags")
	} else if len(tags) == 1 {
		return s.tag
	}

	s.obm = obm
	s.tags = tags

	if s.periodic == nil {
		s.periodic = &task.Periodic{
			Interval: s.interval,
			Execute:  s.run,
		}
		s.periodic.Start()
		s.tag = s.tags[0]
		return s.tag
	}

	return s.tag
}

// periodic execute function
func (s *OptimalStrategy) run() error {
	s.score = 0

	for _, tag := range s.tags {
		scores := make([]float64, 0, s.count)
		go s.testOutboud(tag, scores)
	}

	return nil
}

// Test outbound's network state with multi-round
func (s *OptimalStrategy) testOutboud(tag string, scores []float64) {
	// calculate average score and end test round
	if len(scores) >= int(s.count) {
		var minScore float64 = float64(math.MaxInt64)
		var maxScore float64 = float64(math.MinInt64)
		var sumScore float64
		var score float64

		for _, score := range scores {
			if score < minScore {
				minScore = score
			}
			if score > maxScore {
				maxScore = score
			}
			sumScore += score
		}
		if len(scores) < 3 {
			score = sumScore / float64(len(scores))
		} else {
			score = (sumScore - minScore - maxScore) / float64(s.count-2)
		}
		newError(fmt.Sprintf("Balance OptimalStrategy get %s's score: %.2f", tag, score)).AtDebug().WriteToLog()

		if s.score < score {
			s.score = score
			s.tag = tag
			newError(fmt.Sprintf("Balance OptimalStrategy now pick detour [%s](score: %.2f) from %s", s.tag, s.score, s.tags)).AtInfo().WriteToLog()
		}
		return
	}
	// test outbound by fetch url
	oh := s.obm.GetHandler(tag)
	if oh == nil {
		newError("Wrong OptimalStrategy tag").AtError().WriteToLog()
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				netDestination, err := net.ParseDestination(fmt.Sprintf("%s:%s", network, addr))
				if err != nil {
					return nil, err
				}

				uplinkReader, uplinkWriter := pipe.New()
				downlinkReader, downlinkWriter := pipe.New()
				ctx = session.ContextWithOutbound(
					ctx,
					&session.Outbound{
						Target: netDestination,
					})
				go oh.Dispatch(ctx, &transport.Link{Reader: uplinkReader, Writer: downlinkWriter})

				return net.NewConnection(net.ConnectionInputMulti(uplinkWriter), net.ConnectionOutputMulti(downlinkReader)), nil
			},
			MaxConnsPerHost: 1,
			MaxIdleConns:    1,
		},
		Timeout: s.timeout,
	}
	startAt := time.Now()
	// send http request though this outbound
	req, _ := http.NewRequest("GET", s.url.String(), nil)
	resp, err := client.Do(req)
	// use http response speed or time(no http content) as score
	score := 0.0
	if err != nil {
		newError(err).AtError().WriteToLog()
	} else {
		contentSize := 0
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			contentSize += len(scanner.Bytes())
		}
		if contentSize != 0 {
			score = float64(contentSize) / (float64(time.Now().UnixNano()-startAt.UnixNano()) / float64(time.Second))
		} else {
			// assert http header's Byte size is 100B
			score = 100 / (float64(time.Now().UnixNano()-startAt.UnixNano()) / float64(time.Second))
		}
	}
	// next test round
	client.CloseIdleConnections()
	s.testOutboud(
		tag,
		append(scores, score),
	)
}
