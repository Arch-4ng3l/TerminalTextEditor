package editor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Arch-4ng3l/TextEditor/highlighting"
)

const (
	ESC       = 27
	BACKSPACE = 127
)

type Editor struct {
	cursorX, cursorY int
	lines            []string
	h                highlighting.Highlighter
}

func NewEditor() *Editor {
	return &Editor{
		cursorX: 0,
		cursorY: 0,
		lines:   []string{""},
		h:       *highlighting.NewHighlighter(),
	}
}

func (e *Editor) OpenFile(fileName string) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	lines := strings.Split(string(content)+"\x00", "\n")
	for i := range lines {
		lines[i] = strings.ReplaceAll(strings.ReplaceAll(lines[i], "\t", "    "), "\n", " ")
	}
	e.lines = lines
	e.RefreshScreen(false)
}

func (e *Editor) RefreshScreen(b bool) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	limit1, limit2 := 0, len(e.lines)
	moved := false
	if e.cursorY >= 12 {
		limit1 = e.cursorY - 12
		moved = true
	}
	if limit2-e.cursorY <= 12 {
		moved = true
	} else {
		limit2 = e.cursorY + 12
	}

	for _, l := range e.lines[limit1:limit2] {
		fmt.Println(e.h.Highlight(l))
	}

	if b {
		fmt.Print("\n----NormalMode")
	} else {
		fmt.Print("\n----InputMode")
	}
	fmt.Printf(" X: %d Y:%d \t%d\n", e.cursorX, e.cursorY, len(e.lines))

	if !moved {
		fmt.Printf("\033[%d;%dH", e.cursorY+1, e.cursorX+1)
	} else {
		fmt.Printf("\033[%d;%dH", e.cursorY, e.cursorX+1)
	}

}

func (e *Editor) HandleInputMode(ch byte) bool {
	if ch == ESC {
		e.RefreshScreen(true)
		return true
	}

	e.AddChar(ch)
	return false
}

func (e *Editor) AddChar(ch byte) {
	if ch == BACKSPACE {

		if e.cursorX == 0 {
			if e.cursorY > 0 {
				if e.cursorY == len(e.lines)-1 {
					e.lines = e.lines[:e.cursorY]
				} else {
					e.lines = append(e.lines[:e.cursorY], e.lines[e.cursorY+1:]...)
				}
				e.cursorY--
				e.cursorX = len(e.lines[e.cursorY])
			}
			e.RefreshScreen(false)
			return
		}

		e.lines[e.cursorY] = e.lines[e.cursorY][:e.cursorX-1] + e.lines[e.cursorY][e.cursorX:]
		e.cursorX--

	} else if ch == '\n' {

		e.cursorY++
		e.lines = append(e.lines, "")
		e.cursorX = 0

	} else if ch == '\t' {

		e.lines[e.cursorY] = e.lines[e.cursorY][:e.cursorX] + "    " + e.lines[e.cursorY][e.cursorX:]
		e.cursorX += 4

	} else {

		e.lines[e.cursorY] = e.lines[e.cursorY][:e.cursorX] + string(ch) + e.lines[e.cursorY][e.cursorX:]
		e.cursorX++

	}

	e.RefreshScreen(false)
}

func (e *Editor) HandleNormalMode(ch byte) bool {
	switch ch {
	case 'i':
		if e.cursorX > len(e.lines[e.cursorY]) {
			e.cursorX = len(e.lines[e.cursorY])
		}
		e.RefreshScreen(false)
		return false
	case 'a':
		if e.cursorX >= len(e.lines[e.cursorY]) {
			e.cursorX = len(e.lines[e.cursorY])
		}
		e.cursorX++
		e.RefreshScreen(false)
		return false

	case 'h':
		e.cursorX--

	case 'l':
		e.cursorX++

	case 'k':
		if e.cursorY > 0 {
			e.cursorY--
		}

	case 'j':
		if e.cursorY < len(e.lines)-1 {
			e.cursorY++
		}
	}

	e.RefreshScreen(true)
	return true
}
