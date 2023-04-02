package logger

import (
	"fmt"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

type Logger struct {
	text *text.Text
}

func (l *Logger) Info(s string) {
	l.text.Write(fmt.Sprintf("%s [Info]: %s\n", time.Now().Format("15:04:05"), s), text.WriteCellOpts(cell.FgColor(cell.ColorDefault)))

}

func (l *Logger) Error(err error) {

	l.text.Write(fmt.Sprintf("%s [Error]: %s\n", time.Now().Format("15:04:05"), err.Error()), text.WriteCellOpts(cell.FgColor(cell.ColorRed)))
}

func (l *Logger) GetWidget() *text.Text {
	return l.text
}

func NewLogger() *Logger {
	trimmed, err := text.New(text.RollContent(), text.WrapAtRunes())
	if err != nil {
		panic(err)
	}

	l := Logger{
		text: trimmed,
	}
	return &l
}
