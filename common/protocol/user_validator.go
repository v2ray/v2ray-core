package protocol

type UserValidator interface {
	Add(user *User) error
	Get(timeHash []byte) (*User, Timestamp, bool)
}
