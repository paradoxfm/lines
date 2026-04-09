package algs

type Point struct {
	X, Y float32
}

type Pointi struct {
	X, Y int
}

type Pointf struct {
	X float64
	Y float64
}

type PixelPoint struct {
	X float64
	Y float64
}

type LineEndpoints struct {
	Start PixelPoint
	End   PixelPoint
}

type Cache struct {
	oldWndSize *Pointi

	l1Start *PixelPoint
	l1End   *PixelPoint

	l2Start *PixelPoint
	l2End   *PixelPoint
}

type Algo interface {
	//CalcIntersection() Point

	Init(data string) []string

	Parse(data string)

	Initiated() bool

	//SetInitiated()

	CalcIntersection() Point

	Intersection() *Point

	GetLine1InWindow(wndW, wndH int) (*PixelPoint, *PixelPoint)

	GetLine2InWindow(wndW, wndH int) (*PixelPoint, *PixelPoint)
}
