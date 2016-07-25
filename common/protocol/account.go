package protocol

type Account interface {
	Equals(Account) bool
}
