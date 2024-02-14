package media

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"zmd/util"
)

const (
	MAC_OPEN_COMMAND     = "open"
	WINDOWS_OPEN_COMMAND = "start"

	COMMAND_LINE_MODE = "cmd"
)

var (
	defaultMode = COMMAND_LINE_MODE
)

func playLocalMediaFileByCommandLine(program, mediaFilePath string) {
	if _, err := os.Stat(mediaFilePath); errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("[error] %s -> media file not exists\n", mediaFilePath)
		return
	}

	if !util.HasCommand(program) {
		fmt.Printf("Program '%s' not found\n", program)
		return
	}

	if err := util.ExecCommand(program, []string{mediaFilePath}); err != nil {
		panic(err)
	}
}

func playMediaInMode(mode, program, file string) {
	switch mode {
	case COMMAND_LINE_MODE:
		playLocalMediaFileByCommandLine(program, file)
	default:
		return
	}
}

func playMediaFile(mode, mediaFilePath string) {
	switch runtime.GOOS {
	case "darwin":
		playMediaInMode(mode, MAC_OPEN_COMMAND, mediaFilePath)
	case "windows":
		fmt.Println("windows")
		playMediaInMode(mode, WINDOWS_OPEN_COMMAND, mediaFilePath)
	default:
		return
	}
}

func SetMode(mode string) {
	defaultMode = mode
}

// Play plays a file in default mode unless specified, default mode is COMMAND_LINE, or use SetMode to modify
func Play(file string) {
	playMediaFile(defaultMode, file)
}
