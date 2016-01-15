// +build json

package net

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

func (this *PortRange) UnmarshalJSON(data []byte) error {
	var maybeint int
	err := json.Unmarshal(data, &maybeint)
	if err == nil {
		if maybeint <= 0 || maybeint >= 65535 {
			log.Error("Invalid port [%s]", string(data))
			return InvalidPortRange
		}
		this.From = Port(maybeint)
		this.To = Port(maybeint)
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
			this.From = Port(value)
			this.To = Port(value)
			return nil
		} else if len(pair) == 2 {
			from, err := strconv.Atoi(pair[0])
			if err != nil || from <= 0 || from >= 65535 {
				log.Error("Invalid from port %s", pair[0])
				return InvalidPortRange
			}
			this.From = Port(from)

			to, err := strconv.Atoi(pair[1])
			if err != nil || to <= 0 || to >= 65535 {
				log.Error("Invalid to port %s", pair[1])
				return InvalidPortRange
			}
			this.To = Port(to)

			if this.From > this.To {
				log.Error("Invalid port range %d -> %d", this.From, this.To)
				return InvalidPortRange
			}
			return nil
		}
	}

	return InvalidPortRange
}
