package cmd

import "os/exec"

func newEditorCmd(editor, file string) *exec.Cmd {
	return exec.Command(editor, file)
}
