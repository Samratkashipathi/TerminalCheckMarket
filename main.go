package main

import (
	"log"
	"math"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	headingWidget := widgets.NewParagraph()
	headingWidget.Title = "Name of the Crypto/Stock"
	headingWidget.Text = "BTC"
	headingWidget.SetRect(0, 0, 100, 5)

	sinData := func() [][]float64 {
		n := 220
		data := make([][]float64, 2)
		data[0] = make([]float64, n)
		data[1] = make([]float64, n)
		for i := 0; i < n; i++ {
			data[0][i] = 1 + math.Sin(float64(i)/5)
			data[1][i] = 1 + math.Cos(float64(i)/5)
		}
		return data
	}()

	graphWidget := widgets.NewPlot()
	graphWidget.Title = "1 Day Change"
	graphWidget.Marker = widgets.MarkerDot
	graphWidget.Data = make([][]float64, 2)
	graphWidget.Data[0] = []float64{1, 2, 3, 4, 5}
	graphWidget.Data[1] = sinData[1][4:]
	graphWidget.SetRect(0, 7, 100, 40)
	graphWidget.AxesColor = ui.ColorWhite
	graphWidget.LineColors[0] = ui.ColorCyan
	graphWidget.PlotType = widgets.ScatterPlot

	indianStockWidget := widgets.NewList()
	indianStockWidget.Title = "Stock Watchlist"
	indianStockWidget.Rows = []string{
		"[0] github.com/gizak/termui/v3",
		"[1] Something",
		"[2] New thing",
		"[3] [color](fg:white,bg:green) output",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] foo",
		"[8] bar",
		"[9] baz",
	}
	indianStockWidget.TextStyle = ui.NewStyle(ui.ColorYellow)
	indianStockWidget.WrapText = false
	indianStockWidget.SetRect(110, 0, 140, 17)

	cryptoWidget := widgets.NewList()
	cryptoWidget.Title = "Crypto Watchlist"
	cryptoWidget.Rows = []string{
		"[0] github.com/gizak/termui/v3",
		"[1] Something",
		"[2] New thing",
		"[3] [color](fg:white,bg:green) output",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] foo",
		"[8] bar",
		"[9] baz",
	}
	cryptoWidget.TextStyle = ui.NewStyle(ui.ColorYellow)
	cryptoWidget.WrapText = false
	cryptoWidget.SetRect(110, 20, 140, 38)

	ui.Render(headingWidget, graphWidget, indianStockWidget, cryptoWidget)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}

}
