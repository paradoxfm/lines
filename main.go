package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"lines/algs"
	"log"
	"os"
	"sync"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
)

var context algs.Algo = new(algs.LineArg)
var (
	mu sync.Mutex
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Title("Lines"),
			app.MinSize(unit.Dp(800), unit.Dp(600)),
			app.MaxSize(unit.Dp(800), unit.Dp(600)))
		//window.Option(app.Size(unit.Dp(800), unit.Dp(600)))
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	expl := explorer.NewExplorer(window)

	var startButton widget.Clickable
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e) // This graphics context is used for managing the rendering state.
			fmt.Println("redraw window")
			addTitle(theme, gtx)
			addLabelCoord(theme, gtx)
			addLines(e, &ops)
			addAxies(e, &ops)
			drawPoint(e, &ops)
			if startButton.Clicked(gtx) {
				go openFile(expl, window)
			}
			layout.Inset{Right: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.Button(theme, &startButton, "Открыть").Layout(gtx)
				})
			})
			e.Frame(gtx.Ops) // Pass the drawing operations to the GPU.
		}
	}
}

func drawPoint(e app.FrameEvent, ops *op.Ops) {

	if !context.Initiated() {
		return
	}
	ip := context.Intersection()

	p := algs.PointToScreen(e.Size.X, e.Size.Y, &algs.Pointf{X: float64(ip.X), Y: float64(ip.Y)})

	cx := float64(p.X)
	cy := float64(p.Y)

	r := 6.0

	rect := image.Rect(int(cx-r), int(cy-r), int(cx+r), int(cy+r))

	circle := clip.Ellipse{Min: rect.Min, Max: rect.Max}.Op(ops)
	paint.FillShape(ops, color.NRGBA{R: 120, G: 120, B: 120, A: 255}, circle)
}

func addTitle(theme *material.Theme, gtx layout.Context) {
	title := material.H6(theme, "В6 График пересечения") // про константы я не слышал, компилятор сам разберется ))
	title.Color = color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	title.Alignment = text.Start
	title.Layout(gtx)
}

func addLabelCoord(th *material.Theme, gtx layout.Context) {
	msgX := "-"
	msgY := "-"
	if context.Initiated() {
		i := context.Intersection()
		msgX = fmt.Sprintf("%.2f", i.X)
		msgY = fmt.Sprintf("%.2f", i.Y)
	}
	layout.NE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		title := material.H6(th, fmt.Sprintf("Point: {X: %s, Y: %s}", msgX, msgY))
		title.Color = color.NRGBA{R: 127, G: 0, B: 0, A: 255}
		return title.Layout(gtx)
	})
}

func addAxies(e app.FrameEvent, ops *op.Ops) {
	sz := e.Size
	offset := 30
	clr := color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	frm := algs.PPi(offset, sz.Y/2)
	to := algs.PPi(sz.X-offset, sz.Y/2)
	drawLine(frm, to, clr, ops)
	frm2 := algs.PPi(sz.X/2, offset)
	to2 := algs.PPi(sz.X/2, sz.Y-offset)
	drawLine(frm2, to2, clr, ops)

	//fmt.Println("Draw asies 1 - ", frm, to)
	//fmt.Println("Draw asies 2 - ", frm2, to2)
}

func addLines(e app.FrameEvent, ops *op.Ops) {
	if !context.Initiated() {
		return
	}

	sz := e.Size
	frm, to := context.GetLine1InWindow(sz.X, sz.Y)
	drawLine(frm, to, color.NRGBA{R: 250, G: 200, B: 200, A: 255}, ops)

	frm2, to2 := context.GetLine2InWindow(sz.X, sz.Y)
	drawLine(frm2, to2, color.NRGBA{R: 200, G: 250, B: 200, A: 255}, ops)

	//fmt.Println("Draw line 1 - ", frm, to)
	//fmt.Println("Draw line 2 - ", frm2, to2)
}

func drawLine(frm *algs.PixelPoint, to *algs.PixelPoint, clr color.NRGBA, ops *op.Ops) {
	var line clip.Path
	line.Begin(ops)
	line.MoveTo(f32.Pt(float32(frm.X), float32(frm.Y)))
	line.LineTo(f32.Pt(float32(to.X), float32(to.Y)))
	line.Close()

	paint.FillShape(ops, clr,
		clip.Stroke{
			Path:  line.End(),
			Width: 2,
		}.Op())
}

func openFile(expl *explorer.Explorer, w *app.Window) {
	file, err := expl.ChooseFile(".txt")
	if err != nil {
		log.Printf("rejected: %v", err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file) // io.ReadAll for cross-platform
	if err != nil {
		log.Printf("failed reading file data: %v", err)
		return
	}
	fmt.Println("File content:", string(data))
	context.Init(string(data))
	point := context.CalcIntersection()
	fmt.Println("Point intersection:", point)
	//mu.Lock() // залочка не нужна. все равно событийная отрисовка. наверно убрать
	w.Option(app.Title("Lines")) // помогает для перерисовки
	//mu.Unlock()
}
