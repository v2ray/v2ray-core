package json

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/v2ray/v2ray-core/common/log"
)

var (
	InvalidPortRange = errors.New("Invalid port range.")
)

type PortRange struct {
	from uint16
	to   uint16
}

func (this *PortRange) From() uint16 {
	return this.from
}

func (this *PortRange) To() uint16 {
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
		this.from = uint16(maybeint)
		this.to = uint16(maybeint)
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
			this.from = uint16(value)
			this.to = uint16(value)
			return nil
		} else if len(pair) == 2 {
			from, err := strconv.Atoi(pair[0])
			if err != nil || from <= 0 || from >= 65535 {
				log.Error("Invalid from port %s", pair[0])
				return InvalidPortRange
			}
			this.from = uint16(from)

			to, err := strconv.Atoi(pair[1])
			if err != nil || to <= 0 || to >= 65535 {
				log.Error("Invalid to port %s", pair[1])
				return InvalidPortRange
			}
			this.to = uint16(to)

			if this.from > this.to {
				log.Error("Invalid port range %d -> %d", this.from, this.to)
				return InvalidPortRange
			}
			return nil
		}
	}

	return InvalidPortRange
}
