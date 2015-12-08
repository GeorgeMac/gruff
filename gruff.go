package gruff

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Normaliser func(f float64) func(float64) int

type Printer struct {
	vals                      []float64
	color                     *color.Color
	writer                    *lineWriter
	width, height, top, sides int
	coloured                  bool
	norm                      Normaliser
	curNorm                   func(float64) int
	stop                      chan struct{}
}

func New(width, height int, opts ...Option) *Printer {
	writer := &lineWriter{Writer: bufio.NewWriter(os.Stdout)}
	colorer := color.New(color.BgBlue)
	color.Output = writer

	norm := func(f float64) func(float64) int {
		return func(f float64) int { return int(math.Floor(f)) }
	}

	printer := &Printer{
		vals:   make([]float64, 0),
		height: height,
		top:    5,
		sides:  10,
		color:  colorer,
		writer: writer,
		norm:   norm,
		stop:   make(chan struct{}),
	}

	for _, opt := range opts {
		opt(printer)
	}

	printer.width = width - (printer.sides * 2) - 2

	return printer
}

func (b *Printer) Feed(c <-chan float64) {
	for {
		select {
		case <-b.stop:
			return
		case f, ok := <-c:
			if !ok {
				close(b.stop)
				return
			}
			b.Advance(f)
		}
	}
}

func (b *Printer) Advance(count float64) {
	b.curNorm = b.norm(count)
	if len(b.vals) > b.width {
		b.vals = append(b.vals[1:], count)
	} else {
		b.vals = append(b.vals, count)
	}
	b.render()
}

func (b *Printer) Stop() {
	b.stop <- struct{}{}
}

func (b *Printer) render() {
	b.printBlock(b.top)
	b.printRule("_")
	b.printAllLines()
	b.printRule("_")
	b.printBlock(b.top)
	b.writer.Flush()
}

func (b *Printer) printer(coloured bool) *color.Color {
	if b.coloured && !coloured {
		b.color.DisableColor()
	} else if !b.coloured && coloured {
		b.color.EnableColor()
	}
	b.coloured = coloured
	return b.color
}

func (b *Printer) printAllLines() {
	for i := 0; i < b.height; i++ {
		b.printLine(i)
	}
}

func (b *Printer) printLine(i int) {
	b.printer(false).Printf("%s|", strings.Repeat(" ", b.sides))
	for j := 0; j < b.width; j++ {
		b.printChar(i, j)
	}
	b.printer(false).Printf("|%s", strings.Repeat(" ", b.sides))
	b.writer.commitLine()
}

func (b *Printer) printChar(row, col int) {
	on := b.isOn(row, col)
	toPrint := " "
	if !on && b.isOn(row+1, col) {
		toPrint = "_"
	}
	b.printer(on).Print(toPrint)
}

func (b *Printer) isOn(row, col int) bool {
	diff := b.width - len(b.vals)
	return (col >= diff) && (row > (b.height - b.curNorm(b.vals[col-diff])))
}

func (b *Printer) printBlock(n int) {
	for i := 0; i < n; i++ {
		b.printer(false).Print("")
		b.writer.commitLine()
	}
}

func (b *Printer) printRule(c string) {
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
