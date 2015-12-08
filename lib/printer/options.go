package printer

type Option func(*BarPrinter)

func Normalise(b *BarPrinter) {
	normalise := NewBasicNormaliser(float64(b.height))
	b.norm = normalise.Next
}

func Padding(topBottom, leftRight int) Option {
	return func(b *BarPrinter) {
		b.top = topBottom
		b.sides = leftRight
	}
}
