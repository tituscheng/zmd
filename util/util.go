package util

import (
	"os"
	"os/exec"
)

func ExecCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func HasCommand(path string) bool {
	_, err1 := exec.LookPath(path)
	return err1 == nil
}

func HasFlag(flag string) bool {
	flagMap := make(map[string]bool)
	flagMap["--"+flag] = true
	flagMap["-"+flag] = true

	for _, arg := range os.Args {
		if _, exists := flagMap[arg]; exists {
			return true
		}
	}
	return false
}
