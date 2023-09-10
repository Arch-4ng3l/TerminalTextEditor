package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Arch-4ng3l/TextEditor/editor"
)

func main() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var char = make([]byte, 1)
	normalMode := false
	ed := editor.NewEditor()
	if len(os.Args) > 1 {
		fmt.Println(os.Args)
		ed.OpenFile(os.Args[2])
	}
	ed.RefreshScreen(false)
	for {
		os.Stdin.Read(char)

		if normalMode {
			normalMode = ed.HandleNormalMode(char[0])
		} else {
			normalMode = ed.HandleInputMode(char[0])
		}
	}
}
