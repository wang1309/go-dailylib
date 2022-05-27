package test

import (
	"fmt"
	"github.com/guptarohit/asciigraph"
	"math"
	"testing"
)

// 一个 ascii 图表绘制包
func TestGraph(t *testing.T) {
	data := []float64{3, 4, 9, 6, 2, 4, 5, 8, 5, 10, 2, 7, 2, 5, 6}
	graph := asciigraph.Plot(data)

	fmt.Println(graph)
}

func TestMultiGraph(t *testing.T) {
	data := [][]float64{{0, 1, 2, 3, 3, 3, 2, 0}, {5, 4, 2, 1, 4, 6, 6}}
	graph := asciigraph.PlotMany(data)

	fmt.Println(graph)
}

func TestColorGraph(t *testing.T) {
	data := make([][]float64, 4)

	for i := 0; i < 4; i++ {
		for x := -20; x <= 20; x++ {
			v := math.NaN()
			if r := 20 - i; x >= -r && x <= r {
				v = math.Sqrt(math.Pow(float64(r), 2)-math.Pow(float64(x), 2)) / 2
			}
			data[i] = append(data[i], v)
		}
	}
	graph := asciigraph.PlotMany(data, asciigraph.Precision(0), asciigraph.SeriesColors(
		asciigraph.Red,
		asciigraph.Yellow,
		asciigraph.Green,
		asciigraph.Blue,
	))

	fmt.Println(graph)
}

func TestMultiColorGraph(t *testing.T) {
	data := [][]float64{{5, 1, 5,1,5,5,5,1}}
	graph := asciigraph.PlotMany(data, asciigraph.SeriesColors(
		asciigraph.Red,
		asciigraph.Yellow,
	))

	fmt.Println(graph)
}
