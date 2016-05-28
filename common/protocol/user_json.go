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

	var rawAccount AccountJson
	if err := json.Unmarshal(data, &rawAccount); err != nil {
		return err
	}
	account, err := rawAccount.GetAccount()
	if err != nil {
		return err
	}

	*u = *NewUser(account, UserLevel(rawUserValue.LevelByte), rawUserValue.EmailString)

	return nil
}
