// +build json

package protocol

import "encoding/json"

func (u *User) UnmarshalJSON(data []byte) error {
	type rawUser struct {
		EmailString string `json:"email"`
		LevelByte   byte   `json:"level"`
	}
	var rawUserValue rawUser
	if err := json.Unmarshal(data, &rawUserValue); err != nil {
		return err
	}

	u.Email = rawUserValue.EmailString
	u.Level = UserLevel(rawUserValue.LevelByte)

	return nil
}
