package printer

import "math"

type BasicNormaliser struct {
	min, max, height float64
}

func NewBasicNormaliser(height float64) *BasicNormaliser {
	return &BasicNormaliser{
		height: height,
	}
}

func (n *BasicNormaliser) Next(f float64) func(f float64) int {
	if math.Ceil(f) > n.max {
		n.max = math.Ceil(f)
	} else if f < n.min {
		n.min = f
	}
	return n.normalise
}

func (n *BasicNormaliser) normalise(f float64) int {
	if f == n.max {
		return int(math.Floor(n.height))
	} else if f == n.min {
		return 0
	}

	diff := n.max - n.min
	if diff <= 0 {
		return 0
	}

	return int(math.Floor(((f + math.Abs(n.min)) / diff) * n.height))
}
