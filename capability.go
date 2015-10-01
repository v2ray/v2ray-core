package core

type Capability byte

const (
	TCPConnection = Capability(0x01)
	UDPConnection = Capability(0x02)
)

type Capabilities interface {
	HasCapability(cap Capability) bool
	AddCapability(cap Capability)
}

type listCapabilities struct {
	data []Capability
}

func NewCapabilities() Capabilities {
	return &listCapabilities{
		data: make([]Capability, 0, 16),
	}
}

func (c *listCapabilities) HasCapability(cap Capability) bool {
	for _, v := range c.data {
		if v == cap {
			return true
		}
	}
	return false
}

func (c *listCapabilities) AddCapability(cap Capability) {
	c.data = append(c.data, cap)
}
