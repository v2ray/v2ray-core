package policy

import (
	"time"

	"github.com/golang/protobuf/proto"
)

func (s *Second) Duration() time.Duration {
	return time.Second * time.Duration(s.Value)
}

func (p *Policy) OverrideWith(another *Policy) {
	proto.Merge(p, another)
}
