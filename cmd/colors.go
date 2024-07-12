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
	"log"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var colorName string
var pixels string
var top int

// colorsCmd represents the colors command
var colorsCmd = &cobra.Command{
	Use:   "colors",
	Short: "show greyscale colors in an image",
	Long: `
The 'show colors' command displays a histogram, showing how much of
the infile image is represented by each of 16 named greyscale colors.

Note that the image is *assumed* to be a greyscale!

One of three additional flags can be used:

--color will return a number representing the percentage of that 
named greyscale color in the image

--top will return the top n highest-frequency colors in the image

--pixels shows the histogram only for pixels in the range x,y:x,y
`,
	Run: func(cmd *cobra.Command, args []string) {

		m, _, err := readImage(infile)
		if err != nil {
			log.Fatal(err)
		}

		bounds := m.Bounds()
		numpix := bounds.Max.X * bounds.Max.Y

		// An image's bounds do not necessarily start at (0, 0), so the two loops start
		// at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
		// likely to result in better memory access patterns than X first and Y second.
		var histogram [16]int
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, _, _, _ := m.At(x, y).RGBA()
				// A color's RGBA method returns values in the range [0, 65535].
				// Shifting by 12 reduces this to the range [0, 15].
				histogram[r>>12]++
			}
		}

		var out strings.Builder
		out.WriteString("|Color Name|Pixels|Percent|\n")
		out.WriteString("|----:|-----:|------:|\n")
		for i, x := range histogram {
			pct := float64(x) / float64(numpix) * 100
			out.WriteString(fmt.Sprintf("|%s|%d|%.02f%%|\n", scale[i], x, pct))
		}
		md, err := glamour.Render(out.String(), "dark")
		fmt.Print(md)
	},
}

func init() {
	showCmd.AddCommand(colorsCmd)
	colorsCmd.PersistentFlags().StringVarP(&colorName, "color", "c", "", "greyscale color name (returns percentage)")
	colorsCmd.PersistentFlags().IntVarP(&top, "top", "t", 0, "number of the highest-frequency colors to show")
	colorsCmd.PersistentFlags().StringVarP(&pixels, "pixels", "p", "", "range of pixels to look at (x,y:x,y)")
}
