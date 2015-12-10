package json

import (
	"encoding/json"
	"errors"
	"strings"
)

func UnmarshalStringList(data []byte) ([]string, error) {
	var strarray []string
	if err := json.Unmarshal(data, &strarray); err == nil {
		return strarray, nil
	}

	var rawstr string
	if err := json.Unmarshal(data, &rawstr); err == nil {
		strlist := strings.Split(rawstr, ",")
		return strlist, nil
	}
	return nil, errors.New("Unknown format of a string list: " + string(data))
}
