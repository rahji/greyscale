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
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var infile string
var color string
var top int

var scale = []string{
	"Black",
	"Very Dark Gray",
	"Dark Gray",
	"Medium Dark Gray",
	"Slate Gray",
	"Dim Gray",
	"Light Slate Gray",
	"Gray",
	"Light Gray",
	"Gainsboro",
	"Silver",
	"Light Silver",
	"Very Light Gray",
	"Near White",
	"Off White",
	"White",
}

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "shows greyscale info about an image",
	Long: `
The "show" command displays a histogram, showing how much of the
infile image is represented by each of 16 named greyscale colors.

Note that the image is *assumed* to be a greyscale!

One of two additional flags can be used:

--color will return a number representing the percentage of that 
named greyscale color in the image

--top will return the top n highest-frequency colors in the image
`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("infile=%s color=%s top=%d\n", infile, color, top)

		reader, err := os.Open(infile)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		m, _, err := image.Decode(reader)
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
		out.WriteString("|Color|Pixels|Percent|\n")
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
	rootCmd.AddCommand(showCmd)

	showCmd.PersistentFlags().StringVarP(&infile, "infile", "i", "", "input file (required)")
	showCmd.MarkPersistentFlagRequired("infile")
	showCmd.PersistentFlags().StringVarP(&color, "color", "c", "", "greyscale color name (returns percentage)")
	showCmd.PersistentFlags().IntVarP(&top, "top", "t", 0, "number of the highest-frequency colors to show")
}
