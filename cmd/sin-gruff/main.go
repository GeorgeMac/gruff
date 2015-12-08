package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
	"time"

	"github.com/GeorgeMac/gruff/lib/printer"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	var height, duration int
	flag.IntVar(&height, "h", 20, "Height of bar graph")
	flag.IntVar(&duration, "i", 400, "Number of iterations to radomly generate")
	flag.Parse()

	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan struct{})
	ints := make(chan int)
	printer := printer.NewBarPrinter(w, height, 5, 10)
	go printer.Feed(ints)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		quit, _, _ := reader.ReadRune()
		for quit != 'q' && quit != 'Q' {
			quit, _, _ = reader.ReadRune()
		}
		printer.Stop()
		stop <- struct{}{}
	}()

	go func() {
		for i := 0; i < duration; i++ {
			count := math.Floor(math.Sin(float64(i)/float64(height))*(float64(height)/2)) + (float64(height) / 2)
			ints <- int(count)
			time.Sleep(80 * time.Millisecond)
		}
		close(stop)
	}()
	<-stop
	color.Unset()
}
