package protocol

type TransferType byte

const (
	TransferTypeStream TransferType = 0
	TransferTypePacket TransferType = 1
)

type AddressType byte

const (
	AddressTypeIPv4   AddressType = 1
	AddressTypeDomain AddressType = 2
	AddressTypeIPv6   AddressType = 3
)
