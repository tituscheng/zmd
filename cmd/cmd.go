package cmd

type Command struct {
	Name        string
	Description string
	Func        func()
}

type CommandManager struct {
	List []Command
}

func (cm CommandManager) Has(commandName string) bool {
	for _, command := range cm.List {
		if command.Name == commandName {
			return true
		}
	}
	return false
}

func (cm CommandManager) Show() {
	// Waiting for implementation
}
