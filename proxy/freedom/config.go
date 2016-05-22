package freedom

type DomainStrategy int

const (
	DomainStrategyAsIs  = DomainStrategy(0)
	DomainStrategyUseIP = DomainStrategy(1)
)

type Config struct {
	DomainStrategy DomainStrategy
}
