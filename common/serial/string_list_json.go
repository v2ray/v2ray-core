// +build json

package serial

import (
	"encoding/json"
	"errors"
	"strings"
)

func (this *StringTList) UnmarshalJSON(data []byte) error {
	var strarray []string
	if err := json.Unmarshal(data, &strarray); err == nil {
		*this = *NewStringTList(strarray)
		return nil
	}

	var rawstr string
	if err := json.Unmarshal(data, &rawstr); err == nil {
		strlist := strings.Split(rawstr, ",")
		*this = *NewStringTList(strlist)
		return nil
	}
	return errors.New("Unknown format of a string list: " + string(data))
}
