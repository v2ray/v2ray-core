package protocol

// Account is an user identity used for authentication.
type Account interface {
	Equals(Account) bool
}

// AsAccount is an object can be converted into account.
type AsAccount interface {
	AsAccount() (Account, error)
}
