package dispatcher

import (
	"v2ray.com/core"
	"v2ray.com/core/common/protocol/bittorrent"
	"v2ray.com/core/common/protocol/http"
	"v2ray.com/core/common/protocol/tls"
)

type SniffResult interface {
	Protocol() string
	Domain() string
}

type protocolSniffer func([]byte) (SniffResult, error)

type Sniffer struct {
	sniffer []protocolSniffer
}

func NewSniffer() *Sniffer {
	return &Sniffer{
		sniffer: []protocolSniffer{
			func(b []byte) (SniffResult, error) { return http.SniffHTTP(b) },
			func(b []byte) (SniffResult, error) { return tls.SniffTLS(b) },
			func(b []byte) (SniffResult, error) { return bittorrent.SniffBittorrent(b) },
		},
	}
}

var errUnknownContent = newError("unknown content")

func (s *Sniffer) Sniff(payload []byte) (SniffResult, error) {
	var pendingSniffer []protocolSniffer
	for _, s := range s.sniffer {
		result, err := s(payload)
		if err == core.ErrNoClue {
			pendingSniffer = append(pendingSniffer, s)
			continue
		}

		if err == nil && result != nil {
			return result, nil
		}
	}

	if len(pendingSniffer) > 0 {
		s.sniffer = pendingSniffer
		return nil, core.ErrNoClue
	}

	return nil, errUnknownContent
}
