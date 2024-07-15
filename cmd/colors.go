/*
Copyright © 2024 Rob Duarte <me@robduarte.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"image"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var colorName string
var pixels string
var top int
var nonzero bool
var csv bool

// colorsCmd represents the colors command
var colorsCmd = &cobra.Command{
	Use:   "colors",
	Short: "show greyscale colors in an image",
	Long: `
The 'show colors' command displays a histogram, showing how much of
the infile image is represented by each of 16 named greyscale colors.

Note that the image is *assumed* to be a greyscale!

`,
	Run: func(cmd *cobra.Command, args []string) {

		m, _, err := readImage(infile)
		if err != nil {
			log.Fatal(err)
		}

		bounds := m.Bounds()
		totalPixels := bounds.Max.X * bounds.Max.Y
		pixelsConsidered := 0

		var minRangeX, minRangeY int
		pixelAmount := totalPixels // will get overwritten if --pixels is used

		if pixels != "" {
			// determine the starting coordinates based on --pixels flag

			// turn xy:n into [xy n] strings
			xyxyn := strings.Split(pixels, ":")
			if len(xyxyn) != 2 {
				log.Fatal(fmt.Errorf("--pixels flag must be specied as x,y:n"))
			}
			// turn xy into [x y] ints
			xyxy := strings.Split(xyxyn[0], ",")
			if len(xyxy) != 2 {
				log.Fatal(fmt.Errorf("--pixels flag must be specied as x,y:n"))
			}
			minRangeX, err = strconv.Atoi(xyxy[0])
			if err != nil {
				log.Fatal(fmt.Errorf("x couldn't be converted to a number in --pixels flag"))
			}
			minRangeY, err = strconv.Atoi(xyxy[1])
			if err != nil {
				log.Fatal(fmt.Errorf("y couldn't be converted to a number in --pixels flag"))
			}
			// turn n into n int
			pixelAmount, err = strconv.Atoi(xyxyn[1])
			if err != nil {
				log.Fatal(fmt.Errorf("number of pixels couldn't be converted to a number in --pixels flag"))
			}

			// validate the actual x and y values
			if minRangeX > bounds.Max.X {
				log.Fatal(fmt.Errorf("x value specifed with --pixels is larger than the image width"))
			}
			if minRangeY > bounds.Max.Y {
				log.Fatal(fmt.Errorf("y value specifed with --pixels is larger than the image width"))
			}
		} else {
			// otherwise, just start normally at the minimum bounds of the image (probably 0,0)
			minRangeX = bounds.Min.X
			minRangeY = bounds.Min.Y
		}

		// end pixels are always the max bounds of the image, but later we will
		// break out of the loop early if --pixels was used
		maxRangeX := bounds.Max.X
		maxRangeY := bounds.Max.Y

		var histogram [16]int
		// An image's bounds do not necessarily start at (0, 0), so the two loops start
		// at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
		// likely to result in better memory access patterns than X first and Y second.
	outer:
		for y := minRangeY; y < maxRangeY; y++ {
			for x := minRangeX; x < maxRangeX; x++ {
				grey := getGrey(m, x, y)
				histogram[grey]++
				pixelsConsidered++
				if pixelAmount == pixelsConsidered {
					break outer
				}
			}
		}

		if colorName != "" {
			colorName = cases.Title(language.Und).String(colorName)
			colorIndex := slices.Index(scale, colorName)
			pct := float64(histogram[colorIndex]) / float64(pixelsConsidered) * 100
			fmt.Println(pct)
			os.Exit(0)
		}

		if top > 0 {
			histogram = topValues(histogram, top)
		}

		var out strings.Builder
		if !csv {
			out.WriteString("# Color Histogram\n")
			out.WriteString("||Color Name|Min Value|Max Value|Pixels|Percent|\n")
			out.WriteString("|:--:|----:|----:|----:|-----:|------:|\n")
		}
		for i, x := range histogram {
			maxGreyValueRange := (i << 4) | 0x0F
			minGreyValueRange := maxGreyValueRange - 15
			pct := float64(x) / float64(pixelsConsidered) * 100
			if (top > 0 || nonzero) && pct == 0 {
				// --top causes the value to be zero, so skip it
				// also skip a zero value if --nonzero was specified
				continue
			}
			var outString string
			if csv {
				outString = "%d,%s,%d,%d,%d,%.02f\n"
			} else {
				outString = "|%d|%s|%3d|%3d|%d|%.02f%%|\n"
			}
			out.WriteString(fmt.Sprintf(outString, i, scale[i], minGreyValueRange, maxGreyValueRange, histogram[i], pct))
		}

		if !csv {
			out.WriteString(fmt.Sprintf("\n*Pixels considered: %d of %d*\n", pixelsConsidered, totalPixels))
			md, _ := glamour.Render(out.String(), "dark")
			fmt.Print(md)
		} else {
			fmt.Print(out.String())
		}
		os.Exit(0)
	},
}

func init() {
	showCmd.AddCommand(colorsCmd)
	colorsCmd.PersistentFlags().StringVarP(&colorName, "color", "c", "", "greyscale color name (returns percentage of that color)")
	colorsCmd.PersistentFlags().IntVarP(&top, "top", "t", 0, "filter the histogram to show only the the highest-frequency colors")
	colorsCmd.PersistentFlags().StringVarP(&pixels, "pixels", "p", "", "range of pixels to look at (x,y:n)")
	colorsCmd.PersistentFlags().BoolVarP(&nonzero, "nonzero", "n", false, "only show non-zero results")
	colorsCmd.PersistentFlags().BoolVarP(&csv, "csv", "r", false, "show raw comma-delimited output")
}

// takes an image width and height, the coordinates of a starting pixel, and an amount of pixels to offset from that point
// returns the new coordinates, clipped at the end of the image
func movePixel(w, h, x, y, n int) (int, int) {
	// calculate the starting position in 1D
	startPos := y*w + x

	// calculate the new position in 1D
	newPos := startPos + n
	fmt.Printf("width:%d height:%d x=%d y=%d newPos=%d\n", w, h, x, y, newPos)
	if newPos >= w*h { // clip coordinates to end of the image
		return w - 1, h - 1
	} else {
		// return new position, converted back to 2D coordinates
		newY := newPos / w
		newX := newPos % w
		return newX, newY
	}
}

// takes an image.Image and a set of coordinates
// returns the amount of "grey" at that pixel (actually the amount of red since we just assume green and blue are the same)
// the return value is shifted 12 bits to the right to put it in the range 0-15
func getGrey(m image.Image, x int, y int) int {
	r, _, _, _ := m.At(x, y).RGBA()
	// A color's RGBA method returns values in the range [0, 65535].
	// Shifting by 12 reduces this to the range [0, 15].
	return int(r >> 12)
}

// topValues returns a 16-item array with only the highest `top` values
// the rest of the items have their value set to 0
func topValues(arr [16]int, top int) [16]int {
	if top >= 16 {
		return arr
	}

	// Create a copy of the array and sort it in descending order
	sortedArr := make([]int, 16)
	copy(sortedArr, arr[:])
	sort.Sort(sort.Reverse(sort.IntSlice(sortedArr)))

	// Take the top values
	topValues := sortedArr[:top]

	// Create a map to keep track of the selected values and their counts
	valueCount := make(map[int]int)
	for _, val := range topValues {
		valueCount[val]++
	}

	// Collect the top values from the original array in the same order
	var result [16]int
	i := 0
	for _, val := range arr {
		if valueCount[val] > 0 {
			result[i] = val
			i++
			valueCount[val]--
		}
	}

	return result
}
