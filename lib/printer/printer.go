package printer

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type BarPrinter struct {
	vals                      []int
	color                     *color.Color
	writer                    *lineWriter
	width, height, top, sides int
	shouldClear, coloured     bool
	stop                      chan struct{}
}

func NewBarPrinter(width, height, top, sides int) *BarPrinter {
	writer := &lineWriter{Writer: bufio.NewWriter(os.Stdout)}
	colorer := color.New(color.BgBlue)
	color.Output = writer

	return &BarPrinter{
		vals:   make([]int, 0),
		width:  width - (sides * 2) - 2,
		height: height,
		top:    top,
		sides:  sides,
		color:  colorer,
		writer: writer,
		stop:   make(chan struct{}),
	}
}

func (b *BarPrinter) Feed(c chan int) {
	for {
		select {
		case n := <-c:
			b.AdvanceN(n)
		case <-b.stop:
		}
	}
}

func (b *BarPrinter) AdvanceN(count int) {
	b.vals = append(b.vals, count)
	b.render()
}

func (b *BarPrinter) Stop() {
	b.stop <- struct{}{}
}

func (b *BarPrinter) printer(coloured bool) *color.Color {
	if b.coloured && !coloured {
		b.color.DisableColor()
	} else if !b.coloured && coloured {
		b.color.EnableColor()
	}
	b.coloured = coloured
	return b.color
}

func (b *BarPrinter) printAllLines() {
	for i := 0; i < b.height; i++ {
		b.printLine(i)
	}
}

func (b *BarPrinter) printLine(i int) {
	b.printer(false).Printf("%s|", strings.Repeat(" ", b.sides))
	for j := 0; j < b.width; j++ {
		b.printChar(i, j)
	}
	b.printer(false).Printf("|%s", strings.Repeat(" ", b.sides))
	b.writer.commitLine()
}

func (b *BarPrinter) printChar(row, col int) {
	on := b.isOn(row, col)
	toPrint := " "
	if !on && b.isOn(row+1, col) {
		toPrint = "_"
	}
	b.printer(on).Print(toPrint)
}

func (b *BarPrinter) isOn(row, col int) bool {
	diff := b.width - len(b.vals)
	return (col >= diff) && (row > (b.height - b.vals[col-diff]))
}

func (b *BarPrinter) render() {
	b.printBlock(b.top)
	b.printRule("_")
	b.printAllLines()
	b.printRule("_")
	b.printBlock(b.top)
	b.writer.Flush()
}

func (b *BarPrinter) printBlock(n int) {
	for i := 0; i < n; i++ {
		b.printer(false).Print("")
		b.writer.commitLine()
	}
}

func (b *BarPrinter) printRule(c string) {
	b.printer(false).Printf("%s%s%s", strings.Repeat(" ", b.sides+1), strings.Repeat(c, b.width), strings.Repeat(" ", b.sides+1))
	b.writer.commitLine()
}

type lineWriter struct {
	*bufio.Writer
	current, last int
}

func (l *lineWriter) commitLine() {
	if l.current == 0 {
		fmt.Fprintf(l, "\033[%dA", l.last)
	}

	l.current++
	fmt.Fprintln(l, "")
}

func (l *lineWriter) Flush() error {
	if err := l.Writer.Flush(); err != nil {
		return err
	}
	l.last = l.current
	l.current = 0
	return nil
}
