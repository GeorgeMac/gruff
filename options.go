package gruff

type Option func(*Printer)

func Normalise(b *Printer) {
	normalise := NewBasicNormaliser(float64(b.height))
	b.norm = normalise.Next
}

func Padding(topBottom, leftRight int) Option {
	return func(b *Printer) {
		b.top = topBottom
		b.sides = leftRight
	}
}
