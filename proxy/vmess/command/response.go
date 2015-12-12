package command

var (
	cmdCache = make(map[byte]CommandCreator)
)

func RegisterResponseCommand(id byte, cmdFactory CommandCreator) error {
	cmdCache[id] = cmdFactory
	return nil
}

func CreateResponseCommand(id byte) (Command, error) {
	creator, found := cmdCache[id]
	if !found {
		return nil, ErrorNoSuchCommand
	}
	return creator(), nil
}
