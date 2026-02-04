package internal

import (
	"fmt"
	"os/exec"
	"runtime"
)

func CopyToClipboard(text string) error {
	copyCmd, err := getClipboardCmd()
	if err != nil {
		return err
	}

	pipe, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}

	if _, err := pipe.Write([]byte(text)); err != nil {
		return err
	}

	if err := pipe.Close(); err != nil {
		return err
	}

	return copyCmd.Wait()
}

func getClipboardCmd() (*exec.Cmd, error) {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("pbcopy"), nil
	case "linux":
		if _, err := exec.LookPath("xclip"); err == nil {
			return exec.Command("xclip", "-selection", "clipboard"), nil
		}
		if _, err := exec.LookPath("xsel"); err == nil {
			return exec.Command("xsel", "--clipboard", "--input"), nil
		}
		return nil, fmt.Errorf("no clipboard tool found. Install xclip or xsel")
	default:
		return exec.Command("clip"), nil
	}
}
