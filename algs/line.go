package algs

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type LineArg struct {
	init           bool
	k1, k2, b1, b2 float32
	point          *Point
	cache          *Cache
}

func (l *LineArg) CalcIntersection() Point {
	if l.k2-l.k1 == 0 {
		fmt.Println("No intersection point")
	}
	X := (l.b1 - l.b2) / (l.k2 - l.k1)
	Y := l.k1*X + l.b1
	rez := Point{X, Y}
	l.point = &rez //кешируем
	return rez
}

func (l *LineArg) Intersection() *Point {
	return l.point
}

func (l *LineArg) Init(data string) []string {
	l.reset()
	var lines []string = regexp.MustCompile("\r?\n").Split(data, -1)
	fieldsErr := []string{}
	for _, ln := range lines {
		kv := strings.SplitN(ln, "=", 2)

		switch kv[0] {
		case "K1":
			vr, er := strconv.ParseFloat(kv[1], 32)
			if er == nil {
				l.k1 = float32(vr)
			} else {
				fieldsErr = append(fieldsErr, "K1")
			}
		case "K2":
			vr, er := strconv.ParseFloat(kv[1], 32)
			if er == nil {
				l.k2 = float32(vr)
			} else {
				fieldsErr = append(fieldsErr, "K2")
			}
		case "B1":
			vr, er := strconv.ParseFloat(kv[1], 32)
			if er == nil {
				l.b1 = float32(vr)
			} else {
				fieldsErr = append(fieldsErr, "B1")
			}
		case "B2":
			vr, er := strconv.ParseFloat(kv[1], 32)
			if er == nil {
				l.b2 = float32(vr)
			} else {
				fieldsErr = append(fieldsErr, "B2")
			}
		}
	}
	if len(fieldsErr) == 0 {
		l.init = true // готово у употреблению
	}
	return fieldsErr
}

func parse() float32 {
	return 0
}

func (l *LineArg) reset() {
	l.init = false
	l.k1 = 0
	l.k2 = 0
	l.b1 = 0
	l.b2 = 0
	l.cache = &Cache{}
	l.point = nil
}

func (l *LineArg) Parse(data string) {
	l.init = true
}

func (l *LineArg) Initiated() bool {
	return l.init
}

func DrawLineThroughWindow(wTotalPx, hTotalPx int, p Pointf) LineEndpoints {
	const padding = 30

	// размеры активной области
	activeW := float64(wTotalPx - 2*padding)
	activeH := float64(hTotalPx - 2*padding)

	// центр активной области в пикселях экрана (от 0 до W/H)
	centerX := float64(wTotalPx) / 2.0
	centerY := float64(hTotalPx) / 2.0

	// 3. Переводим относительные координаты P (-100 до 100) в пиксели относительно центра
	// Диапазон -100..100 соответствует половине активной ширины/высоты (например, -400..400 пикселей)
	pixelPerUnitX := activeW / 200.0 // 200 - это полный диапазон (-100 до 100)
	pixelPerUnitY := activeH / 200.0

	Px := p.X * pixelPerUnitX
	Py := p.Y * pixelPerUnitY

	// 4. Находим параметры k и b для формулы y = kx + b (в системе координат относительно центра)
	// Если линия проходит через центр (0,0), k = Py/Px, b = 0.
	// Если Point совпадает с центром (Px=0, Py=0), линия не определена, можно вернуть заглушку или обработать ошибку.
	if Px == 0 && Py == 0 {
		return LineEndpoints{}
	}

	k := Py / Px
	b := float64(0) // Линия проходит через начало координат в этой системе

	// 5. Определяем границы активной области в системе координат относительно центра
	halfW := activeW / 2.0
	halfH := activeH / 2.0

	// 6. Находим крайние точки линии (где она пересекает границы активной области)
	// Проверяем пересечения с вертикальными границами (x = +/- halfW)
	y1 := k*(-halfW) + b
	y2 := k*(halfW) + b

	var start, end PixelPoint

	// Если Y-координаты пересечения находятся в пределах активной высоты (-halfH до halfH)
	if math.Abs(y1) <= halfH && math.Abs(y2) <= halfH {
		start.X = -halfW
		start.Y = y1
		end.X = halfW
		end.Y = y2
	} else {
		// Иначе линия пересекает горизонтальные границы (y = +/- halfH)
		x1 := (halfH - b) / k
		x2 := (-halfH - b) / k
		start.X = x1
		start.Y = halfH
		end.X = x2
		end.Y = -halfH
	}

	// 7. Преобразуем точки обратно в экранные пиксели (от 0 до WTotalPx)
	return LineEndpoints{
		Start: PixelPoint{
			X: start.X + centerX,
			Y: start.Y + centerY,
		},
		End: PixelPoint{
			X: end.X + centerX,
			Y: end.Y + centerY,
		},
	}
}

// вспомогательные функции конвертеры чтоб не загромождать код
func PPf(x, y float32) *PixelPoint {
	return &PixelPoint{float64(x), float64(y)}
}

func PPi(x, y int) *PixelPoint {
	return &PixelPoint{float64(x), float64(y)}
}

