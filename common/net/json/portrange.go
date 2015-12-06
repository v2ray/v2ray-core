package json

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

var (
	InvalidPortRange = errors.New("Invalid port range.")
)

type PortRange struct {
	from v2net.Port
	to   v2net.Port
}

func (this *PortRange) From() v2net.Port {
	return this.from
}

func (this *PortRange) To() v2net.Port {
	return this.to
}

func (this *PortRange) UnmarshalJSON(data []byte) error {
	var maybeint int
	err := json.Unmarshal(data, &maybeint)
	if err == nil {
		if maybeint <= 0 || maybeint >= 65535 {
			log.Error("Invalid port [%s]", string(data))
			return InvalidPortRange
		}
		this.from = v2net.Port(maybeint)
		this.to = v2net.Port(maybeint)
		return nil
	}

	var maybestring string
	err = json.Unmarshal(data, &maybestring)
	if err == nil {
		pair := strings.SplitN(maybestring, "-", 2)
		if len(pair) == 1 {
			value, err := strconv.Atoi(pair[0])
			if err != nil || value <= 0 || value >= 65535 {
				log.Error("Invalid from port %s", pair[0])
				return InvalidPortRange
			}
			this.from = v2net.Port(value)
			this.to = v2net.Port(value)
			return nil
		} else if len(pair) == 2 {
			from, err := strconv.Atoi(pair[0])
			if err != nil || from <= 0 || from >= 65535 {
				log.Error("Invalid from port %s", pair[0])
				return InvalidPortRange
			}
			this.from = v2net.Port(from)

			to, err := strconv.Atoi(pair[1])
			if err != nil || to <= 0 || to >= 65535 {
				log.Error("Invalid to port %s", pair[1])
				return InvalidPortRange
			}
			this.to = v2net.Port(to)

			if this.from > this.to {
				log.Error("Invalid port range %d -> %d", this.from, this.to)
				return InvalidPortRange
			}
			return nil
		}
	}

	return InvalidPortRange
}
