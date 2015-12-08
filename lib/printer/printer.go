package printer

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type BarPrinter struct {
	representation        [][]bool
	color                 *color.Color
	writer                *lineWriter
	height, top, sides    int
	shouldClear, coloured bool
	stop                  chan struct{}
}

func NewBarPrinter(width, height, top, sides int) *BarPrinter {
	representation := make([][]bool, width-(sides*2)-2)
	for i := 0; i < len(representation); i++ {
		representation[i] = make([]bool, height)
	}

	writer := &lineWriter{Writer: bufio.NewWriter(os.Stdout)}
	colorer := color.New(color.BgBlue)
	color.Output = writer

	return &BarPrinter{
		representation: representation,
		height:         height,
		top:            top,
		sides:          sides,
		color:          colorer,
		writer:         writer,
		stop:           make(chan struct{}),
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
	diff := b.height - count
	bar := make([]bool, b.height)
	for j := 0; j < len(bar); j++ {
		bar[j] = j >= diff
	}
	b.Advance(bar)
}

func (b *BarPrinter) Advance(bar []bool) {
	b.representation = append(b.representation[1:], bar[0:b.height])
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
	for j := 0; j < len(b.representation); j++ {
		b.printChar(j, i)
	}
	b.printer(false).Printf("|%s", strings.Repeat(" ", b.sides))
	b.writer.commitLine()
}

func (b *BarPrinter) printChar(x, y int) {
	bar := b.representation[x]
	toPrint := " "
	if y+1 < len(bar) && bar[y] != bar[y+1] {
		toPrint = "_"
	}
	b.printer(bar[y]).Print(toPrint)
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
	b.printer(false).Printf("%s%s%s", strings.Repeat(" ", b.sides+1), strings.Repeat(c, len(b.representation)), strings.Repeat(" ", b.sides+1))
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
