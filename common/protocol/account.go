package protocol

type Account interface {
	Equals(Account) bool
}

type AsAccount interface {
	AsAccount() (Account, error)
}

type NewAccountFactory func() AsAccount
