package printer

import "math"

type NormalisedBarPrinter struct {
	*BarPrinter
	min, max float64
}

func NewNormalisedBarPrinter(width, height, top, sides int) *NormalisedBarPrinter {
	return &NormalisedBarPrinter{
		BarPrinter: NewBarPrinter(width, height, top, sides),
	}
}

func (n *NormalisedBarPrinter) FeedF(fs <-chan float64) {
	for f := range fs {
		n.AdvanceF(f)
	}
}

func (n *NormalisedBarPrinter) AdvanceF(f float64) {
	if math.Ceil(f) > n.max {
		n.max = math.Ceil(f)
	} else if f < n.min {
		n.min = f
	}
	n.AdvanceN(n.normalise(f))
}

func (n *NormalisedBarPrinter) normalise(f float64) int {
	if f == n.max {
		return n.height
	} else if f == n.min {
		return 0
	}

	diff := n.max - n.min
	if diff <= 0 {
		return 0
	}

	return int(math.Floor(((f + math.Abs(n.min)) / diff) * float64(n.height)))
}