func (l *LineArg) GetLine1InWindow(wndW, wndH int) (*PixelPoint, *PixelPoint) {
	// кешируем для максимально быстрой перерисовки если окно поменяло размер
	c := l.cache
	wsz := c.oldWndSize // старый размер окна
	if wsz != nil && wsz.X == wndW && wsz.Y == wndH && c.l1End != nil {
		return c.l1Start, c.l1End
	}
	c.oldWndSize = &Pointi{X: wndW, Y: wndH}
	fm, to := getLineInWindow(wndW, wndH, float64(l.k1), float64(l.b1), l.point)
	c.l1Start = &fm
	c.l1End = &to
	return c.l1Start, c.l1End
}

// можно подрефачить на одну функцию, но сегодня уже лень
// TODO: add params
func (l *LineArg) GetLine2InWindow(wndW, wndH int) (*PixelPoint, *PixelPoint) {
	// кешируем для максимально быстрой перерисовки если окно поменяло размер
	c := l.cache
	wsz := c.oldWndSize // старый размер окна
	if wsz != nil && wsz.X == wndW && wsz.Y == wndH && c.l2End != nil {
		return c.l2Start, c.l2End
	}
	c.oldWndSize = &Pointi{X: wndW, Y: wndH}
	fm, to := getLineInWindow(wndW, wndH, float64(l.k2), float64(l.b2), l.point)
	c.l2Start = &fm
	c.l2End = &to
	return c.l2Start, c.l2End
}

// а вот тут у нас основная магия по преобразованию координат линий к координам окна
// возвращает точки пересечения прямой y = kx + b с границами активной области окна
func getLineInWindow(wndW, wndH int, k, b float64, pRef *Point) (PixelPoint, PixelPoint) {
	const padding = 40.0

	// активная область
	activeW := float64(wndW) - 2*padding
	activeH := float64(wndH) - 2*padding

	// масштаб (пикселей на 1 относительную единицу)
	// середина четверти (pRef.X) соответствует activeW / 4
	scale := (activeW / 4.0) / math.Abs(float64(pRef.X))
	scale2 := (activeH / 4.0) / math.Abs(float64(pRef.Y))

	if scale2 < scale { //выбираем меньший, чтоб все влезлоы
		scale = scale2
	}

	// uраницы активной области в относительных координатах
	// центр окна —>(0,0). границы:
	minRelX := -(activeW / 2.0) / scale
	maxRelX := (activeW / 2.0) / scale
	minRelY := -(activeH / 2.0) / scale
	maxRelY := (activeH / 2.0) / scale

	// gоиск точек пересечения прямой y = kx + b с рамкой окна
	var points []Pointf

	// точка на левой границе
	yLeft := k*minRelX + b
	if yLeft >= minRelY && yLeft <= maxRelY {
		points = append(points, Pointf{minRelX, yLeft})
	}
	// точка на правой границе
	yRight := k*maxRelX + b
	if yRight >= minRelY && yRight <= maxRelY {
		points = append(points, Pointf{maxRelX, yRight})
	}
	// точка на нижней границе
	if k != 0 {
		xBottom := (minRelY - b) / k
		if xBottom > minRelX && xBottom < maxRelX {
			points = append(points, Pointf{xBottom, minRelY})
		}
		// точка на верхней границе
		xTop := (maxRelY - b) / k
		if xTop > minRelX && xTop < maxRelX {
			points = append(points, Pointf{xTop, maxRelY})
		}
	}

	// если линия не попала в окно, хотя такого быть не должно. мб стоит кинуть ошибку?
	if len(points) < 2 {
		return PixelPoint{}, PixelPoint{}
	}

	// перевод из относительных координат в пиксели экрана
	centerX := float64(wndW) / 2.0
	centerY := float64(wndH) / 2.0

	toPixel := func(p Pointf) PixelPoint {
		return PixelPoint{
			X: centerX + (p.X * scale),
			// Y обычно инвертирован (вниз растет),
			// если нужно "математическое" отображение (вверх растет), ставим минус.
			Y: centerY - (p.Y * scale),
		}
	}

	p1 := toPixel(points[0])
	p2 := toPixel(points[1])
	return p1, p2
}

// тут по сути копипаста для расчета и отрисовки самой точки
// мб подрефачить
func PointToScreen(wndW, wndH int, p *Pointf) Pointf {
	const padding = 40.0

	// доступная область
	activeW := float64(wndW) - 2*padding
	activeH := float64(wndH) - 2*padding

	// масштаб (пикселей на 1 относительную единицу)
	// середина четверти (pRef.X) соответствует activeW / 4
	scale := (activeW / 4.0) / math.Abs(float64(p.X))
	scale2 := (activeH / 4.0) / math.Abs(float64(p.Y))

	if scale2 < scale { //выбираем меньший, чтоб все влезлоы
		scale = scale2
	}

	centerX := float64(wndW) / 2.0
	centerY := float64(wndH) / 2.0

	return Pointf{
		X: centerX + (p.X * scale),
		Y: centerY - (p.Y * scale),
	}
}
