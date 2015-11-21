package testing

type PortRange struct {
	FromValue uint16
	ToValue   uint16
}

func (this *PortRange) From() uint16 {
	return this.FromValue
}

func (this *PortRange) To() uint16 {
	return this.ToValue
}
