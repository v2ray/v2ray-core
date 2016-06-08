// +build json

package collect

import (
	"encoding/json"
	"errors"
	"strings"
)

func (this *StringList) UnmarshalJSON(data []byte) error {
	var strarray []string
	if err := json.Unmarshal(data, &strarray); err == nil {
		*this = *NewStringList(strarray)
		return nil
	}

	var rawstr string
	if err := json.Unmarshal(data, &rawstr); err == nil {
		strlist := strings.Split(rawstr, ",")
		*this = *NewStringList(strlist)
		return nil
	}
	return errors.New("Unknown format of a string list: " + string(data))
}
