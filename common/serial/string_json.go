// +build json

package serial

import (
	"encoding/json"
)

func (this *StringLiteral) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*this = StringLiteral(str)
	return nil
}
